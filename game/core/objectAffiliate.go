package core

import (
	"github.com/go-gl/mathgl/mgl32"
)

// 基础依附对象
type ObjectAffiliate struct {
	// 继承基础对象
	Object
	// 父节点
	Parent IObjectWorld
	// 相对父节点偏移
	Offset mgl32.Vec2
	// 大小
	Size mgl32.Vec2
	// 锚点布局类型
	AnchorType AnchorType
}

var _ IObject = (*ObjectAffiliate)(nil)

// 初始化
func (o *ObjectAffiliate) Init() {
	o.Object.Init()
	o.AnchorType = AnchorTypeCenter
}

// 非接口实现

// 获取父节点
func (o *ObjectAffiliate) GetParent() IObjectWorld {
	return o.Parent
}

// 设置父亲节点
func (o *ObjectAffiliate) SetParent(parent IObjectWorld) {
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
	o.SetOffsetByAnchorType(o.AnchorType)
}

// 设置缩放比例
func (o *ObjectAffiliate) SetScale(scale float32) {
	o.Size = o.Size.Mul(scale)
	o.SetOffsetByAnchorType(o.AnchorType)
}

// 获取锚点布局类型
func (o *ObjectAffiliate) GetAnchorType() AnchorType {
	return o.AnchorType
}

// 设置锚点布局类型
func (o *ObjectAffiliate) SetAnchorType(anchorType AnchorType) {
	o.AnchorType = anchorType
}

// 根据锚点类型设置偏移
func (o *ObjectAffiliate) SetOffsetByAnchorType(anchor AnchorType) {
	switch anchor {
	case AnchorTypeTopLeft:
		o.Offset = mgl32.Vec2{0, 0}
	case AnchorTypeTopCenter:
		o.Offset = mgl32.Vec2{-o.Size.X() / 2, 0}
	case AnchorTypeTopRight:
		o.Offset = mgl32.Vec2{-o.Size.X(), 0}
	case AnchorTypeCenterLeft:
		o.Offset = mgl32.Vec2{0, -o.Size.Y() / 2}
	case AnchorTypeCenter:
		o.Offset = mgl32.Vec2{-o.Size.X() / 2, -o.Size.Y() / 2}
	case AnchorTypeCenterRight:
		o.Offset = mgl32.Vec2{-o.Size.X(), -o.Size.Y() / 2}
	case AnchorTypeBottomLeft:
		o.Offset = mgl32.Vec2{0, -o.Size.Y()}
	case AnchorTypeBottomCenter:
		o.Offset = mgl32.Vec2{-o.Size.X() / 2, -o.Size.Y()}
	case AnchorTypeBottomRight:
		o.Offset = mgl32.Vec2{-o.Size.X(), -o.Size.Y()}
	default:
		o.Offset = mgl32.Vec2{0, 0}
	}
}
