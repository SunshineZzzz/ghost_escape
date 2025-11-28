package world

import (
	"ghost_escape/game/affiliate"
	"ghost_escape/game/core"

	"github.com/go-gl/mathgl/mgl32"
)

// 法术
type Spell struct {
	// 继承基础世界对象
	core.ObjectWorld
	// 精灵动画
	spriteAnim core.IObjectAnima
	// 伤害值
	damage float32
}

var _ core.IObject = (*Spell)(nil)
var _ core.IObjectScreen = (*Spell)(nil)
var _ core.IObjectWorld = (*Spell)(nil)

// 创建法术
func AddSpellChild(parent core.IObject, filePath string, pos mgl32.Vec2, damage, scale float32, anchor core.AnchorType) *Spell {
	spell := &Spell{}
	spell.Init()
	spell.damage = damage
	spell.spriteAnim = affiliate.AddSpriteAnimChild(spell, filePath, scale, anchor)
	spell.spriteAnim.SetLoop(false)
	size := spell.spriteAnim.GetSize()
	spell.Collider = affiliate.AddColliderChild(spell, size, core.ColliderTypeCircle, anchor)
	spell.SetPosition(pos)
	if parent != nil {
		parent.AddChild(spell)
	}
	return spell
}

// 更新
func (s *Spell) Update(dt float32) {
	s.ObjectWorld.Update(dt)
	if s.spriteAnim.GetFinish() {
		s.NeedRemove = true
	}
	s.attack()
}

// 非接口实现

// 攻击
func (s *Spell) attack() {
	objects := core.GetInstance().GetCurrentScene().GetChildWorld()
	// 遍历所有对象
	for e := objects.Front(); e != nil; e = e.Next() {
		object := e.Value.(core.IObjectWorld)
		if object.GetType() != core.ObjectTypeEnemy {
			continue
		}
		// 检查碰撞
		if s.Collider.IsColliding(object.GetCollider()) {
			// 敌人受到伤害
			object.TakeDamage(s.damage)
		}
	}
}
