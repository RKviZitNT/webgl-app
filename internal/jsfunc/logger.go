//go:build js

package jsfunc

import "syscall/js"

func LogInfo(msg string) {
	js.Global().Get("console").Call("log", msg)
}

func LogWarn(msg string) {
	js.Global().Get("console").Call("warn", msg)
}

func LogError(msg string) {
	js.Global().Get("console").Call("error", msg)
}
