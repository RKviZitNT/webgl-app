//go:build js

package resourceloader

import (
	"errors"
	"syscall/js"
)

func LoadFile(path string, onSuccess fileCallback, onError errorCallback) {
	promise := js.Global().Call("loadFile", path)

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
