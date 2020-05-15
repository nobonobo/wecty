package main

import (
	"github.com/nobonobo/wecty"
)

// Render ...
func (c *TopView) Render() wecty.HTML {
	return wecty.Tag("form", 		
		wecty.Event("submit", c.OnSubmit),
		wecty.Tag("button", 
			wecty.Text("Submit"),
		),
	)
}
