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
	// 敌人
	ObjectTypeEnemy
)

// 锚点布局定义，用于确定父亲的位置
type AnchorType int

const (
	// 无锚点
	AnchorTypeNone AnchorType = iota
	// 顶部左，绘制完成后，父亲位于绘制矩形左上角
	AnchorTypeTopLeft
	// 顶部中心，绘制完成后，父亲位于绘制矩形顶部中心
	AnchorTypeTopCenter
	// 顶部右，绘制完成后，父亲位于绘制矩形右上角
	AnchorTypeTopRight
	// 中心左，绘制完成后，父亲位于绘制矩形中心左
	AnchorTypeCenterLeft
	// 中心，绘制完成后，父亲位于绘制矩形中心
	AnchorTypeCenter
	// 中心右，绘制完成后，父亲位于绘制矩形中心右
	AnchorTypeCenterRight
	// 底部左边，绘制完成后，父亲位于绘制矩形底部左
	AnchorTypeBottomLeft
	// 底部中心，绘制完成后，父亲位于绘制矩形底部中心
	AnchorTypeBottomCenter
	// 底部右，绘制完成后，父亲位于绘制矩形底部右
	AnchorTypeBottomRight
)
