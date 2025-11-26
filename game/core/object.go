package core

import (
	"github.com/SunshineZzzz/purego-sdl3/sdl"
)

// 对象抽象
type IObject interface {
	// 初始化
	Init()
	// 处理事件
	HandleEvent(*sdl.Event)
	// 更新
	Update(dt float32)
	// 渲染
	Render()
	// 清理
	Clean()
}

// 基础对象
type Object struct {
}

var _ IObject = (*Object)(nil)

// 初始化
func (o *Object) Init() {
}

// 处理事件
func (o *Object) HandleEvent(*sdl.Event) {
}

// 更新
func (o *Object) Update(float32) {
}

// 渲染
func (o *Object) Render() {
}

// 清理
func (o *Object) Clean() {
}

// 非接口实现

// 获取游戏实例
func (o *Object) Game() *Game {
	return GetInstance()
}
