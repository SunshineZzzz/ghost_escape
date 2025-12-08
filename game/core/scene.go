package core

import (
	"container/list"

	"github.com/SunshineZzzz/purego-sdl3/sdl"
	"github.com/go-gl/mathgl/mgl32"
)

// 场景抽象
type IScene interface {
	IObject
	// 世界坐标转换为屏幕坐标
	WorldToScreen(mgl32.Vec2) mgl32.Vec2
	// 屏幕坐标转换为世界坐标
	ScreenToWorld(mgl32.Vec2) mgl32.Vec2
	// 获取摄像机位置
	GetCameraPosition() mgl32.Vec2
	// 设置摄像机位置
	SetCameraPosition(mgl32.Vec2)
	// 获取世界大小
	GetWorldSize() mgl32.Vec2
	// 获取世界对象孩子
	GetChildWorld() *list.List
	// 获取屏幕对象孩子
	GetChildScreen() *list.List
}

// 基础场景
type Scene struct {
	// 继承基础对象
	Object
	// 世界大小
	WorldSize mgl32.Vec2
	// 摄像机位置
	CameraPositon mgl32.Vec2
	// 世界对象孩子
	ChildrenWorld list.List
	// 屏幕对象孩子
	ChildrenScreen list.List
	// 是否暂停
	IsPause bool
}

var _ IObject = (*Scene)(nil)
var _ IScene = (*Scene)(nil)

// 世界坐标转换为屏幕坐标
func (s *Scene) WorldToScreen(worldPosition mgl32.Vec2) mgl32.Vec2 {
	// 世界坐标-摄像机位置=屏幕坐标
	return worldPosition.Sub(s.CameraPositon)
}

// 屏幕坐标转换为世界坐标
func (s *Scene) ScreenToWorld(screenPosition mgl32.Vec2) mgl32.Vec2 {
	// 屏幕坐标+摄像机位置=世界坐标
	return screenPosition.Add(s.CameraPositon)
}

// 获取摄像机位置(世界坐标系)
func (s *Scene) GetCameraPosition() mgl32.Vec2 {
	return s.CameraPositon
}

// 设置摄像机位置(世界坐标系)
func (s *Scene) SetCameraPosition(pos mgl32.Vec2) {
	s.CameraPositon[0] = mgl32.Clamp(pos.X(), -30.0, s.WorldSize.X()-s.Game().GetScreenSize().X()+30.0)
	s.CameraPositon[1] = mgl32.Clamp(pos.Y(), -30.0, s.WorldSize.Y()-s.Game().GetScreenSize().Y()+30.0)
}

// 获取世界大小
func (s *Scene) GetWorldSize() mgl32.Vec2 {
	return s.WorldSize
}

// 初始化
func (s *Scene) Init() {
	s.Object.Init()
	s.ChildrenWorld.Init()
	s.ChildrenScreen.Init()
	// 我这里重写了AddChild所以需要设置Self
	s.Object.Self = s
	s.IsPause = false
}

// 处理事件
func (s *Scene) HandleEvent(event *sdl.Event) {
	for e := s.ChildrenScreen.Front(); e != nil; e = e.Next() {
		if e.Value.(IObject).GetIsActive() {
			e.Value.(IObject).HandleEvent(event)
		}
	}
	if s.IsPause {
		return
	}
	s.Object.HandleEvent(event)
	for e := s.ChildrenWorld.Front(); e != nil; e = e.Next() {
		if e.Value.(IObject).GetIsActive() {
			e.Value.(IObject).HandleEvent(event)
		}
	}
}

// 更新
func (s *Scene) Update(dt float32) {
	if !s.IsPause {
		s.Object.Update(dt)
		for e := s.ChildrenWorld.Front(); e != nil; {
			next := e.Next()
			if e.Value.(IObject).GetNeedRemove() {
				s.ChildrenWorld.Remove(e)
				e.Value.(IObject).SetActive(false)
			}
			if e.Value.(IObject).GetIsActive() {
				e.Value.(IObject).Update(dt)
			}
			e = next
		}
	}

	for e := s.ChildrenScreen.Front(); e != nil; {
		next := e.Next()
		if e.Value.(IObject).GetNeedRemove() {
			s.ChildrenScreen.Remove(e)
			e.Value.(IObject).SetActive(false)
		}
		if e.Value.(IObject).GetIsActive() {
			e.Value.(IObject).Update(dt)
		}
		e = next
	}
}

// 渲染
func (s *Scene) Render() {
	s.Object.Render()
	for e := s.ChildrenWorld.Front(); e != nil; e = e.Next() {
		if e.Value.(IObject).GetIsActive() {
			e.Value.(IObject).Render()
		}
	}
	for e := s.ChildrenScreen.Front(); e != nil; e = e.Next() {
		if e.Value.(IObject).GetIsActive() {
			e.Value.(IObject).Render()
		}
	}
}

// 清理
func (s *Scene) Clean() {
	s.Object.Clean()
	for e := s.ChildrenWorld.Front(); e != nil; e = e.Next() {
		e.Value.(IObject).Clean()
	}
	for e := s.ChildrenScreen.Front(); e != nil; e = e.Next() {
		e.Value.(IObject).Clean()
	}
	s.ChildrenWorld.Init()
	s.ChildrenScreen.Init()
}

// 增加孩子
func (s *Scene) AddChild(child IObject) {
	switch child.GetType() {
	case ObjectTypeWorld, ObjectTypeEnemy:
		s.ChildrenWorld.PushBack(child)
	case ObjectTypeScreen:
		s.ChildrenScreen.PushBack(child)
	default:
		s.Object.AddChild(child)
	}
}

// 移除孩子
func (s *Scene) RemoveChild(child IObject) {
	switch child.GetType() {
	case ObjectTypeWorld:
		for e := s.ChildrenWorld.Front(); e != nil; e = e.Next() {
			if e.Value.(IObject) == child {
				s.ChildrenWorld.Remove(e)
			}
		}
	case ObjectTypeScreen:
		for e := s.ChildrenScreen.Front(); e != nil; e = e.Next() {
			if e.Value.(IObject) == child {
				s.ChildrenScreen.Remove(e)
			}
		}
	default:
		s.Object.RemoveChild(child)
	}
}

// 获取世界对象孩子
func (s *Scene) GetChildWorld() *list.List {
	return &s.ChildrenWorld
}

// 获取屏幕对象孩子
func (s *Scene) GetChildScreen() *list.List {
	return &s.ChildrenScreen
}

// 暂停
func (s *Scene) Pause() {
	s.IsPause = true
	s.Game().PauseAllMusic()
	s.Game().PauseAllEffects()
}

// 恢复
func (s *Scene) Resume() {
	s.IsPause = false
	s.Game().ResumeAllMusic()
	s.Game().ResumeAllEffects()
}
