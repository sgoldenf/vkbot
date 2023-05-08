package client

import "encoding/json"

type keyboard struct {
	Buttons [][]button `json:"buttons"`
}

type button struct {
	Action action `json:"action"`
}

type action struct {
	Type    string                 `json:"type"`
	Label   string                 `json:"label"`
	Payload map[string]interface{} `json:"payload"`
}

type buttonType string

const (
	returnButton buttonType = "return"
)

func newKeyboardMap() map[string]*keyboard {
	returnButton := button{
		Action: action{
			Type:    "callback",
			Label:   "Return",
			Payload: map[string]interface{}{"button": "return"},
		},
	}
	return map[string]*keyboard{
		"main": {
			Buttons: [][]button{
				{
					{
						Action: action{
							Type:    "callback",
							Label:   "Button 1",
							Payload: map[string]interface{}{"layer": "main", "button": "1", "keyboard": "1"},
						},
					},
					{
						Action: action{
							Type:    "callback",
							Label:   "Button 2",
							Payload: map[string]interface{}{"layer": "main", "button": "2", "keyboard": "2"},
						},
					},
					{
						Action: action{
							Type:    "callback",
							Label:   "Button 3",
							Payload: map[string]interface{}{"layer": "main", "button": "3", "keyboard": "3"},
						},
					},
					{
						Action: action{
							Type:    "callback",
							Label:   "Button 4",
							Payload: map[string]interface{}{"layer": "main", "button": "4", "keyboard": "4"},
						},
					},
				},
			},
		},
		"1": {
			Buttons: [][]button{
				{
					{
						Action: action{
							Type:    "callback",
							Label:   "Button 1.1",
							Payload: map[string]interface{}{"layer": "1", "button": "1.1"},
						},
					},
					{
						Action: action{
							Type:    "callback",
							Label:   "Button 1.2",
							Payload: map[string]interface{}{"layer": "1", "button": "1.2"},
						},
					},
				},
				{
					returnButton,
				},
			},
		},
		"2": {
			Buttons: [][]button{
				{
					{
						Action: action{
							Type:    "callback",
							Label:   "Button 2.1",
							Payload: map[string]interface{}{"layer": "2", "button": "2.1"},
						},
					},
					{
						Action: action{
							Type:    "callback",
							Label:   "Button 2.2",
							Payload: map[string]interface{}{"layer": "2", "button": "2.2"},
						},
					},
				},
				{
					returnButton,
				},
			},
		},
		"3": {
			Buttons: [][]button{
				{
					{
						Action: action{
							Type:    "callback",
							Label:   "Button 3.1",
							Payload: map[string]interface{}{"layer": "3", "button": "3.1"},
						},
					},
					{
						Action: action{
							Type:    "callback",
							Label:   "Button 3.2",
							Payload: map[string]interface{}{"layer": "3", "button": "3.2"},
						},
					},
				},
				{
					returnButton,
				},
			},
		},
		"4": {
			Buttons: [][]button{
				{
					{
						Action: action{
							Type:    "callback",
							Label:   "Button 4.1",
							Payload: map[string]interface{}{"layer": "4", "button": "4.1"},
						},
					},
					{
						Action: action{
							Type:    "callback",
							Label:   "Button 4.2",
							Payload: map[string]interface{}{"layer": "4", "button": "4.2"},
						},
					},
				},
				{
					returnButton,
				},
			},
		},
	}
}

func (k *keyboard) toJSON() ([]byte, error) {
	return json.Marshal(k)
}
