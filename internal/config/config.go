//go:build js

package config

import (
	"syscall/js"
	"webgl-app/internal/resourceloader"
	"webgl-app/internal/utils"
)

type Window struct {
	Width     float64
	Height    float64
	FrameRate int64
}

type Config struct {
	Debug  bool
	Window Window
}

var ProgramConfig = Config{
	Debug: false,
	Window: Window{
		Width:     1600,
		Height:    900,
		FrameRate: 60,
	},
}

func LoadSources(path string, v any) error {
	var (
		errLoad error
	)

	done := make(chan struct{}, 0)

	resourceloader.LoadFile(path,
		func(src js.Value) {
			errLoad = utils.ParseStringToJSON(src.String(), &v)
			close(done)
		},
		func(err error) {
			errLoad = err
			close(done)
		})

	<-done

	if errLoad != nil {
		return errLoad
	}

	return nil
}
