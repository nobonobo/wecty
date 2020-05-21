package views

import (
	"github.com/nobonobo/wecty"
	"github.com/nobonobo/wecty/examples/jsonrpc/frontend/components"
)

// Render ...
func (c *Login) Render() wecty.HTML {
	return wecty.Tag("body", 
		&components.Header{},
		wecty.Tag("main", 			
			wecty.Class{
				"container": true,
			},
			wecty.Tag("form", 				
				wecty.Attr("method", "post"),
				wecty.Attr("action", "/rpc/login"),
				wecty.Attr("enctype", "multipart/form-data"),
				wecty.Class{
					"form-horizontal": true,
				},
				wecty.Event("submit", c.OnSubmit),
				wecty.Tag("div", 					
					wecty.Class{
						"columns": true,
					},
					wecty.Tag("div", 						
						wecty.Class{
							"column": true,
							"col-mx-auto": true,
						},
						wecty.Attr("style", "max-width: 640px;"),
						wecty.Tag("div", 							
							wecty.Class{
								"card": true,
							},
							wecty.Tag("div", 								
								wecty.Class{
									"card-header": true,
								},
								wecty.Tag("div", 									
									wecty.Class{
										"card-title": true,
										"h5": true,
									},
									wecty.Text("Login"),
								),
							),
							wecty.Tag("div", 								
								wecty.Class{
									"card-body": true,
								},
								wecty.Tag("div", 									
									wecty.Class{
										"form-group": true,
									},
									wecty.Tag("div", 										
										wecty.Class{
											"col-3": true,
											"col-sm-12": true,
										},
										wecty.Tag("label", 											
											wecty.Class{
												"form-label": true,
											},
											wecty.Attr("for", "email"),
											wecty.Text("Email:"),
										),
									),
									wecty.Tag("div", 										
										wecty.Class{
											"col-9": true,
											"col-sm-12": true,
										},
										
										wecty.Tag("input", 											
											wecty.Class{
												"form-input": true,
											},
											wecty.Attr("id", "email"),
											wecty.Attr("name", "email"),
											wecty.Attr("type", "email"),
											wecty.Attr("placeholder", "Email"),
											wecty.Attr("pattern", "[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+.[a-zA-Z0-9-.]+$"),
										),
									),
								),
								wecty.Tag("div", 									
									wecty.Class{
										"form-group": true,
									},
									wecty.Tag("div", 										
										wecty.Class{
											"col-3": true,
											"col-sm-12": true,
										},
										wecty.Tag("label", 											
											wecty.Class{
												"form-label": true,
											},
											wecty.Attr("for", "password"),
											wecty.Text("Password:"),
										),
									),
									wecty.Tag("div", 										
										wecty.Class{
											"col-9": true,
											"col-sm-12": true,
										},
										
										wecty.Tag("input", 											
											wecty.Class{
												"form-input": true,
											},
											wecty.Attr("id", "password"),
											wecty.Attr("name", "password"),
											wecty.Attr("type", "password"),
											wecty.Attr("placeholder", "Password"),
											wecty.Attr("pattern", "^([a-zA-Z0-9-_]{8,})$"),
										),
									),
								),
							),
							wecty.Tag("div", 								
								wecty.Class{
									"card-footer": true,
								},
								wecty.Tag("div", 									
									wecty.Class{
										"form-group": true,
									},
									wecty.Tag("div", 										
										wecty.Class{
											"col-3": true,
											"col-sm-12": true,
										},
									),
									wecty.Tag("div", 										
										wecty.Class{
											"col-9": true,
											"col-sm-12": true,
										},
										wecty.Tag("button", 											
											wecty.Class{
												"btn": true,
												"btn-primary": true,
											},
											wecty.Text("Login"),
										),
									),
								),
							),
						),
					),
				),
			),
		),
	)
}
