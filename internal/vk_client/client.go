package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	vkAPI      string = "https://api.vk.com/method/"
	apiVersion string = "5.131"
)

type VkClient struct {
	http.Client   `json:"-"`
	ApiToken      string               `json:"-"`
	GroupID       string               `json:"-"`
	Session       sessionData          `json:"response"`
	errorLog      *log.Logger          `json:"-"`
	randomizer    *rand.Rand           `json:"-"`
	keyboards     map[string]*keyboard `json:"-"`
	errorHabdlers map[int]func(int)
}

type sessionData struct {
	Server string `json:"server"`
	Key    string `json:"key"`
	Ts     string `json:"ts"`
	Wait   string `json:"-"`
}

type pollData struct {
	NewTs   string   `json:"ts"`
	Updates []update `json:"updates"`
}

type update struct {
	EventType string                 `json:"type"`
	Object    map[string]interface{} `json:"object"`
}

func New(accessToken, groupID string) *VkClient {
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	vkClient := &VkClient{
		ApiToken:   accessToken,
		GroupID:    groupID,
		errorLog:   errorLog,
		randomizer: r,
		keyboards:  newKeyboardMap(),
	}
	vkClient.setHandlers()
	vkClient.setLongPollServer()
	return vkClient
}

func (c *VkClient) setLongPollServer() {
	params := map[string]string{"group_id": c.GroupID, "v": apiVersion}
	resp, err := c.get(vkAPI, getLongPollServerMethod, params)
	if err != nil {
		c.errorLog.Println(err)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.errorLog.Println(err)
	}
	if err := json.Unmarshal(bodyBytes, c); err != nil {
		c.errorLog.Fatal(err)
	}
	c.Session.Wait = "25"
}

func (c *VkClient) Poll() {
	urlQuery := c.Session.Server
	params := map[string]string{"act": "a_check", "key": c.Session.Key, "ts": c.Session.Ts, "wait": c.Session.Wait}
	resp, err := c.get(urlQuery, "", params)
	if err != nil {
		c.errorLog.Println(err)
		return
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.errorLog.Println(err)
		return
	}
	c.processPollData(bodyBytes)
}

func (c *VkClient) processPollData(bodyBytes []byte) {
	var data pollData
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		c.errorLog.Println(err)
		return
	}
	if strings.HasPrefix(string(bodyBytes), `{"failed":`) {
		c.processError(bodyBytes)
		return
	}
	if len(data.Updates) > 0 {
		c.Session.Ts = data.NewTs
		for _, u := range data.Updates {
			c.processUpdate(&u)
		}
	}
}

func (c *VkClient) processError(bodyBytes []byte) {
	var failed struct {
		Code int `json:"failed"`
		Ts   int `json:"ts"`
	}
	if err := json.Unmarshal(bodyBytes, &failed); err != nil {
		c.errorLog.Println(err)
		return
	}
	if handler, ok := c.errorHabdlers[failed.Code]; !ok {
		c.errorLog.Println("Unknown error code:", failed.Code)
	} else {
		handler(failed.Ts)
	}
}

func (c *VkClient) processUpdate(u *update) {
	if u.EventType == messageNew {
		message := u.Object["message"].(map[string]interface{})
		c.sendMessage(
			uint64(message["from_id"].(float64)),
			greetingMessage,
			"main",
		)
	} else if u.EventType == messageEvent {
		c.processMessageEvent(u.Object)
	} else {
		c.errorLog.Println("event type not implemented: " + u.EventType)
	}
}

func (c *VkClient) processMessageEvent(object map[string]interface{}) {
	c.sendMessageEventAnswer(object)
	userID := uint64(object["user_id"].(float64))
	payload := object["payload"].(map[string]interface{})
	if button := buttonType(payload["button"].(string)); button == returnButton {
		c.sendMessage(userID, greetingMessage, "main")
	} else {
		message := "You Chose Button " + payload["button"].(string) +
			" in layer " + payload["layer"].(string)
		var keyboard string
		if k, ok := payload["keyboard"]; ok {
			keyboard = k.(string)
		}
		c.sendMessage(userID, message, keyboard)
	}
}

func (c *VkClient) sendMessageEventAnswer(object map[string]interface{}) {
	userIdFloat := object["user_id"].(float64)
	userID := fmt.Sprint(uint64(userIdFloat))
	peerID := object["peer_id"].(float64)
	eventID := object["event_id"].(string)
	params := map[string]string{
		"event_id": eventID,
		"user_id":  userID,
		"peer_id":  fmt.Sprint(uint64(peerID)),
	}
	resp, err := c.get(vkAPI, sendMessageEventAnswerMethod, params)
	if err != nil {
		c.errorLog.Println(err)
	}
	resp.Body.Close()
}

func (c *VkClient) sendMessage(userID uint64, message, keyboard string) {
	params := map[string]string{
		"user_id":   fmt.Sprint(userID),
		"random_id": fmt.Sprint(c.randomizer.Int63()),
		"message":   message,
	}
	if keyboard != "" {
		keyboardBytes, err := c.keyboards[keyboard].toJSON()
		if err != nil {
			c.errorLog.Println(err)
			return
		}
		params["keyboard"] = string(keyboardBytes)
	}
	resp, err := c.get(vkAPI, messagesSendMethod, params)
	if err != nil {
		c.errorLog.Println(err)
	}
	resp.Body.Close()
}

func (c *VkClient) get(urlPath, method string, params map[string]string) (*http.Response, error) {
	urlPath += method
	queryParams := url.Values{}
	for k, v := range params {
		queryParams.Set(k, v)
	}
	queryParams.Set("v", apiVersion)
	urlPath += "?" + queryParams.Encode()
	req, err := http.NewRequest("GET", urlPath, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.ApiToken)
	return c.Do(req)
}
