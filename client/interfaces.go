package client

import (
	"giveaway/data/owner"
)

type SuspendsThread interface {
	Sleep()
}

type AcceptsThreadSuspender interface {
	SetSuspender(SuspendsThread)
}

type HasDateAttribute interface {
	GetCreationDate() int64
}

type HasOwner interface {
	GetOwner() *owner.Owner
}

type ITask interface {
	FetchData()
	//StopDataFetching()
	DropData()
	DecideWinner()
}
