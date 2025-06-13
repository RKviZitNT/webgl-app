//go:build js

package assetsmanager

import (
	"fmt"
	"sync"
	"syscall/js"
	"webgl-app/internal/graphics/webgl"
	"webgl-app/internal/jsfunc"
	"webgl-app/internal/resourceloader"
	"webgl-app/internal/utils"
)

type AssetsSources struct {
	MetaDataSrc map[string]string `json:"metadata"`
	TexturesSrc map[string]string `json:"textures"`
}

var ASrc AssetsSources

type AssetsManager struct {
	metadata map[string]*js.Value
	textures map[string]*webgl.Texture
	mu       sync.RWMutex
}

func NewAssetsManager() *AssetsManager {
	return &AssetsManager{
		metadata: make(map[string]*js.Value),
		textures: make(map[string]*webgl.Texture),
	}
}

func (a *AssetsManager) Load(glCtx *webgl.GLContext, assetsConfig string) error {
	var (
		loadErr error
	)

	js.Global().Call("setLoadingProgress", 10, "Loading assets config...")
	assetsSrc, err := a.loadAssetsConfig(assetsConfig)
	if err != nil {
		return err
	}

	js.Global().Call("setLoadingProgress", 10, "Loading metadata...")
	err = a.loadMetadata(assetsSrc.MetaDataSrc)
	if err != nil {
		return err
	}

	js.Global().Call("setLoadingProgress", 55, "Loading textures...")
	err = a.loadTextures(glCtx, assetsSrc.TexturesSrc)
	if err != nil {
		return err
	}

	return loadErr
}

func (a *AssetsManager) loadAssetsConfig(assetsConfig string) (*AssetsSources, error) {
	var (
		loadErr   error
		assetsSrc AssetsSources
	)

	done := make(chan struct{}, 0)

	resourceloader.LoadFile(assetsConfig,
		func(src js.Value) {
			loadErr = utils.ParseStringToJSON(src.String(), &assetsSrc)
			close(done)
		},
		func(err error) {
			loadErr = err
			close(done)
		})

	<-done

	if loadErr != nil {
		return nil, loadErr
	}

	return &assetsSrc, nil
}

func (a *AssetsManager) loadMetadata(srcPaths map[string]string) error {
	var (
		loadErr  error
		errMutex sync.Mutex
		wg       sync.WaitGroup
	)

	total := len(srcPaths)
	wg.Add(total)

	resourceloader.LoadFiles(srcPaths,
		func(name string, src js.Value) {
			a.addMetadata(name, src)
			wg.Done()
		},
		func(err error) {
			errMutex.Lock()
			if loadErr == nil {
				loadErr = err
			}
			errMutex.Unlock()
			wg.Done()
		},
		func(loaded int) {
			jsfunc.LogInfo(fmt.Sprintf("Loaded %d/%d metadata", loaded, total))
			js.Global().Call("setLoadingProgress", 10+(float64(loaded)/float64(total)*45), fmt.Sprintf("Loading metadata... %d/%d", loaded, total))
		})

	wg.Wait()

	return loadErr
}

func (a *AssetsManager) loadTextures(glCtx *webgl.GLContext, srcPaths map[string]string) error {
	var (
		loadErr  error
		errMutex sync.Mutex
		wg       sync.WaitGroup
	)

	total := len(srcPaths)
	wg.Add(total)

	resourceloader.LoadImages(srcPaths,
		func(name string, img js.Value) {
			a.addTexture(name, webgl.NewTexture(glCtx.GL, img))
			wg.Done()
		},
		func(err error) {
			errMutex.Lock()
			if loadErr == nil {
				loadErr = err
			}
			errMutex.Unlock()
			wg.Done()
		},
		func(loaded int) {
			jsfunc.LogInfo(fmt.Sprintf("Loaded %d/%d textures", loaded, total))
			js.Global().Call("setLoadingProgress", 55+(float64(loaded)/float64(total)*45), fmt.Sprintf("Loading textures... %d/%d", loaded, total))
		})

	wg.Wait()

	return loadErr
}

func (a *AssetsManager) addMetadata(name string, src js.Value) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.metadata[name] = &src
}

func (a *AssetsManager) GetMetadata(name string) *js.Value {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.metadata[name]
}

func (a *AssetsManager) addTexture(name string, texture *webgl.Texture) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.textures[name] = texture
}

func (a *AssetsManager) GetTexture(name string) *webgl.Texture {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.textures[name]
}
