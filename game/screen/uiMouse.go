package screen

import (
	"ghost_escape/game/affiliate"
	"ghost_escape/game/core"
)

// 鼠标UI
type UIMouse struct {
	// 继承基础屏幕对象
	core.ObjectScreen
	// 鼠标精灵图1
	sprite1 *affiliate.Sprite
	// 鼠标精灵图2
	sprite2 *affiliate.Sprite
	// 图片切换计时器
	timer float32
}

var _ core.IObject = (*UIMouse)(nil)
var _ core.IObjectScreen = (*UIMouse)(nil)

func AddUIMouse(parent core.IObject, spritePath1, spritePath2 string, scale float32, anchor core.AnchorType) *UIMouse {
	uiMouse := &UIMouse{}
	uiMouse.Init()
	uiMouse.sprite1 = affiliate.AddSpriteChild(uiMouse, spritePath1, scale, anchor)
	uiMouse.sprite2 = affiliate.AddSpriteChild(uiMouse, spritePath2, scale, anchor)
	if parent != nil {
		parent.AddChild(uiMouse)
	}
	return uiMouse
}

// 更新
func (u *UIMouse) Update(dt float32) {
	u.ObjectScreen.Update(dt)
	u.timer += dt
	if u.timer < 0.3 {
		u.sprite1.SetActive(true)
		u.sprite2.SetActive(false)
	} else if u.timer < 0.6 {
		u.sprite1.SetActive(false)
		u.sprite2.SetActive(true)
	} else {
		u.timer = 0
	}
	u.SetRenderPosition(core.GetInstance().GetMousePosition())
}

// 非阶接口实现

// 获取鼠标精灵图1
func (u *UIMouse) GetSprite1() *affiliate.Sprite {
	return u.sprite1
}

// 设置鼠标精灵图1
func (u *UIMouse) SetSprite1(un *affiliate.Sprite) {
	u.sprite1 = un
}

// 获取鼠标精灵图2
func (o *UIMouse) GetSprite2() *affiliate.Sprite {
	return o.sprite2
}

// 设置鼠标精灵图2
func (u *UIMouse) SetSprite2(un *affiliate.Sprite) {
	u.sprite2 = un
}
