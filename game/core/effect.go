package core

import "github.com/go-gl/mathgl/mgl32"

// 特效组件
type Effect struct {
	// 继承基础世界对象
	ObjectWorld
	// 精灵动画
	spriteAnim IObjectAnima
	// 动画结束后需要添加到场景中的对象
	nextObject IObjectWorld
}

var _ IObject = (*Effect)(nil)
var _ IObjectScreen = (*Effect)(nil)
var _ IObjectWorld = (*Effect)(nil)

// 创建特效组件
func AddEffect(parent IObject, spriteAnim IObjectAnima, pos mgl32.Vec2, nextObject IObjectWorld) *Effect {
	effect := &Effect{}
	effect.ObjectWorld.Init()
	effect.spriteAnim = spriteAnim
	effect.spriteAnim.SetLoop(false)
	effect.SetPosition(pos)
	effect.SetNextObject(nextObject)
	if parent != nil {
		parent.AddChild(effect)
	}
	if spriteAnim.GetParent() == nil {
		spriteAnim.SetParent(effect)
		effect.AddChild(spriteAnim)
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
func (s *Effect) SetNextObject(nextObject IObjectWorld) {
	s.nextObject = nextObject
}

// 获取特效播放完成后需要添加到场景中的对象
func (s *Effect) GetNextObject() IObjectWorld {
	return s.nextObject
}
