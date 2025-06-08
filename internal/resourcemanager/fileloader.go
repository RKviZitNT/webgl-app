//go:build js

package resourcemanager

import (
	"errors"
	"syscall/js"
)

func LoadFile(path string, onSuccess FileCallback, onError ErrorCallback) {
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
