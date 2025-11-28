package world

import (
	"ghost_escape/game/affiliate"
	"ghost_escape/game/core"

	"github.com/go-gl/mathgl/mgl32"
)

// 特效
type Effect struct {
	// 继承基础世界对象
	core.ObjectWorld
	// 精灵动画
	spriteAnim core.IObjectAnima
	// 动画结束后需要添加到场景中的对象
	nextObject core.IObjectWorld
}

var _ core.IObject = (*Effect)(nil)
var _ core.IObjectScreen = (*Effect)(nil)
var _ core.IObjectWorld = (*Effect)(nil)

// 创建特效
func AddEffect(parent core.IObject, filePath string, pos mgl32.Vec2, scale float32, anchorType core.AnchorType, nextObject core.IObjectWorld) *Effect {
	effect := &Effect{}
	effect.Init()
	effect.spriteAnim = affiliate.AddSpriteAnimChild(effect, filePath, scale, anchorType)
	effect.spriteAnim.SetLoop(false)
	effect.SetPosition(pos)
	effect.SetNextObject(nextObject)
	if parent != nil {
		parent.AddChild(effect)
	}
	return effect
}

// 更新
func (s *Effect) Update(dt float32) {
	s.ObjectWorld.Update(dt)
	s.checkFinish()
}

// 检查特效是否播放完毕
func (s *Effect) checkFinish() {
	if s.spriteAnim.GetFinish() {
		s.NeedRemove = true
		if s.nextObject != nil {
			s.Game().GetCurrentScene().SafeAddChild(s.nextObject)
		}
	}
}

// 设置特效播放完成后需要添加到场景中的对象
func (s *Effect) SetNextObject(nextObject core.IObjectWorld) {
	s.nextObject = nextObject
}

// 获取特效播放完成后需要添加到场景中的对象
func (s *Effect) GetNextObject() core.IObjectWorld {
	return s.nextObject
}
