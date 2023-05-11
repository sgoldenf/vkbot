package client

import "strconv"

func (c *VkClient) setHandlers() {
	c.errorHabdlers = map[int]func(int){
		historyOutdatedError: c.handleHistoryOutdatedError,
		keyExpiredError:      c.handleKeyExpiredError,
		lostInformationError: c.handleLostInformationError,
	}
}

func (c *VkClient) handleHistoryOutdatedError(ts int) {
	c.Session.Ts = strconv.Itoa(ts)
}

func (c *VkClient) handleKeyExpiredError(ts int) {
	c.setLongPollServer()
}

func (c *VkClient) handleLostInformationError(ts int) {
	c.setLongPollServer()
}
