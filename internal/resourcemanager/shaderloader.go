//go:build js

package resourcemanager

import (
	"errors"
	"syscall/js"
)

func LoadShader(path string, onSuccess ShaderCallback, onError ErrorCallback) {
	promise := js.Global().Call("loadShader", path)

	promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if onSuccess != nil {
			onSuccess(args[0].String())
		}
		return nil
	}))

	promise.Call("catch", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if onError != nil {
			onError(errors.New(args[0].Get("message").String()))
		}
		return nil
	}))
}

// func LoadShader(path string) string {
// 	ch := make(chan string, 1)

// 	promise := js.Global().Call("loadShader", path)
// 	promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
// 		ch <- args[0].String()
// 		return nil
// 	}))
// 	promise.Call("catch", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
// 		println("Shader load error:", args[0].Get("message").String())
// 		ch <- ""
// 		return nil
// 	}))

// 	return <-ch
// }
