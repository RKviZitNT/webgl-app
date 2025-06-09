//go:build js

package resourceloader

import (
	"errors"
	"syscall/js"
)

func LoadImage(path string, onSuccess imageCallback, onError errorCallback) {
	promise := js.Global().Call("loadImage", path)

	promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if onSuccess != nil {
			onSuccess(args[0])
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

func LoadImages(assetsSrc map[string]string, onSuccess imagesCallback, onError errorCallback, onProgress progressCallback) {
	total := len(assetsSrc)
	loaded := 0

	for name, path := range assetsSrc {
		promise := js.Global().Call("loadImage", path)

		promise.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			loaded++
			if onProgress != nil {
				onProgress(total, loaded)
			}
			if onSuccess != nil {
				onSuccess(name, args[0])
			}
			return nil
		}))

		promise.Call("catch", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if onProgress != nil {
				onProgress(total, loaded)
			}
			if onError != nil {
				onError(errors.New(args[0].Get("message").String()))
			}
			return nil
		}))
	}
}
