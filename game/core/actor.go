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
	// 角色属性组件
	Stats *Stats
}

var _ IObject = (*Actor)(nil)
var _ IObjectScreen = (*Actor)(nil)

// 初始化
func (a *Actor) Init() {
	a.ObjectWorld.Init()
	a.MaxSpeed = 100.0
}

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

// 移动
func (a *Actor) Move(dt float32) {
	newPos := a.Position.Add(a.Velocity.Mul(dt))
	a.Position[0] = mgl32.Clamp(newPos.X(), 0.0, a.Game().GetWorldSize().X())
	a.Position[1] = mgl32.Clamp(newPos.Y(), 0.0, a.Game().GetWorldSize().Y())
	a.SetPosition(a.Position)
}

// 获取角色属性
func (a *Actor) GetStats() *Stats {
	return a.Stats
}

// 设置角色属性
func (a *Actor) SetStats(stats *Stats) {
	a.Stats = stats
}

// 被伤害
func (a *Actor) TakeDamage(damage float32) {
	if a.Stats == nil {
		return
	}
	a.Stats.TakeDamage(damage)
}

// 获取是否活着
func (a *Actor) GetAlive() bool {
	if a.Stats == nil {
		return false
	}
	return a.Stats.GetAlive()
}
