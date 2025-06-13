//go:build js

package config

import (
	"syscall/js"
	"webgl-app/internal/resourceloader"
	"webgl-app/internal/utils"
)

type Window struct {
	Width     float64 `json:"width"`
	Height    float64 `json:"height"`
	FrameRate int64   `json:"frame_rate"`
}

type Config struct {
	Window Window `json:"window"`
}

var GlobalConfig Config

func LoadConfig(path string) error {
	var (
		loadErr error
	)

	done := make(chan struct{}, 0)

	resourceloader.LoadFile(path,
		func(src js.Value) {
			loadErr = utils.ParseStringToJSON(src.String(), &GlobalConfig)
			close(done)
		},
		func(err error) {
			loadErr = err
			close(done)
		})

	<-done

	if loadErr != nil {
		return loadErr
	}

	return nil
}
