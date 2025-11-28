package affiliate

import (
	"ghost_escape/game/core"

	"github.com/go-gl/mathgl/mgl32"
)

// 精灵图动画组件
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
	// 是否循环播放
	loop bool
	// 是否播放完毕
	isFinish bool
}

var _ core.IObject = (*SpriteAnim)(nil)
var _ ISprite = (*SpriteAnim)(nil)
var _ core.IObjectAnima = (*SpriteAnim)(nil)
var _ core.IObjectAffiliate = (*SpriteAnim)(nil)

// 添加精灵图动画组件到父对象中
func AddSpriteAnimChild(parent core.IObjectWorld, filePath string, scale float32, anchorType core.AnchorType) *SpriteAnim {
	child := &SpriteAnim{}
	child.Init()
	child.SetTexture(core.CreateTexture(filePath))
	child.SetScale(scale)
	child.SetParent(parent)
	child.SetAnchorType(anchorType)
	parent.AddChild(child)
	return child
}

// 创建精灵图动画组件
func CteateSpriteAnimChild(filePath string, scale float32, anchorType core.AnchorType) *SpriteAnim {
	child := &SpriteAnim{}
	child.Init()
	child.SetTexture(core.CreateTexture(filePath))
	child.SetScale(scale)
	child.SetAnchorType(anchorType)
	return child
}

// 初始化
func (s *SpriteAnim) Init() {
	s.Sprite.Init()
	s.fps = 10.0
	s.loop = true
	s.isFinish = false
}

// 更新
func (s *SpriteAnim) Update(dt float32) {
	if s.isFinish {
		return
	}

	s.Sprite.Update(dt)
	s.frameTimer += dt
	// 动画帧计时器超过播放一帧所需的时间时，切换到下一帧
	if s.frameTimer >= 1.0/s.fps {
		s.currentFrame++
		// 当前帧超过总帧数时，重置当前帧为0
		if s.currentFrame >= s.totalFrame {
			s.currentFrame = 0.0
			// 如果不是循环播放，标记为播放完毕
			if !s.loop {
				s.isFinish = true
				return
			}
		}
		// 重置动画帧计时器
		s.frameTimer = 0.0
	}
	s.Texture.SrcRect.X = s.Texture.SrcRect.W * s.currentFrame
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

// 获取是否循环播放
func (s *SpriteAnim) GetLoop() bool {
	return s.loop
}

// 设置是否循环播放
func (s *SpriteAnim) SetLoop(loop bool) {
	s.loop = loop
}

// 获取是否播放完毕
func (s *SpriteAnim) GetFinish() bool {
	return s.isFinish
}

// 设置是否播放完毕
func (s *SpriteAnim) SetFinish(finish bool) {
	s.isFinish = finish
}

// 非接口实现

// 设置纹理
func (s *SpriteAnim) SetTexture(texture *core.Texture) {
	s.Texture = texture
	s.totalFrame = texture.SrcRect.W / texture.SrcRect.H
	texture.SrcRect.W = texture.SrcRect.H
	s.Size = mgl32.Vec2{texture.SrcRect.W, texture.SrcRect.H}
}
