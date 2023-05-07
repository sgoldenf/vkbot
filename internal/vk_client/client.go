package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const (
	vkAPI                   string = "https://api.vk.com/method/"
	apiVersion              string = "5.131"
	getLongPollServerMethod string = "groups.getLongPollServer"
)

type VkClient struct {
	http.Client `json:"-"`
	ApiToken    string      `json:"-"`
	Session     sessionData `json:"response"`
	errorLog    *log.Logger `json:"-"`
	randomizer  *rand.Rand  `json:"-"`
	params      *jsonParams `json:"-"`
}

type sessionData struct {
	Server string `json:"server"`
	Key    string `json:"key"`
	Ts     string `json:"ts"`
	Wait   string `json:"-"`
}

type jsonParams struct {
	firstLayerKeyboardJSON string
}

type pollData struct {
	NewTs   string   `json:"ts"`
	Updates []update `json:"updates"`
}

type update struct {
	EventType event                  `json:"type"`
	Object    map[string]interface{} `json:"object"`
}

func New(accessToken, groupID string) *VkClient {
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	vkClient := &VkClient{ApiToken: accessToken, errorLog: errorLog, randomizer: r}
	vkClient.createJSONParams()
	params := map[string]string{"group_id": groupID, "v": apiVersion}
	resp, err := vkClient.get(vkAPI, getLongPollServerMethod, params)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		errorLog.Fatal(err)
	}
	if err := json.Unmarshal(bodyBytes, &vkClient); err != nil {
		errorLog.Fatal(err)
	}
	vkClient.Session.Wait = "25"
	return vkClient
}

func (c *VkClient) createJSONParams() {
	bytes, err := json.Marshal(firstLayerKeyboard)
	if err != nil {
		c.errorLog.Fatal(err)
	}
	c.params = &jsonParams{firstLayerKeyboardJSON: string(bytes)}
}

func (c *VkClient) get(urlPath, method string, params map[string]string) (*http.Response, error) {
	urlPath += method
	querySymbol := "?"
	for k, v := range params {
		urlPath += querySymbol + k + "=" + v
		querySymbol = "&"
	}
	req, err := http.NewRequest("GET", urlPath, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.ApiToken)
	return c.Do(req)
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
	var data pollData
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		c.errorLog.Println(err)
		return
	}
	c.Session.Ts = data.NewTs
	for _, u := range data.Updates {
		c.processUpdate(&u)
	}
}

func (c *VkClient) processUpdate(u *update) {
	if u.EventType == messageNew {
		c.welcome(u.Object)
	} else {
		fmt.Println(u.EventType) // TODO: process message_event
	}
}

func (c *VkClient) welcome(object map[string]interface{}) {
	if message, ok := object["message"].(map[string]interface{}); ok {
		userIdFloat := message["from_id"].(float64)
		params := map[string]string{
			"user_id":   fmt.Sprint(uint64(userIdFloat)),
			"random_id": fmt.Sprint(c.randomizer.Int63()),
			"message":   "Welcome!",
			"keyboard":  c.params.firstLayerKeyboardJSON,
			"v":         apiVersion,
		}
		_, err := c.get(vkAPI, "messages.send", params)
		if err != nil {
			c.errorLog.Println(err)
		}
	}
}
