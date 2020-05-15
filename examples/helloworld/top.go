package main

import (
	"syscall/js"

	"github.com/nobonobo/wecty"
)

//go:generate wecty generate -c TopView -p main top.html

// TopView ...
type TopView struct {
	wecty.Core
}

// OnSubmit ...
func (c *TopView) OnSubmit(ev js.Value) interface{} {
	ev.Call("preventDefault")
	println("submit!")
	return nil
}
