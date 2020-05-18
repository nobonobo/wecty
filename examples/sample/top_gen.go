package main

import (
	"github.com/nobonobo/wecty"
)

// Render ...
func (c *Top) Render() wecty.HTML {
	return wecty.Tag("body", 
		wecty.Tag("form", 			
			wecty.Event("submit", c.OnSubmit),
			wecty.Tag("input", 				
				wecty.Attr("id", "name"),
				wecty.Attr("type", "text"),
			),
		),
		wecty.Tag("h2", 
			wecty.Text(c.text),
		),
	)
}
