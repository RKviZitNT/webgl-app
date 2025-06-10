//go:build js

package resourceloader

import "syscall/js"

type fileCallback func(source string)
type imageCallback func(img js.Value)
type imagesCallback func(name string, img js.Value)
type progressCallback func(loaded int)
type errorCallback func(err error)
