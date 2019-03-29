package commands

import (
	"giveaway/client/api"
	"giveaway/instagram/account"
	"giveaway/instagram/account/repository"
	"net/http"
)

type Command interface {
	Handle()
	SetChannel(chan interface{})
}

type ReLoginAccountCommand struct {
	acc *account.Account
	ch  chan interface{}
}

func (c *ReLoginAccountCommand) SetChannel(channel chan interface{}) {
	c.ch = channel
}

func (c *ReLoginAccountCommand) Handle() {
	c.acc.Cookies = make(map[string][]*http.Cookie, 0)
	cl := api.NewApiClient()
	cl.SetAccount(c.acc)

	success, _ := cl.Login()
	if success {
		c.acc.Status = account.LoggedIn
	} else {
		c.acc.Status = account.CheckPoint
	}
	repo := repository.GetRepositoryInstance()
	repo.Save(c.acc)
	c.ch <- c.acc
	close(c.ch)
}

func MakeNewReLoginCommand(acc *account.Account) *ReLoginAccountCommand {
	c := &ReLoginAccountCommand{}
	c.acc = acc
	repo := repository.GetRepositoryInstance()
	c.acc.Status = account.Maintenance
	repo.Save(acc)
	return c
}
