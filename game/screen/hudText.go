package screen

import (
	"ghost_escape/game/affiliate"
	"ghost_escape/game/core"

	"github.com/go-gl/mathgl/mgl32"
)

// HUD文本
type HudText struct {
	// 继承基础屏幕对象
	core.ObjectScreen
	// 文本标签
	textLabel *affiliate.TextLabel
	// 背景图
	spriteBG *affiliate.Sprite
	// 背景图大小
	bgSize mgl32.Vec2
}

// 创建HUD文本
func AddHudTextChild(parent core.IObject, text string, renderPos mgl32.Vec2, bgSize mgl32.Vec2, fontPath string,
	fontSize float32, bgPath string, anchor core.AnchorType) *HudText {

	h := &HudText{}
	h.Init()
	h.SetRenderPosition(renderPos)
	h.setSpriteBG(affiliate.AddSpriteChild(h, bgPath, 1.0, anchor))
	h.setBgSize(bgSize)
	h.setTextLabel(affiliate.AddTextLabelChild(h, text, fontPath, fontSize, anchor))
	if parent != nil {
		parent.AddChild(h)
	}
	return h
}

// 设置背景图
func (h *HudText) setSpriteBG(sprite *affiliate.Sprite) {
	h.spriteBG = sprite
}

// 设置背景图大小
func (h *HudText) setBgSize(size mgl32.Vec2) {
	h.bgSize = size
	h.spriteBG.SetSize(size)
}

// 设置文本标签
func (h *HudText) setTextLabel(label *affiliate.TextLabel) {
	h.textLabel = label
}

// 设置文本
func (h *HudText) SetText(text string) {
	if h.textLabel == nil {
		return
	}
	h.textLabel.SetText(text)
}

// 根据文本设置背景大小
func (h *HudText) SetBgSizeByText(margin float32) {
	textSize := h.textLabel.GetSize()
	h.setBgSize(mgl32.Vec2{textSize.X() + margin, textSize.Y() + margin})
}
