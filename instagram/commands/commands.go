package commands

import (
	"giveaway/client/api"
	"giveaway/instagram/account"
	"net/http"
)

type Command interface {
	Handle()
	SetChannel(chan interface{})
}



type ReLoginAccountCommand struct {
	acc *account.Account
	ch chan interface{}
}

func (c *ReLoginAccountCommand) SetChannel(channel chan interface{}) {
	c.ch = channel
}

func (c *ReLoginAccountCommand) Handle() {
	c.acc.Cookies = make(map[string][]*http.Cookie, 0)
	cl := api.NewApiClient(c.acc.Proxy)
	cl.SetAccount(c.acc)
	cl.QeSync()
	cl.LauncherSync()
	cl.Login()
	c.ch <- c.acc
	close(c.ch)
}