package views

import (
	"syscall/js"

	"github.com/nobonobo/wecty"
	"github.com/nobonobo/wecty/examples/jsonrpc/frontend/utils"
)

//go:generate wecty generate -c Login -p views login.html

// Login ...
type Login struct {
	wecty.Core
}

// OnSubmit ...
func (c *Login) OnSubmit(ev js.Value) interface{} {
	ev.Call("preventDefault")
	utils.Login(ev.Get("target"))
	return nil
}
