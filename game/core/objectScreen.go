package core

import (
	"github.com/go-gl/mathgl/mgl32"
)

// 屏幕对象抽象
type IObjectScreen interface {
	// 继承基础对象接口
	IObject
	// 设置渲染(屏幕)位置
	SetRenderPosition(mgl32.Vec2)
	// 获取渲染(屏幕)位置
	GetRenderPosition() mgl32.Vec2
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

// 初始化
func (o *ObjectScreen) Init() {
	o.Object.Init()
	o.ObjectType = ObjectTypeScreen
}

// 获取渲染(屏幕)位置
func (o *ObjectScreen) GetRenderPosition() mgl32.Vec2 {
	return o.RenderPosition
}

// 非接口实现
