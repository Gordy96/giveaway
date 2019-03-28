package client

import (
	"giveaway/instagram/account"
)

type WorksWithAccount interface {
	GetAccount() *account.Account
	SetAccount(*account.Account)
}

type AuthenticatesAccount interface {
	WorksWithAccount
	Login() (bool, error)
}


type SuspendsThread interface {
	Sleep()
}

type AcceptsThreadSuspender interface {
	SetSuspender(SuspendsThread)
}

type IRule interface {
	Validate(interface{}) (bool, error)
}

type HasDateAttribute interface {
	GetCreationDate() int32
}