package views

import (
	"github.com/nobonobo/wecty"
)

// Render ...
func (c *Index) Render() wecty.HTML {
	return wecty.Tag("body", 
		c.Header(),
		wecty.Tag("main", 			
			wecty.Class{
				"container": true,
			},
			wecty.Tag("h2", 
				wecty.Text("Index:"),
				wecty.Text(wecty.GetURL().String()),
			),
		),
	)
}
