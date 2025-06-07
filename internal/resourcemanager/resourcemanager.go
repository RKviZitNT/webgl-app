//go:build js

package resourcemanager

import "syscall/js"

type ShaderCallback func(source string)
type ImageCallback func(img js.Value)
type ProgressCallback func(total, loaded, loadErrors int)
type ErrorCallback func(err error)
