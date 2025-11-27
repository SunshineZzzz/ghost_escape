package affiliate

import (
	"ghost_escape/game/core"

	"github.com/go-gl/mathgl/mgl32"
)

// 精灵图动画
type SpriteAnim struct {
	// 继承精灵图
	Sprite
	// 动画总帧数
	totalFrame float32
	// 动画帧率，1秒多少张图片，totalFrame/fps = 播放时间
	fps float32
	// 动画当前帧
	currentFrame float32
	// 动画帧计时器
	frameTimer float32
}

var _ core.IObject = (*Sprite)(nil)
var _ ISprite = (*Sprite)(nil)

func AddSpriteAnimChild(parent core.IObjectScreen, filePath string, scale float32) *SpriteAnim {
	child := &SpriteAnim{}
	child.Init()
	child.SetTexture(core.CreateTexture(filePath))
	child.SetScale(scale)
	child.SetParent(parent)
	parent.AddChild(child)
	return child
}

// 初始化
func (s *SpriteAnim) Init() {
	s.Sprite.Init()
	s.fps = 10.0
}

// 更新
func (s *SpriteAnim) Update(dt float32) {
	s.Sprite.Update(dt)
	s.frameTimer += dt
	// 动画帧计时器超过播放一帧所需的时间时，切换到下一帧
	if s.frameTimer >= 1.0/s.fps {
		s.currentFrame++
		// 当前帧超过总帧数时，重置当前帧为0
		if s.currentFrame >= s.totalFrame {
			s.currentFrame = 0.0
		}
		// 重置动画帧计时器
		s.frameTimer = 0.0
	}
	s.Texture.SrcRect.X = s.Texture.SrcRect.W * s.currentFrame
}

// 设置纹理
func (s *SpriteAnim) SetTexture(texture *core.Texture) {
	s.Texture = texture
	s.totalFrame = texture.SrcRect.W / texture.SrcRect.H
	texture.SrcRect.W = texture.SrcRect.H
	s.Size = mgl32.Vec2{texture.SrcRect.W, texture.SrcRect.H}
}

// 获取当前帧
func (s *SpriteAnim) GetCurrentFrame() float32 {
	return s.currentFrame
}

// 设置当前帧
func (s *SpriteAnim) SetCurrentFrame(frame float32) {
	s.currentFrame = frame
}

// 获取总帧数
func (s *SpriteAnim) GetTotalFrame() float32 {
	return s.totalFrame
}

// 设置总帧数
func (s *SpriteAnim) SetTotalFrame(totalFrame float32) {
	s.totalFrame = totalFrame
}

// 获取帧率
func (s *SpriteAnim) GetFps() float32 {
	return s.fps
}

// 设置帧率
func (s *SpriteAnim) SetFps(fps float32) {
	s.fps = fps
}

// 获取动画帧计时器
func (s *SpriteAnim) GetFrameTimer() float32 {
	return s.frameTimer
}

// 设置动画帧计时器
func (s *SpriteAnim) SetFrameTimer(frameTimer float32) {
	s.frameTimer = frameTimer
}
