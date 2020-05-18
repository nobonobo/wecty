package main

import (
	"log"

	"github.com/nobonobo/wecty"
)

func main() {
	log.SetFlags(log.Lshortfile)
	wecty.RenderBody(&Top{})
	select {}
}
