package core

import (
	"github.com/go-gl/mathgl/mgl32"
)

// 基础世界对象
type ObjectWorld struct {
	// 继承基础屏幕对象
	ObjectScreen
	// 世界位置
	Position mgl32.Vec2
}

var _ IObject = (*ObjectWorld)(nil)
var _ IObjectScreen = (*ObjectWorld)(nil)

// 更新
func (o *ObjectWorld) Update(dt float32) {
	o.ObjectScreen.Update(dt)
	o.RenderPosition = o.Game().GetCurrentScene().WorldToScreen(o.Position)
}

// 设置渲染(屏幕)位置
func (o *ObjectWorld) SetRenderPosition(pos mgl32.Vec2) {
	o.RenderPosition = pos
	o.Position = o.Game().GetCurrentScene().ScreenToWorld(pos)
}

// 初始化
func (o *ObjectWorld) Init() {
	o.ObjectScreen.Init()
	o.ObjectType = ObjectTypeWorld
}

// 非接口实现

// 获取世界位置
func (o *ObjectWorld) GetPosition() mgl32.Vec2 {
	return o.Position
}

// 设置世界位置
func (o *ObjectWorld) SetPosition(pos mgl32.Vec2) {
	o.Position = pos
	o.RenderPosition = o.Game().GetCurrentScene().WorldToScreen(pos)
}
