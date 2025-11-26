package core

import (
	"github.com/go-gl/mathgl/mgl32"
)

// 基础依附对象
type ObjectAffiliate struct {
	// 继承基础对象
	Object
	// 父节点
	Parent IObjectScreen
	// 相对父节点偏移
	Offset mgl32.Vec2
	// 大小
	Size mgl32.Vec2
}

var _ IObject = (*ObjectAffiliate)(nil)

// 初始化
func (o *ObjectAffiliate) Init() {
	o.Object.Init()
}

// 非接口实现

// 获取父节点
func (o *ObjectAffiliate) GetParent() IObjectScreen {
	return o.Parent
}

// 设置父亲节点
func (o *ObjectAffiliate) SetParent(parent IObjectScreen) {
	o.Parent = parent
}

// 获取相对父节点偏移
func (o *ObjectAffiliate) GetOffset() mgl32.Vec2 {
	return o.Offset
}

// 设置相对父节点偏移
func (o *ObjectAffiliate) SetOffset(offset mgl32.Vec2) {
	o.Offset = offset
}

// 获取大小
func (o *ObjectAffiliate) GetSize() mgl32.Vec2 {
	return o.Size
}

// 设置大小
func (o *ObjectAffiliate) SetSize(size mgl32.Vec2) {
	o.Size = size
}
