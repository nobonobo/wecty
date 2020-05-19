package main

import (
	"github.com/nobonobo/wecty"
)

type Top struct {
	wecty.Core
}

func (c *Top) Render() wecty.HTML {
	return wecty.Tag("body",
		wecty.TagWithNS("svg", "http://www.w3.org/2000/svg",
			wecty.Attr("width", "14"),
			wecty.Attr("height", "16"),
			wecty.Attr("viewBox", "0 0 14 16"),
			wecty.TagWithNS("path", "http://www.w3.org/2000/svg",
				wecty.Attr("fill-rule", "evenodd"),
				wecty.Attr("d", "M7 2.3c3.14 0 5.7 2.56 5.7 5.7s-2.56 5.7-5.7 5.7A5.71 5.71 0 0 1 1.3 8c0-3.14 2.56-5.7 5.7-5.7zM7 1C3.14 1 0 4.14 0 8s3.14 7 7 7 7-3.14 7-7-3.14-7-7-7zm1 3H6v5h2V4zm0 6H6v2h2v-2z"),
			),
		),
	)
}

func main() {
	top := &Top{}
	wecty.RenderBody(top)
	select {}
}
