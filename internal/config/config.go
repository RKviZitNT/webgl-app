//go:build js

package config

import (
	"syscall/js"
	"webgl-app/internal/jsfunc"
	"webgl-app/internal/resourceloader"
	"webgl-app/internal/utils"
)

// ----- Program config -----

type Window struct {
	Width     float64 `json:"width"`
	Height    float64 `json:"height"`
	FrameRate int64   `json:"frame_rate"`
}

type ProgramConfig struct {
	Debug  bool   `json:"debug"`
	Window Window `json:"window"`
}

// ----- Shaders config -------

type Shaders struct {
	Vertex   string `json:"vertex"`
	Fragment string `json:"fragment"`
}

type ShadersConfig struct {
	TextureShaders Shaders `json:"texture_shaders"`
	DebugShaders   Shaders `json:"debug_shaders"`
}

//  ----- Assets config  -----

type AssetsConfig struct {
	MetaDataSrc map[string]string `json:"metadata"`
	TexturesSrc map[string]string `json:"textures"`
}

// -----------------------

var (
	ProgramConf ProgramConfig
	ShadersConf ShadersConfig
	AssetsConf  AssetsConfig
)

func LoadConfigs(programConfigPath, shadersConfigPath, assetsConfigPath string) error {
	var (
		loadErr error
	)

	done := make(chan struct{}, 0)

	resourceloader.LoadFile(programConfigPath,
		func(src js.Value) {
			loadErr = utils.ParseStringToJSON(src.String(), &ProgramConf)

			resourceloader.LoadFile(shadersConfigPath,
				func(src js.Value) {
					loadErr = utils.ParseStringToJSON(src.String(), &ShadersConf)

					resourceloader.LoadFile(assetsConfigPath,
						func(src js.Value) {
							loadErr = utils.ParseStringToJSON(src.String(), &AssetsConf)
							close(done)
						},
						func(err error) {
							loadErr = err
							close(done)
						})
				},
				func(err error) {
					loadErr = err
					close(done)
				})
		},
		func(err error) {
			loadErr = err
			close(done)
		})

	<-done

	if loadErr != nil {
		return loadErr
	}

	jsfunc.LogInfo("Configs loaded")

	return nil
}
