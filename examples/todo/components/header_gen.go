package components

import (
	"github.com/nobonobo/wecty"
)

// Render ...
func (c *Header) Render() wecty.HTML {
	return wecty.Tag("header", 		
		wecty.Class{
			"navbar": true,
		},
		wecty.Tag("section", 			
			wecty.Class{
				"navbar-section": true,
			},
			wecty.Tag("a", 				
				wecty.Attr("href", "#/"),
				wecty.Class{
					"navbar-brand": true,
					"mr-2": true,
				},
				wecty.Text("TODO"),
			),
			wecty.Tag("a", 				
				wecty.Attr("href", "#/login"),
				wecty.Class{
					"btn": true,
					"btn-link": true,
				},
				wecty.Text("Login"),
			),
			wecty.Tag("a", 				
				wecty.Attr("href", "#/"),
				wecty.Class{
					"btn": true,
					"btn-link": true,
				},
				wecty.Text("List"),
			),
		),
	)
}
