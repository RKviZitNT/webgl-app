//go:build js

package resourceloader

import "syscall/js"

type fileCallback func(source js.Value)
type filesCallback func(name string, source js.Value)
type imageCallback func(img js.Value)
type imagesCallback func(name string, img js.Value)
type progressCallback func(loaded int)
type errorCallback func(err error)
