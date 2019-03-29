package client

import "giveaway/data"

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
	GetOwner() *data.Owner
}
