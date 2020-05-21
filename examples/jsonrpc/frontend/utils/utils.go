package utils

import (
	"syscall/js"

	"github.com/nobonobo/wecty"
)

// Login ...
func Login(form js.Value) {
	println("login")
	fd := js.Global().Get("FormData").New(form)
	js.Global().Call("fetch", form.Get("action"), map[string]interface{}{
		"method": form.Get("method"),
		"body":   fd,
	}).Call("then",
		// success ...
		wecty.Callback1(func(res js.Value) interface{} {
			if res.Get("status").Int() != 200 {
				js.Global().Call("alert", "login failed")
				return nil
			}
			ref := wecty.GetURL().Query().Get("ref")
			if len(ref) > 0 {
				wecty.Navigate(ref)
			} else {
				wecty.Navigate("/")
			}
			return nil
		}),
		// fail ...
		wecty.Callback1(func(err js.Value) interface{} {
			js.Global().Call("alert", err)
			return nil
		}),
	)
}

// Logout ...
func Logout() {
	println("logout")
	js.Global().Call("fetch", "/rpc/logout", map[string]interface{}{
		"method": "post",
	}).Call("then",
		// success ...
		wecty.Callback1(func(res js.Value) interface{} {
			wecty.Navigate("/logout?ref=/")
			return nil
		}),
		// fail ...
		wecty.Callback1(func(err js.Value) interface{} {
			js.Global().Call("alert", err)
			return nil
		}),
	)
}
