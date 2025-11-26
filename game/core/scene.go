package core

import (
	"github.com/go-gl/mathgl/mgl32"
)

// 场景抽象
type IScene interface {
	IObject
	// 世界坐标转换为屏幕坐标
	WorldToScreen(mgl32.Vec2) mgl32.Vec2
	// 屏幕坐标转换为世界坐标
	ScreenToWorld(mgl32.Vec2) mgl32.Vec2
	// 设置摄像机位置
	SetCameraPosition(mgl32.Vec2)
	// 获取世界大小
	GetWorldSize() mgl32.Vec2
}

// 基础场景
type Scene struct {
	// 继承基础对象
	Object
	// 世界大小
	WorldSize mgl32.Vec2
	// 摄像机位置
	CameraPositon mgl32.Vec2
	// 场景中的物体
	Objects []IObject
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

// 设置摄像机位置(世界坐标系)
func (s *Scene) SetCameraPosition(pos mgl32.Vec2) {
	s.CameraPositon[0] = mgl32.Clamp(pos.X(), -30.0, s.WorldSize.X()-s.Game().GetScreenSize().X()+30.0)
	s.CameraPositon[1] = mgl32.Clamp(pos.Y(), -30.0, s.WorldSize.Y()-s.Game().GetScreenSize().Y()+30.0)
}

// 获取世界大小
func (s *Scene) GetWorldSize() mgl32.Vec2 {
	return s.WorldSize
}
