package affiliate

import (
	"ghost_escape/game/core"

	"github.com/SunshineZzzz/purego-sdl3/sdl"
	"github.com/go-gl/mathgl/mgl32"
)

// 进度条依附对象
type AffiliateBar struct {
	// 继承基础依附对象
	core.ObjectAffiliate
	// 百分比
	Percent mgl32.Vec2
	// 高进度条颜色
	ColorHigh sdl.FColor
	// 中等进度条颜色
	ColorMid sdl.FColor
	// 低进度条颜色
	ColorLow sdl.FColor
}

var _ core.IObject = (*AffiliateBar)(nil)
var _ core.IObjectAffiliate = (*AffiliateBar)(nil)

// 创建进度条依附对象
func AddAffiliateBarChild(parent core.IObjectScreen, size mgl32.Vec2, anchor core.AnchorType) *AffiliateBar {
	bar := &AffiliateBar{}
	bar.Init()
	bar.SetAnchorType(anchor)
	bar.SetSize(size)
	if parent != nil {
		bar.SetParent(parent)
		parent.AddChild(bar)
	}
	return bar
}

// 初始化进度条
func (b *AffiliateBar) Init() {
	b.ObjectAffiliate.Init()
	b.Percent = mgl32.Vec2{0.0, 1.0}
	// 绿色
	b.ColorHigh = sdl.FColor{R: 0.0, G: 1.0, B: 0.0, A: 1.0}
	// 橙色
	b.ColorMid = sdl.FColor{R: 1.0, G: 0.65, B: 0.0, A: 1.0}
	// 红色
	b.ColorLow = sdl.FColor{R: 1.0, G: 0.0, B: 0.0, A: 1.0}
}

// 渲染进度条
func (b *AffiliateBar) Render() {
	b.ObjectAffiliate.Render()
	pos := b.Parent.GetRenderPosition().Add(b.Offset)
	if b.Percent.X() > 0.75 {
		core.GetInstance().RenderHBar(pos, b.Size, b.Percent, b.ColorHigh)
	} else if b.Percent.X() > 0.3 {
		core.GetInstance().RenderHBar(pos, b.Size, b.Percent, b.ColorMid)
	} else {
		core.GetInstance().RenderHBar(pos, b.Size, b.Percent, b.ColorLow)
	}
}

// 非接口实现

// 获取进度条百分比
func (b *AffiliateBar) GetPercent() mgl32.Vec2 {
	return b.Percent
}

// 设置进度条百分比
func (b *AffiliateBar) SetPercent(percent mgl32.Vec2) {
	b.Percent = percent
}

// 设置高进度条颜色
func (b *AffiliateBar) SetColorHigh(color sdl.FColor) {
	b.ColorHigh = color
}

// 获取高进度条颜色
func (b *AffiliateBar) GetColorHigh() sdl.FColor {
	return b.ColorHigh
}

// 设置中等进度条颜色
func (b *AffiliateBar) SetColorMid(color sdl.FColor) {
	b.ColorMid = color
}

// 获取中等进度条颜色
func (b *AffiliateBar) GetColorMid() sdl.FColor {
	return b.ColorMid
}

// 设置低进度条颜色
func (b *AffiliateBar) SetColorLow(color sdl.FColor) {
	b.ColorLow = color
}

// 获取低进度条颜色
func (b *AffiliateBar) GetColorLow() sdl.FColor {
	return b.ColorLow
}
