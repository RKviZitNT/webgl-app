//go:build js

package resourcemanager

import "syscall/js"

func LoadImage(path string) js.Value {
	ch := make(chan js.Value, 1)

	promise := js.Global().Call("loadImage", path)
	promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ch <- args[0]
		return nil
	}))
	promise.Call("catch", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		println("Image load error:", args[0].Get("message").String())
		ch <- js.Null()
		return nil
	}))

	return <-ch
}
