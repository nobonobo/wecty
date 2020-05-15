package main

import (
	"github.com/nobonobo/wecty"
)

func main() {
	wecty.RenderBody(wecty.Tag("body", &TopView{}))
	select {}
}
