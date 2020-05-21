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
				wecty.Text("JSONRPC-Sample"),
			),
		),
		wecty.Tag("section", 			
			wecty.Class{
				"navbar-section": true,
			},
			wecty.Tag("div", 				
				wecty.Class{
					"dropdown": true,
				},
				wecty.Tag("a", 					
					wecty.Attr("href", "#"),
					wecty.Class{
						"btn": true,
						"btn-link": true,
						"dropdown-toggle": true,
					},
					wecty.Attr("tabindex", "0"),
					wecty.Text(c.UserName),
					wecty.Tag("i", 						
						wecty.Class{
							"icon": true,
							"icon-caret": true,
						},
					),
				),
				wecty.Tag("ul", 					
					wecty.Class{
						"menu": true,
					},
					wecty.Tag("li", 
						wecty.Tag("a", 							
							wecty.Attr("href", "#"),
							wecty.Event("click", c.OnClickLogout),
							wecty.Tag("i", 								
								wecty.Class{
									"icon": true,
									"icon-shutdown": true,
								},
							),
							wecty.Text("Logout"),
						),
					),
				),
			),
		),
	)
}
