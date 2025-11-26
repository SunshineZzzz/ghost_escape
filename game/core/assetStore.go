package core

import (
	"fmt"
	"strconv"

	"github.com/SunshineZzzz/purego-sdl3/img"
	"github.com/SunshineZzzz/purego-sdl3/sdl"
	"github.com/SunshineZzzz/purego-sdl3/ttf"
)

// 资源管理器
type AssetStore struct {
	// SDL渲染器
	sdlRenderer *sdl.Renderer
	// 存储所有加载的纹理
	textures map[string]*sdl.Texture
	// 存储所有加载的字体
	fonts map[string]*ttf.Font
	// 存储所有加载的声音
	sounds map[string]ISound
}

// 创建资源管理器
func CreateAssetStore(sdlRenderer *sdl.Renderer) *AssetStore {
	return &AssetStore{
		sdlRenderer: sdlRenderer,
		textures:    make(map[string]*sdl.Texture),
		fonts:       make(map[string]*ttf.Font),
		sounds:      make(map[string]ISound),
	}
}

// 清理
func (a *AssetStore) Clean() {
	for _, texture := range a.textures {
		sdl.DestroyTexture(texture)
	}
	for _, font := range a.fonts {
		ttf.CloseFont(font)
	}
	for _, sound := range a.sounds {
		sound.Close()
	}
	a.textures = make(map[string]*sdl.Texture)
	a.fonts = make(map[string]*ttf.Font)
	a.sounds = make(map[string]ISound)
}

// 载入图片素材
func (a *AssetStore) loadImage(filePath string) error {
	texture := img.LoadTexture(a.sdlRenderer, filePath)
	if texture == nil {
		return fmt.Errorf("load image error,%s", sdl.GetError())
	}
	a.textures[filePath] = texture
	return nil
}

// 载入声音素材
func (a *AssetStore) loadSound(filePath string) error {
	sound, err := NewSound(filePath)
	if err != nil {
		return fmt.Errorf("load sound error,%s", err.Error())
	}
	a.sounds[filePath] = sound
	return nil
}

// 载入字体素材
func (a *AssetStore) loadFont(filePath string, fontSize float32) error {
	font := ttf.OpenFont(filePath, fontSize)
	if font == nil {
		return fmt.Errorf("load font error,%s", sdl.GetError())
	}
	a.fonts[filePath+strconv.Itoa(int(fontSize))] = font
	return nil
}

// 获取图片素材
func (a *AssetStore) GetImage(filePath string) (*sdl.Texture, error) {
	t, ok := a.textures[filePath]
	if ok {
		return t, nil
	}
	err := a.loadImage(filePath)
	if err != nil {
		return nil, err
	}
	return a.textures[filePath], nil
}

// 获取声音素材
func (a *AssetStore) GetSound(filePath string) (ISound, error) {
	sound, ok := a.sounds[filePath]
	if ok {
		return sound, nil
	}
	err := a.loadSound(filePath)
	if err != nil {
		return nil, err
	}
	return a.sounds[filePath], nil
}

// 获取字体素材
func (a *AssetStore) GetFont(filePath string, fontSize float32) (*ttf.Font, error) {
	font, ok := a.fonts[filePath+strconv.Itoa(int(fontSize))]
	if ok {
		return font, nil
	}
	err := a.loadFont(filePath, fontSize)
	if err != nil {
		return nil, err
	}
	return a.fonts[filePath+strconv.Itoa(int(fontSize))], nil
}
