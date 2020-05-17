package main

import "github.com/nobonobo/wecty"

func main() {
	wecty.RenderBody(&Top{})
	select {}
}
