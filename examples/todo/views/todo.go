package views

import (
	"github.com/nobonobo/wecty"
	"github.com/nobonobo/wecty/examples/todo/components"
)

//go:generate wecty generate -c Todo -p views todo.html

// Todo ...
type Todo struct {
	wecty.Core
	items components.ItemList
}

// Header ...
func (c *Todo) Header() wecty.Markup {
	return &components.Header{}
}

// Items ...
func (c *Todo) Items() wecty.Markup {
	return c.items.Items()
}
