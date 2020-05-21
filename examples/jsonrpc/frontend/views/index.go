package views

import (
	"syscall/js"

	"github.com/nobonobo/wecty"
	"github.com/nobonobo/wecty/examples/jsonrpc/frontend/components"
)

//go:generate wecty generate -c Index -p views index.html

// Index ...
type Index struct {
	wecty.Core
	header *components.Header
}

func (c *Index) Header() *components.Header {
	if c.header == nil {
		c.header = &components.Header{}
	}
	return c.header
}

func (c *Index) Mount() {
	client := js.Global().Get("JsonRpcClient").New()
	client.Call("request", "Service.User").Call("then",
		// success ...
		wecty.Callback1(func(res js.Value) interface{} {
			js.Global().Get("console").Call("log", res)
			c.Header().UserName = res.Get("Name").String()
			go wecty.Rerender(c.Header())
			return nil
		}),
		// fail ...
		wecty.Callback1(func(err js.Value) interface{} {
			js.Global().Get("console").Call("log", err)
			wecty.Navigate("/login?ref=/")
			return nil
		}),
	)
}
