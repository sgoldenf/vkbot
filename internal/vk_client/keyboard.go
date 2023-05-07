package client

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

var (
	firstLayerButton1 button = button{
		Action: action{
			Type:    "callback",
			Label:   "Button1",
			Payload: map[string]interface{}{"button": "1"},
		},
	}
	firstLayerButton2 button = button{
		Action: action{
			Type:    "callback",
			Label:   "Button2",
			Payload: map[string]interface{}{"button": "2"},
		},
	}
	firstLayerButton3 button = button{
		Action: action{
			Type:    "callback",
			Label:   "Button3",
			Payload: map[string]interface{}{"button": "3"},
		},
	}
	firstLayerButton4 button = button{
		Action: action{
			Type:    "callback",
			Label:   "Button4",
			Payload: map[string]interface{}{"button": "4"},
		},
	}
	firstLayerKeyboard = keyboard{
		Buttons: [][]button{{firstLayerButton1, firstLayerButton2, firstLayerButton3, firstLayerButton4}},
	}
)
