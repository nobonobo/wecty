package main

import (
	"fmt"
	"syscall/js"

	"github.com/nobonobo/wecty"
)

type Sub struct {
	wecty.Core
}

func (c *Sub) String() string {
	return fmt.Sprintf("sub(%p)", c)
}

func (c *Sub) Mount() {
	println("sub mount")
}

func (c *Sub) Unmount() {
	println("sub unmount")
}

func (c *Sub) Render() wecty.HTML {
	return wecty.Tag("h3", wecty.Text("Hello!"))
}

//go:generate wecty generate -c Top top.html

type Top struct {
	wecty.Core
	text  string
	sub   *Sub
	count int
}

func (c *Top) OnSubmit(ev js.Value) interface{} {
	ev.Call("preventDefault")
	c.count++
	c.text = fmt.Sprintf("hoge:%d", c.count) //time.Now().Format(time.RFC3339Nano)
	wecty.Rerender(c)
	return nil
}

func (c *Top) Mount() {
	//log.Print("top mount")
}

func (c *Top) Unmount() {
	//log.Print("top unmount")
}
