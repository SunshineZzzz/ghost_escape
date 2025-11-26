package affiliate

import (
	"ghost_escape/game/core"

	"github.com/go-gl/mathgl/mgl32"
)

// 精灵图
type Sprite struct {
	// 继承基础依附对象
	core.ObjectAffiliate
	// 纹理
	Texture *core.Texture
}

var _ core.IObject = (*Sprite)(nil)

// 初始化
func (s *Sprite) Init() {
	s.ObjectAffiliate.Init()
}

// 渲染
func (s *Sprite) Render() {
	if s.Texture == nil {
		return
	}
	if s.Parent == nil {
		return
	}
	pos := s.Parent.GetRenderPosition().Add(s.Offset)
	core.GetInstance().RenderTexture(s.Texture, pos, s.Size)
}

// 设置纹理
func (s *Sprite) SetTexture(texture *core.Texture) {
	s.Texture = texture
	s.Size = mgl32.Vec2{texture.SrcRect.W, texture.SrcRect.H}
}
