//go:build js

package jsfunc

import "syscall/js"

func LogInfo(msg interface{}) {
	js.Global().Get("console").Call("log", msg)
}

func LogWarn(msg interface{}) {
	js.Global().Get("console").Call("warn", msg)
}

func LogError(msg interface{}) {
	js.Global().Get("console").Call("error", msg)
}
