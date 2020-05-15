package components

import (
	"github.com/nobonobo/wecty"
)

// Render ...
func (c *Item) Render() wecty.HTML {
	return wecty.Tag("li", 
		wecty.Tag("label", 			
			wecty.Class{
				"form-checkbox": true,
				"form-inline": true,
			},
			
			wecty.Tag("input", 				
				wecty.Attr("type", "checkbox"),
				wecty.Event("click", c.OnClick),
			),
			wecty.Tag("i", 				
				wecty.Class{
					"form-icon": true,
				},
			),
			wecty.Text(c.Title),
		),
	)
}
