package core

// 物体类型
type ObjectType int

const (
	// 无类型
	ObjectTypeNone ObjectType = iota
	// 屏幕对象
	ObjectTypeScreen
	// 世界对象
	ObjectTypeWorld
)

// 锚点布局定义
type AnchorType int

const (
	// 无锚点
	AnchorTypeNone AnchorType = iota
	// 顶部左
	AnchorTypeTopLeft
	// 顶部中心
	AnchorTypeTopCenter
	// 顶部右
	AnchorTypeTopRight
	// 中心左
	AnchorTypeCenterLeft
	// 中心
	AnchorTypeCenter
	// 中心右
	AnchorTypeCenterRight
	// 底部左边
	AnchorTypeBottomLeft
	// 底部中心
	AnchorTypeBottomCenter
	// 底部右
	AnchorTypeBottomRight
)
