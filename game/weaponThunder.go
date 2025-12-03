package game

import (
	"ghost_escape/game/core"
	"ghost_escape/game/raw"
	"ghost_escape/game/world"

	"github.com/SunshineZzzz/purego-sdl3/sdl"
)

// 雷武器组件
type WeaponThunder struct {
	// 继承基础武器组件
	raw.Weapon
}

var _ core.IObject = (*WeaponThunder)(nil)

// 创建雷武器组件
func AddWeaponThunderChild(parent *core.Actor, cooldown float32, manaCost float32) *WeaponThunder {
	w := &WeaponThunder{}
	w.Init()
	w.SetParent(parent)
	w.SetCooldown(cooldown)
	w.SetManaCost(manaCost)
	if parent != nil {
		parent.AddChild(w)
	}
	return w
}

// 处理事件
func (w *WeaponThunder) HandleEvent(event *sdl.Event) {
	w.Weapon.HandleEvent(event)
	// 处理攻击事件
	if event.Type() == sdl.EventMouseButtonDown {
		if event.Button().Button == uint8(sdl.ButtonLeft) {
			if w.CanAttack() {
				w.Game().PlaySound("assets/sound/big-thunder.mp3", false)
				pos := core.GetInstance().GetMousePosition().Add(core.GetInstance().GetCurrentScene().GetCameraPosition())
				spell := world.AddSpellChild(nil, "assets/effect/Thunderstrike w blur.png", pos, 40.0, 3.0, core.AnchorTypeCenter)
				// 攻击
				w.Attack(pos, spell)
			}
		}
	}
}

// 获取技能使用恢复百分比
func (w *WeaponThunder) GetSkillPercent() float32 {
	return w.CooldownTimer / w.Cooldown
}
