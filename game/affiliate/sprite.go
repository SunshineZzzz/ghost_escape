package affiliate

import (
	"ghost_escape/game/core"

	"github.com/go-gl/mathgl/mgl32"
)

// 精灵图抽象
type ISprite interface {
	// 设置纹理
	SetTexture(*core.Texture)
	// 设置反转
	SetFlip(bool)
	// 获取是否反转
	GetFlip() bool
	// 设置角度
	SetAngle(float64)
	// 获取角度
	GetAngle() float64
}

// 精灵图组件
type Sprite struct {
	// 继承基础依附对象
	core.ObjectAffiliate
	// 纹理
	Texture *core.Texture
}

var _ core.IObject = (*Sprite)(nil)
var _ ISprite = (*Sprite)(nil)

// 添加精灵图子对象
func AddSpriteChild(parent core.IObjectWorld, filePath string, scale float32) *Sprite {
	child := &Sprite{}
	child.Init()
	child.SetTexture(core.CreateTexture(filePath))
	child.SetScale(scale)
	child.SetParent(parent)
	parent.AddChild(child)
	return child
}

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

// 设置反转
func (s *Sprite) SetFlip(flip bool) {
	s.Texture.IsFlip = flip
}

// 获取是否反转
func (s *Sprite) GetFlip() bool {
	return s.Texture.IsFlip
}

// 设置角度
func (s *Sprite) SetAngle(angle float64) {
	s.Texture.Angle = angle
}

// 获取角度
func (s *Sprite) GetAngle() float64 {
	return s.Texture.Angle
}

// 非接口实现
