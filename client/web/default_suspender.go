package web

import "time"

type defaultSuspender struct {

}

func (c *defaultSuspender) Sleep() {
	time.Sleep(time.Second * 2)
}
