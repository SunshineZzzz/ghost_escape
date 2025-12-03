package screen

import (
	"ghost_escape/game/affiliate"
	"ghost_escape/game/core"

	"github.com/SunshineZzzz/purego-sdl3/sdl"
	"github.com/go-gl/mathgl/mgl32"
)

// 按钮HUD
type HudButton struct {
	// 继承基础屏幕对象
	core.ObjectScreen
	// 正常精灵图状态
	normalSprite *affiliate.Sprite
	// 悬停精灵图状态
	hoverSprite *affiliate.Sprite
	// 按压精灵图状态
	pressSprite *affiliate.Sprite
	// 是否悬停
	isHover bool
	// 是否按压
	isPress bool
	// 是否触发
	isTrigger bool
}

func AddHudButtonChild(parent core.IObject, renderPos mgl32.Vec2, normalPath string, hoverPath string,
	pressPath string, scale float32, anchor core.AnchorType) *HudButton {
	b := &HudButton{}
	b.Init()
	b.SetRenderPosition(renderPos)
	b.normalSprite = affiliate.AddSpriteChild(b, normalPath, scale, anchor)
	b.hoverSprite = affiliate.AddSpriteChild(b, hoverPath, scale, anchor)
	b.pressSprite = affiliate.AddSpriteChild(b, pressPath, scale, anchor)
	b.normalSprite.SetActive(true)
	b.hoverSprite.SetActive(false)
	b.pressSprite.SetActive(false)
	if parent != nil {
		parent.AddChild(b)
	}
	return b
}

// 处理事件
func (b *HudButton) HandleEvent(event *sdl.Event) {
	b.ObjectScreen.HandleEvent(event)
	if event.Type() == sdl.EventMouseButtonDown {
		if event.Button().Button == uint8(sdl.ButtonLeft) {
			if b.isHover {
				b.isPress = true
				b.Game().PlaySound("assets/sound/UI_button08.wav", false)
			}
		}
	} else if event.Type() == sdl.EventMouseButtonUp {
		if event.Button().Button == uint8(sdl.ButtonLeft) {
			b.isPress = false
			if b.isHover {
				b.isTrigger = true
			}
		}
	}
}

// 更新
func (b *HudButton) Update(dt float32) {
	b.ObjectScreen.Update(dt)
	b.checkHover()
	b.checkState()
}

// 非接口实现

// 检查悬停
func (b *HudButton) checkHover() {
	newHover := false
	pos := b.GetRenderPosition().Add(b.normalSprite.GetOffset())
	size := b.normalSprite.GetSize()
	if b.Game().IsMouseInRect(pos, pos.Add(size)) {
		newHover = true
	} else {
		newHover = false
	}
	if newHover != b.isHover {
		b.isHover = newHover
		if b.isHover && !b.isPress {
			b.Game().PlaySound("assets/sound/UI_button12.wav", false)
		}
	}
}

// 检查状态
func (b *HudButton) checkState() {
	if !b.isPress && !b.isHover {
		b.normalSprite.SetActive(true)
		b.hoverSprite.SetActive(false)
		b.pressSprite.SetActive(false)
	} else if !b.isPress && b.isHover {
		b.normalSprite.SetActive(false)
		b.hoverSprite.SetActive(true)
		b.pressSprite.SetActive(false)
	} else {
		b.normalSprite.SetActive(false)
		b.hoverSprite.SetActive(false)
		b.pressSprite.SetActive(true)
	}
}

// 是否触发
func (b *HudButton) GetIsTrigger() bool {
	if b.isTrigger {
		b.isTrigger = false
		return true
	}
	return false
}
