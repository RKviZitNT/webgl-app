//go:build js

package resourcemanager

import "syscall/js"

func LoadShader(path string) string {
	ch := make(chan string, 1)

	promise := js.Global().Call("loadShader", path)
	promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ch <- args[0].String()
		return nil
	}))
	promise.Call("catch", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		println("Shader load error:", args[0].Get("message").String())
		ch <- ""
		return nil
	}))

	return <-ch
}
