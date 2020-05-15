package components

import (
	"syscall/js"

	"github.com/nobonobo/wecty"
)

//go:generate wecty generate -c Item -p components item.html

// Item ...
type Item struct {
	wecty.Core
	Title string
}

// OnClick ...
func (c *Item) OnClick(ev js.Value) interface{} {
	println("click:", c.Title)
	return nil
}
