package core

import (
	"github.com/go-gl/mathgl/mgl32"
)

// 屏幕对象抽象
type IObjectScreen interface {
	IObject
	// 设置渲染(屏幕)位置
	SetRenderPosition(mgl32.Vec2)
}

// 基础屏幕对象
type ObjectScreen struct {
	// 继承基础对象
	Object
	// 渲染(屏幕)位置
	RenderPosition mgl32.Vec2
}

var _ IObject = (*ObjectScreen)(nil)
var _ IObjectScreen = (*ObjectScreen)(nil)

// 设置渲染(屏幕)位置
func (o *ObjectScreen) SetRenderPosition(pos mgl32.Vec2) {
	o.RenderPosition = pos
}

// 非接口实现

// 获取渲染(屏幕)位置
func (o *ObjectScreen) GetRenderPosition() mgl32.Vec2 {
	return o.RenderPosition
}
