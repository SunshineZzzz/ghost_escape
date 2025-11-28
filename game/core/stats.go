package core

// 基础状态组件
type Stats struct {
	// 继承基础对象
	Object
	// 父节点
	Parent *Actor
	// 血量
	Health float32
	// 最大血量
	MaxHealth float32
	// 法力
	Mana float32
	// 法力重新恢复速度
	ManaRegenSpeed float32
	// 最大法力
	MaxMana float32
	// 伤害
	Damage float32
	// 受伤后无敌时间有多长
	InvincibleTime float32
	// 无敌计时器
	InvincibleTimer float32
	// 是否活着
	IsAlive bool
	// 是否无敌
	IsInvincible bool
}

var _ IObject = (*Stats)(nil)

func AddStatusChild(parent *Actor, maxHealth, maxMana, damage, manaRegenSpeed float32) *Stats {
	stats := &Stats{}
	stats.Init()
	stats.Parent = parent
	stats.MaxHealth = maxHealth
	stats.Health = maxHealth
	stats.MaxMana = maxMana
	stats.Mana = maxMana
	stats.ManaRegenSpeed = manaRegenSpeed
	stats.Damage = damage
	if parent != nil {
		parent.AddChild(stats)
	}
	return stats
}

// 初始化
func (s *Stats) Init() {
	s.Object.Init()
	s.IsAlive = true
	s.InvincibleTime = 1.5
	s.InvincibleTimer = 0.0
	s.IsInvincible = false
}

// 更新状态
func (s *Stats) Update(dt float32) {
	s.Object.Update(dt)
	s.regenMana(dt)
	if s.IsInvincible {
		s.InvincibleTimer += dt
		if s.InvincibleTimer >= s.InvincibleTime {
			s.IsInvincible = false
			s.InvincibleTimer = 0.0
		}
	}
}

// 非接口实现

// 能否使用法力
func (s *Stats) CanUseMana(manaCost float32) bool {
	return s.Mana >= manaCost
}

// 使用法力
func (s *Stats) UseMana(manaCost float32) {
	s.Mana -= manaCost
	if s.Mana < 0.0 {
		s.Mana = 0.0
	}
}

// 恢复法力
func (s *Stats) regenMana(dt float32) {
	s.Mana += s.ManaRegenSpeed * dt
	if s.Mana > s.MaxMana {
		s.Mana = s.MaxMana
	}
}

// 被伤害
func (s *Stats) TakeDamage(damage float32) {
	if s.IsInvincible {
		return
	}

	s.Health -= damage
	if s.Health < 0.0 {
		s.Health = 0.0
		s.IsAlive = false
	}
	// fmt.Printf("damage: %f, health: %f\n", damage, s.Health)
	s.IsInvincible = true
	s.InvincibleTimer = 0.0
}

// 获取生命值
func (s *Stats) GetHealth() float32 {
	return s.Health
}

// 设置生命值
func (s *Stats) SetHealth(health float32) {
	s.Health = health
}

// 获取最大生命值
func (s *Stats) GetMaxHealth() float32 {
	return s.MaxHealth
}

// 设置最大生命值
func (s *Stats) SetMaxHealth(maxHealth float32) {
	s.MaxHealth = maxHealth
}

// 获取法力
func (s *Stats) GetMana() float32 {
	return s.Mana
}

// 设置法力
func (s *Stats) SetMana(mana float32) {
	s.Mana = mana
}

// 获取最大法力
func (s *Stats) GetMaxMana() float32 {
	return s.MaxMana
}

// 设置最大法力
func (s *Stats) SetMaxMana(maxMana float32) {
	s.MaxMana = maxMana
}

// 获取伤害
func (s *Stats) GetDamage() float32 {
	return s.Damage
}

// 设置伤害
func (s *Stats) SetDamage(damage float32) {
	s.Damage = damage
}

// 获取法力恢复速度
func (s *Stats) GetManaRegenSpeed() float32 {
	return s.ManaRegenSpeed
}

// 设置法力恢复速度
func (s *Stats) SetManaRegenSpeed(manaRegenSpeed float32) {
	s.ManaRegenSpeed = manaRegenSpeed
}

// 获取是否活着
func (s *Stats) GetAlive() bool {
	return s.IsAlive
}

// 设置是否活着
func (s *Stats) SetAlive(alive bool) {
	s.IsAlive = alive
}

// 设置父亲节点
func (s *Stats) SetParent(parent *Actor) {
	s.Parent = parent
}

// 获取父亲节点
func (s *Stats) GetParent() *Actor {
	return s.Parent
}

// 设置无敌
func (s *Stats) SetInvincible(invincible bool) {
	s.IsInvincible = invincible
}

// 获取是否无敌
func (s *Stats) GetInvincible() bool {
	return s.IsInvincible
}
