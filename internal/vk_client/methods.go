package client

type method string

const (
	getLongPollServerMethod      method = "groups.getLongPollServer"
	messagesSendMethod           method = "messages.send"
	sendMessageEventAnswerMethod method = "messages.sendMessageEventAnswer"
)
