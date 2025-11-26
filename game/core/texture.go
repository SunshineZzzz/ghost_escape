package core

import (
	"github.com/SunshineZzzz/purego-sdl3/sdl"
)

// 纹理
type Texture struct {
	// 底层纹理对象
	Texture *sdl.Texture
	// 纹理原始矩形区域，目标渲染矩形区域在ObjectAffiliate中已经有了
	SrcRect sdl.FRect
	// 角度
	Angle float64
	// 是否反转
	IsFlip bool
}

// 创建纹理
func CreateTexture(filePath string) *Texture {
	texture, err := GetInstance().GetAssetStore().GetImage(filePath)
	if err != nil {
		panic(err)
	}
	var w, h float32
	sdl.GetTextureSize(texture, &w, &h)
	return &Texture{
		Texture: texture,
		SrcRect: sdl.FRect{X: 0.0, Y: 0.0, W: w, H: h},
	}
}
