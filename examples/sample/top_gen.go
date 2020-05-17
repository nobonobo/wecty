package main

import (
	"github.com/nobonobo/wecty"
)

// Render ...
func (c *Top) Render() wecty.HTML {
	return wecty.Tag("body", 
		wecty.Tag("button", 			
			wecty.Event("click", c.OnClick),
			wecty.Text("Submit"),
		),
		wecty.Tag("h2", 
			wecty.Text(c.text),
		),
	)
}
