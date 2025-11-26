package core

import (
	"github.com/go-gl/mathgl/mgl32"
)

// 基础角色
type Actor struct {
	// 继承基础世界对象
	ObjectWorld
	// 速度，X/Y每秒移动像素
	Velocity mgl32.Vec2
	// 最大速度大小
	MaxSpeed float32
}

var _ IObject = (*Actor)(nil)
var _ IObjectScreen = (*Actor)(nil)

// 非接口实现

// 设置速度
func (a *Actor) SetVelocity(v mgl32.Vec2) {
	a.Velocity = v
}

// 获取速度
func (a *Actor) GetVelocity() mgl32.Vec2 {
	return a.Velocity
}

// 设置最大速度大小
func (a *Actor) SetMaxSpeed(speed float32) {
	a.MaxSpeed = speed
}

// 获取最大速度大小
func (a *Actor) GetMaxSpeed() float32 {
	return a.MaxSpeed
}
