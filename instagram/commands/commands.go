package commands

import (
	"giveaway/client/api"
	"giveaway/instagram/account"
	"giveaway/instagram/account/repository"
	"giveaway/utils/logger"
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

	success, err := cl.Login()
	if success {
		c.acc.Status = account.Available
		logger.DefaultLogger().Infof("account %s (%s) is now logged in", c.acc.Username, c.acc.Id)
	} else {
		c.acc.Status = account.CheckPoint
		logger.DefaultLogger().Infof("account %s (%s) can`t login (%s)", c.acc.Username, c.acc.Id, err.Error())
	}
	repo := repository.GetRepositoryInstance()
	repo.Save(c.acc)
	c.ch <- c.acc
	close(c.ch)
}

func MakeNewReLoginCommand(acc *account.Account) *ReLoginAccountCommand {
	logger.DefaultLogger().Infof("account %s (%s) needs re-login", acc.Username, acc.Id)
	c := &ReLoginAccountCommand{}
	c.acc = acc
	repo := repository.GetRepositoryInstance()
	c.acc.Status = account.Maintenance
	repo.Save(acc)
	return c
}
