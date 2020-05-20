package main

import (
	"github.com/nobonobo/wecty"
)

func main() {
	//log.SetFlags(log.Lshortfile)
	wecty.RenderBody(&Top{})
	select {}
}
