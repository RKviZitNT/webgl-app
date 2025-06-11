//go:build js

package assetsmanager

import (
	"fmt"
	"sync"
	"syscall/js"
	"webgl-app/internal/graphics/texture"
	"webgl-app/internal/graphics/webgl"
	"webgl-app/internal/jsfunc"
	"webgl-app/internal/resourceloader"
)

type AssetsSources struct {
	MetaDataSrc map[string]string
	TexturesSrc map[string]string
}

var ASrc AssetsSources

type AssetsManager struct {
	metadata map[string]*js.Value
	textures map[string]*texture.Texture
	mu       sync.RWMutex
}

func NewAssetsManager() *AssetsManager {
	ASrc = AssetsSources{
		MetaDataSrc: map[string]string{
			"warrior_anim": "assets/meta/warrior_anim_data.json",
		},
		TexturesSrc: map[string]string{
			"background1": "assets/images/backgrounds/background1.jpg",
			"warrior":     "assets/sprites/warrior/warrior_spritesheet.png",
		},
	}

	return &AssetsManager{
		metadata: make(map[string]*js.Value),
		textures: make(map[string]*texture.Texture),
	}
}

func (a *AssetsManager) Load(glCtx *webgl.GLContext, assetsSrc AssetsSources) error {
	var (
		loadErr error
	)

	js.Global().Call("setLoadingProgress", 10, "Loading metadata...")
	a.loadMetadata(assetsSrc.MetaDataSrc)

	js.Global().Call("setLoadingProgress", 55, "Loading textures...")
	a.loadTextures(glCtx, assetsSrc.TexturesSrc)

	return loadErr
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
			a.addTexture(name, texture.NewTexture(glCtx.GL, img))
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

func (a *AssetsManager) addTexture(name string, texture *texture.Texture) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.textures[name] = texture
}

func (a *AssetsManager) GetTexture(name string) *texture.Texture {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.textures[name]
}
