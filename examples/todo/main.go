package main

import (
	"github.com/nobonobo/wecty"
	"github.com/nobonobo/wecty/examples/todo/views"
)

func main() {
	/*
		wecty.AddMeta("viewport", "width=device-width")
		wecty.AddStylesheet("https://unpkg.com/spectre.css/dist/spectre.min.css")
		wecty.AddStylesheet("https://unpkg.com/spectre.css/dist/spectre-exp.min.css")
		wecty.AddStylesheet("https://unpkg.com/spectre.css/dist/spectre-icons.min.css")
		wecty.AddStylesheet("assets/app.css")
	*/
	wecty.AddScript("assets/jsonformdata.js")
	router := wecty.NewRouter()
	router.Handle("/login", func(key string) {
		wecty.SetTitle("Login")
		wecty.RenderBody(&views.Form{})
	})
	router.Handle("/", func(key string) {
		wecty.SetTitle("TodoList")
		wecty.RenderBody(&views.Todo{})
	})
	if err := router.Start(); err != nil {
		println(err)
		wecty.RenderBody(wecty.NotFoundPage())
	}
	select {}
}
