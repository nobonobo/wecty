package components

import (
	"syscall/js"

	"github.com/nobonobo/wecty"
	"github.com/nobonobo/wecty/examples/jsonrpc/frontend/utils"
)

//go:generate wecty generate -c Header -p components header.html

// Header ...
type Header struct {
	wecty.Core
	UserName string
}

// OnClickLogout ...
func (c *Header) OnClickLogout(ev js.Value) interface{} {
	utils.Logout()
	return nil
}
