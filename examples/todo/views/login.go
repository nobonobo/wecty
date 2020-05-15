package views

import (
	"syscall/js"

	"github.com/nobonobo/wecty"
	"github.com/nobonobo/wecty/examples/todo/components"
)

//go:generate wecty generate -c Form -p views login.html

// Form ...
type Form struct {
	wecty.Core
}

// Header ...
func (c *Form) Header() wecty.Markup {
	return &components.Header{}
}

// OnSubmit ...
func (c *Form) OnSubmit(ev js.Value) interface{} {
	ev.Call("preventDefault")
	println("submit!")
	return nil
}
