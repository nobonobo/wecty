package main

import (
	"log"
	"syscall/js"
	"time"

	"github.com/nobonobo/wecty"
)

type Sub struct {
	wecty.Core
}

func (c *Sub) Mount() {
	log.Print("sub mount")
}

func (c *Sub) Unmount() {
	log.Print("sub unmount")
}

func (c *Sub) Render() wecty.HTML {
	return wecty.Tag("h3", wecty.Text("Hello!"))
}

//go:generate wecty generate -c Top top.html

type Top struct {
	wecty.Core
	text string
}

func (c *Top) OnSubmit(ev js.Value) interface{} {
	ev.Call("preventDefault")
	c.text = time.Now().Format(time.RFC3339Nano)
	wecty.Rerender(c)
	return nil
}

func (c *Top) Mount() {
	log.Print("top mount")
}

func (c *Top) Unmount() {
	log.Print("top unmount")
}
