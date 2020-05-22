package wecty

import "syscall/js"

// Callback0 ...
func Callback0(fn func() interface{}) js.Func {
	var cb js.Func
	cb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer cb.Release()
		return fn()
	})
	return cb
}

// Callback1 ...
func Callback1(fn func(res js.Value) interface{}) js.Func {
	var cb js.Func
	cb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer cb.Release()
		return fn(args[0])
	})
	return cb
}

// CallbackN ...
func CallbackN(fn func(res []js.Value) interface{}) js.Func {
	var cb js.Func
	cb = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer cb.Release()
		return fn(args)
	})
	return cb
}
