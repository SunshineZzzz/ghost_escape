package affiliate

import (
	"ghost_escape/game/core"

	"github.com/SunshineZzzz/purego-sdl3/ttf"
	"github.com/go-gl/mathgl/mgl32"
)

// 文本标签组件
type TextLabel struct {
	// 继承基础依附对象
	core.ObjectAffiliate
	// 文本对象
	ttfText *ttf.Text
	// 字体路径
	fontPath string
	// 字体大小
	fontSize float32
}

// 创建文本标签组件
func AddTextLabelChild(parent core.IObjectScreen, text string, fontPath string, fontSize float32, anchor core.AnchorType) *TextLabel {
	label := &TextLabel{}
	label.Init()
	label.SetFont(fontPath, fontSize)
	label.SetText(text)
	label.SetAnchorType(anchor)
	label.updateSize()
	if parent != nil {
		label.SetParent(parent)
		parent.AddChild(label)
	}
	return label
}

// 清除
func (t *TextLabel) Clear() {
	if t.ttfText != nil {
		ttf.DestroyText(t.ttfText)
		t.ttfText = nil
	}
}

// 渲染
func (t *TextLabel) Render() {
	t.ObjectAffiliate.Render()
	if t.ttfText == nil {
		return
	}
	pos := t.Parent.GetRenderPosition().Add(t.Offset)
	ttf.DrawRendererText(t.ttfText, pos.X(), pos.Y())
}

// 非接口实现

// 设置字体
func (t *TextLabel) SetFont(fontPath string, fontSize float32) {
	t.fontPath = fontPath
	t.fontSize = fontSize
	font, err := t.Game().GetAssetStore().GetFont(fontPath, fontSize)
	if err != nil {
		return
	}
	if t.ttfText != nil {
		t.Clean()
	}
	t.ttfText = t.Game().CreateTTFText("", fontPath, fontSize)
	ttf.SetTextFont(t.ttfText, font)
}

// 设置字体路径
func (t *TextLabel) SetFontPath(fontPath string) {
	t.Clear()
	t.SetFont(fontPath, t.fontSize)
}

// 设置字体大小
func (t *TextLabel) SetFontSize(fontSize float32) {
	t.Clear()
	t.SetFont(t.fontPath, fontSize)
}

// 设置文本
func (t *TextLabel) SetText(text string) {
	if t.ttfText == nil {
		return
	}
	ttf.SetTextString(t.ttfText, text, uint64(len(text)))
}

// 更新文本大小
func (t *TextLabel) updateSize() {
	if t.ttfText == nil {
		return
	}
	var w, h int32
	ttf.GetTextSize(t.ttfText, &w, &h)
	t.SetSize(mgl32.Vec2{float32(w), float32(h)})
}
