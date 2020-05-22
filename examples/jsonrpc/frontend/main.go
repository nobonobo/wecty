package main

import (
	"syscall/js"

	"github.com/nobonobo/wecty"
	"github.com/nobonobo/wecty/examples/jsonrpc/frontend/views"
)

func main() {
	wecty.AddStylesheet("/assets/app.css")
	wecty.AddStylesheet("https://unpkg.com/spectre.css/dist/spectre.min.css")
	wecty.AddStylesheet("https://unpkg.com/spectre.css/dist/spectre-exp.min.css")
	wecty.AddStylesheet("https://unpkg.com/spectre.css/dist/spectre-icons.min.css")
	obj := <-wecty.LoadModule([]string{"JsonRpcClient"}, "/assets/jsonrpcclient.js")
	js.Global().Set("JsonRpcClient", obj)
	router := wecty.NewRouter()
	router.Handle("/", func(key string) {
		wecty.RenderBody(&views.Index{})
	})
	router.Handle("/login", func(key string) {
		wecty.RenderBody(&views.Login{})
	})
	if err := router.Start(); err != nil {
		println(err)
		wecty.RenderBody(wecty.NotFoundPage())
	}
	select {}
}
