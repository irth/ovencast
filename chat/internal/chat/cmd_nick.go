package chat

import "regexp"

// TODO: support something more fun than ascii only
var NickRegex = regexp.MustCompile(`^[a-zA-Z0-9\._-]{3,32}$`)

func validateNick(n string) bool {
	return NickRegex.MatchString(n)
}

func (c *Chat) handleNick(state *ClientState, cmd NickCommand) {
	currentNick := state.nick
	newNick := cmd.Request

	if newNick == currentNick {
		cmd.OK(nil)
		return
	}

	if !validateNick(newNick) {
		cmd.Err("nick invalid")
		return
	}

	c.nicksLock.Lock()
	defer c.nicksLock.Unlock()
	if _, ok := c.nicks[newNick]; ok {
		cmd.Err("nick already taken")
		return
	}

	c.nicks[newNick] = struct{}{}
	delete(c.nicks, currentNick)

	state.nick = newNick
	state.ready = true // allow sending messages after a nick has been set

	cmd.OK(nil)
	// TODO: broadcast nick change
}
