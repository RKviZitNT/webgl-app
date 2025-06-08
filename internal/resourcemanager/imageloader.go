//go:build js

package resourcemanager

import (
	"errors"
	"syscall/js"
)

func LoadImage(path string, onSuccess ImageCallback, onError ErrorCallback) {
	promise := js.Global().Call("loadImage", path)

	promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		onSuccess(args[0])
		return nil
	}))

	promise.Call("catch", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		onError(errors.New(args[0].Get("message").String()))
		return nil
	}))
}

func LoadImages(paths []string, onSuccess ImageCallback, onError ErrorCallback, onProgress ProgressCallback) {
	total := len(paths)
	loaded := 0
	loadErrors := 0

	for _, path := range paths {
		promise := js.Global().Call("loadImage", path)

		promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			loaded++
			if onProgress != nil {
				onProgress(total, loaded, loadErrors)
			}
			if onSuccess != nil {
				onSuccess(args[0])
			}
			return nil
		}))

		promise.Call("catch", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if onProgress != nil {
				onProgress(total, loaded, loadErrors)
			}
			if onError != nil {
				onError(errors.New(args[0].Get("message").String()))
			}
			return nil
		}))
	}
}
