package raw

import (
	"ghost_escape/game/core"
	"ghost_escape/game/world"

	"github.com/go-gl/mathgl/mgl32"
)

// 基础武器组件
type Weapon struct {
	// 继承基础对象
	core.Object
	// 父节点
	Parent *core.Actor
	// 冷却时间
	Cooldown float32
	// 冷却计时器
	CooldownTimer float32
	// 法力消耗
	ManaCost float32
}

var _ core.IObject = (*Weapon)(nil)

// 初始化
func (w *Weapon) Init() {
	w.Object.Init()
	w.CooldownTimer = 0.0
	w.ManaCost = 0.0
	w.Cooldown = 1.0
}

// 更新
func (w *Weapon) Update(dt float32) {
	w.Object.Update(dt)
	// 很可能大于1
	w.CooldownTimer += dt
}

// 非接口实现

// 是否可以攻击
func (w *Weapon) CanAttack() bool {
	if w.CooldownTimer < w.Cooldown {
		return false
	}
	if !w.Parent.GetStats().CanUseMana(w.ManaCost) {
		return false
	}
	return true
}

// 攻击
func (w *Weapon) Attack(pos mgl32.Vec2, spell *world.Spell) {
	if spell == nil {
		return
	}
	w.Parent.GetStats().UseMana(w.ManaCost)
	w.CooldownTimer = 0.0
	spell.SetPosition(pos)
	core.GetInstance().GetCurrentScene().AddChild(spell)
}

// 设置父节点
func (w *Weapon) SetParent(parent *core.Actor) {
	w.Parent = parent
}

// 获取父节点
func (w *Weapon) GetParent() *core.Actor {
	return w.Parent
}

// 获取法力消耗
func (w *Weapon) GetManaCost() float32 {
	return w.ManaCost
}

// 设置法力消耗
func (w *Weapon) SetManaCost(manaCost float32) {
	w.ManaCost = manaCost
}

// 设置冷却时间
func (w *Weapon) SetCooldown(cooldown float32) {
	w.Cooldown = cooldown
}

// 获取冷却时间
func (w *Weapon) GetCooldown() float32 {
	return w.Cooldown
}
