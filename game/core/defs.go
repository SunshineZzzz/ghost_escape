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
