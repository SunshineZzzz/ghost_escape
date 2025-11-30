package screen

import (
	"ghost_escape/game/affiliate"
	"ghost_escape/game/core"

	"github.com/go-gl/mathgl/mgl32"
)

// 状态HUD
type HudStats struct {
	// 继承基础屏幕对象
	core.ObjectScreen
	// 目标角色
	target *core.Actor
	// 血条图标
	healthBarIcon affiliate.ISprite
	// 血条精灵图
	healthBar affiliate.ISprite
	// 血条背景精灵图
	healthBarBg affiliate.ISprite
	// 法力条图标
	manaBarIcon affiliate.ISprite
	// 法力条精灵图
	manaBar affiliate.ISprite
	// 法力条背景精灵图
	manaBarBg affiliate.ISprite
	// 血量百分比
	healthPercent float32
	// 法力百分比
	manaPercent float32

	// 更加详细的设置大小和位置
}

// 创建状态HUD
func AddHudStats(parent core.IObject, target *core.Actor, renderPosition mgl32.Vec2) *HudStats {
	// 创建状态HUD对象
	hudStats := &HudStats{}
	hudStats.Init()
	hudStats.SetRenderPosition(renderPosition)
	hudStats.SetTarget(target)
	if parent != nil {
		parent.AddChild(hudStats)
	}
	return hudStats
}

// 初始化
func (h *HudStats) Init() {
	h.ObjectScreen.Init()
	h.healthBarBg = affiliate.AddSpriteChild(h, "assets/UI/bar_bg.png", 3.0, core.AnchorTypeCenterLeft)
	h.healthBarBg.SetOffset(h.healthBarBg.GetOffset().Add(mgl32.Vec2{30.0, 0.0}))
	h.healthBar = affiliate.AddSpriteChild(h, "assets/UI/bar_red.png", 3.0, core.AnchorTypeCenterLeft)
	h.healthBar.SetOffset(h.healthBar.GetOffset().Add(mgl32.Vec2{30.0, 0.0}))
	h.healthBarIcon = affiliate.AddSpriteChild(h, "assets/UI/Red Potion.png", 0.5, core.AnchorTypeCenterLeft)

	h.manaBarBg = affiliate.AddSpriteChild(h, "assets/UI/bar_bg.png", 3.0, core.AnchorTypeCenterLeft)
	h.manaBarBg.SetOffset(h.manaBarBg.GetOffset().Add(mgl32.Vec2{300.0, 0.0}))
	h.manaBar = affiliate.AddSpriteChild(h, "assets/UI/bar_blue.png", 3.0, core.AnchorTypeCenterLeft)
	h.manaBar.SetOffset(h.manaBar.GetOffset().Add(mgl32.Vec2{300.0, 0.0}))
	h.manaBarIcon = affiliate.AddSpriteChild(h, "assets/UI/Blue Potion.png", 0.5, core.AnchorTypeCenterLeft)
	h.manaBarIcon.SetOffset(h.manaBarIcon.GetOffset().Add(mgl32.Vec2{270.0, 0.0}))
}

// 更新
func (h *HudStats) Update(dt float32) {
	h.ObjectScreen.Update(dt)
	h.updateHealthBar()
	h.updateManaBar()
}

// 非接口实现

// 设置目标
func (h *HudStats) SetTarget(target *core.Actor) {
	h.target = target
}

// 获取目标
func (h *HudStats) GetTarget() *core.Actor {
	return h.target
}

// 设置血量百分比
func (h *HudStats) SetHealthPercent(healthPercent float32) {
	h.healthPercent = healthPercent
}

// 获取血量百分比
func (h *HudStats) GetHealthPercent() float32 {
	return h.healthPercent
}

// 设置法力百分比
func (h *HudStats) SetManaPercent(manaPercent float32) {
	h.manaPercent = manaPercent
}

// 获取法力百分比
func (h *HudStats) GetManaPercent() float32 {
	return h.manaPercent
}

// 设置血条精灵图
func (h *HudStats) SetHealthBar(healthBar affiliate.ISprite) {
	h.healthBar = healthBar
}

// 获取血条精灵图
func (h *HudStats) GetHealthBar() affiliate.ISprite {
	return h.healthBar
}

// 设置血条背景精灵图
func (h *HudStats) SetHealthBarBg(healthBarBg affiliate.ISprite) {
	h.healthBarBg = healthBarBg
}

// 获取血条背景精灵图
func (h *HudStats) GetHealthBarBg() affiliate.ISprite {
	return h.healthBarBg
}

// 更新血条
func (h *HudStats) updateHealthBar() {
	if h.target == nil || h.healthBarBg == nil || h.target.GetStats() == nil {
		return
	}
	h.healthBar.SetPercent(mgl32.Vec2{h.target.GetStats().GetHealth() / h.target.GetStats().GetMaxHealth(), 1.0})
}

// 更新法力条
func (h *HudStats) updateManaBar() {
	if h.target == nil || h.manaBarBg == nil || h.target.GetStats() == nil {
		return
	}
	h.manaBar.SetPercent(mgl32.Vec2{h.target.GetStats().GetMana() / h.target.GetStats().GetMaxMana(), 1.0})
}
