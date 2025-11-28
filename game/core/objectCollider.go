package core

import (
	"github.com/go-gl/mathgl/mgl32"
)

// 碰撞器类型
type ColliderType int

const (
	// 圆形碰撞器，size的X轴为直径，默认Y=X
	ColliderTypeCircle ColliderType = iota
)

// 碰撞器组件抽象
type IObjectCollider interface {
	// 碰撞器组件抽象
	IObject
	// 获取碰撞器类型
	GetColliderType() ColliderType
	// 设置碰撞器类型
	SetColliderType(ColliderType)
	// 是否发生碰撞
	IsColliding(IObjectCollider) bool
	// 获取父节点
	GetParent() IObjectScreen
	// 设置父亲节点
	SetParent(IObjectScreen)
	// 获取相对父节点偏移
	GetOffset() mgl32.Vec2
	// 设置相对父节点偏移
	SetOffset(mgl32.Vec2)
	// 获取大小
	GetSize() mgl32.Vec2
	// 设置大小
	SetSize(mgl32.Vec2)
	// 设置缩放比例
	SetScale(float32)
}

// 基础碰撞器
type ObjectCollider struct {
	// 继承基础依附对象
	ObjectAffiliate
	// 碰撞器类型
	Type ColliderType
}

var _ IObject = (*ObjectCollider)(nil)
var _ IObjectCollider = (*ObjectCollider)(nil)
var _ IObjectAffiliate = (*ObjectCollider)(nil)

// 获取碰撞器类型
func (o *ObjectCollider) GetColliderType() ColliderType {
	return o.Type
}

// 设置碰撞器类型
func (o *ObjectCollider) SetColliderType(t ColliderType) {
	o.Type = t
}

// 是否发生碰撞
func (o *ObjectCollider) IsColliding(other IObjectCollider) bool {
	panic("not implemented")
}
