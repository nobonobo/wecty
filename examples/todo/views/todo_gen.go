package views

import (
	"github.com/nobonobo/wecty"
	"github.com/nobonobo/wecty/examples/todo/components"
)

// Render ...
func (c *Todo) Render() wecty.HTML {
	return wecty.Tag("body", 
		&components.Header{},
		wecty.Tag("main", 			
			wecty.Class{
				"container": true,
			},
			c.Items(),
		),
	)
}
