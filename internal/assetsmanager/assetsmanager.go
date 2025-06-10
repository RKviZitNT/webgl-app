//go:build js

package assetsmanager

import (
	"fmt"
	"sync"
	"syscall/js"
	"webgl-app/internal/graphics/texture"
	"webgl-app/internal/graphics/webgl"
	"webgl-app/internal/resourceloader"
)

type AssetsManager struct {
	textures map[string]*texture.Texture
	mu       sync.RWMutex
}

var AssetsSrc = map[string]string{
	"background1": "assets/images/backgrounds/background1.jpg",
	"warrior":     "assets/sprites/warrior/spritesheet.png",
}

func NewAssetsManager() *AssetsManager {
	return &AssetsManager{
		textures: make(map[string]*texture.Texture),
	}
}

func (a *AssetsManager) Load(glCtx *webgl.GLContext, assetsSrc map[string]string) error {
	var (
		loadErr  error
		errMutex sync.Mutex
		wg       sync.WaitGroup
	)

	total := len(assetsSrc)
	wg.Add(total)

	js.Global().Call("setLoadingProgress", 10, "Loading assets...")
	resourceloader.LoadImages(assetsSrc,
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
			js.Global().Get("console").Call("log", fmt.Sprintf("Loaded %d/%d assets", loaded, total))
			js.Global().Call("setLoadingProgress", 10+(float64(loaded)/float64(total)*90), "Loading assets...")
		})

	wg.Wait()
	return loadErr
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
