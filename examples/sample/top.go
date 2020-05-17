package main

import (
	"syscall/js"
	"time"

	"github.com/nobonobo/wecty"
)

//go:generate wecty generate -c Top top.html

type Top struct {
	wecty.Core
	text string
}

func (c *Top) OnClick(ev js.Value) interface{} {
	c.text = time.Now().Format(time.RFC3339Nano)
	wecty.Rerender(c)
	return nil
}
