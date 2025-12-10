package game

import (
	"ghost_escape/game/affiliate"
	"ghost_escape/game/core"
	"ghost_escape/game/world"

	"github.com/SunshineZzzz/purego-sdl3/sdl"
	"github.com/go-gl/mathgl/mgl32"
)

// 玩家
type Player struct {
	// 继承基础角色
	core.Actor
	// 雷武器组件
	Weapon *WeaponThunder
	// 空闲精灵动画
	spriteIdleAnim *affiliate.SpriteAnim
	// 移动精灵动画
	spriteMoveAnim *affiliate.SpriteAnim
	// 是否在运动
	isMoving bool
	// 受伤闪烁timer
	flashTimer *core.Timer
	// 死亡特效
	deadEffect *world.Effect
}

var _ core.IObject = (*Player)(nil)
var _ core.IObjectScreen = (*Player)(nil)
var _ core.IActor = (*Player)(nil)

// 初始化
func (p *Player) Init() {
	p.Actor.Init()
	p.MaxSpeed = 500.0
	p.spriteIdleAnim = affiliate.AddSpriteAnimChild(p, "assets/sprite/ghost-idle.png", 2.0, core.AnchorTypeCenter)
	p.spriteMoveAnim = affiliate.AddSpriteAnimChild(p, "assets/sprite/ghost-move.png", 2.0, core.AnchorTypeCenter)
	p.spriteIdleAnim.SetActive(true)
	p.spriteMoveAnim.SetActive(false)
	p.isMoving = false
	p.Collider = affiliate.AddColliderChild(p, p.spriteIdleAnim.GetSize().Mul(0.5), core.ColliderTypeCircle, core.AnchorTypeCenter)
	p.Stats = core.AddStatusChild(&p.Actor, 100.0, 100.0, 40.0, 10.0)
	// 雷武器组件
	p.Weapon = AddWeaponThunderChild(&p.Actor, 2.0, 40.0)
	// affiliate.AddTextLabelChild(p, "这是主角", "assets/font/VonwaonBitmap-16px.ttf", 16.0, core.AnchorTypeCenter)

	p.deadEffect = world.AddEffectChild(nil, "assets/effect/1764.png", mgl32.Vec2{0.0, 0.0}, 2.0, core.AnchorTypeCenter, nil)

	// 受伤闪烁timer
	p.flashTimer = core.AddTimerChild(p, 0.4)
	p.flashTimer.Start()
}

// 处理事件
func (p *Player) HandleEvent(event *sdl.Event) {
	p.Actor.HandleEvent(event)
}

// 更新
func (p *Player) Update(dt float32) {
	p.Actor.Update(dt)
	// 速度慢慢减速
	p.Velocity = p.Velocity.Mul(0.9)
	p.keyboardControl()
	p.Move(dt)
	p.checkState()
	p.syncCamera()
	p.checkIsDead()
}

// 渲染
func (p *Player) Render() {
	if p.Stats.GetInvincible() && p.flashTimer.GetProcess() < 0.5 {
		return
	}
	p.Actor.Render()
}

// 清理
func (p *Player) Clean() {
	p.Actor.Clean()
}

// 非接口实现

// 键盘控制
func (p *Player) keyboardControl() {
	currentKeyStates := sdl.GetKeyboardState()
	if currentKeyStates[sdl.ScancodeW] {
		p.Velocity[1] = -p.MaxSpeed
	}
	if currentKeyStates[sdl.ScancodeS] {
		p.Velocity[1] = p.MaxSpeed
	}
	if currentKeyStates[sdl.ScancodeA] {
		p.Velocity[0] = -p.MaxSpeed
	}
	if currentKeyStates[sdl.ScancodeD] {
		p.Velocity[0] = p.MaxSpeed
	}
}

// 同步相机
func (p *Player) syncCamera() {
	// 相机跟着玩家一起移动
	p.Game().GetCurrentScene().SetCameraPosition(p.Position.Sub(p.Game().GetScreenSize().Mul(0.5)))
}

// 检查状态
func (p *Player) checkState() {
	if p.Velocity.X() < 0.0 {
		p.spriteIdleAnim.SetFlip(true)
		p.spriteMoveAnim.SetFlip(true)
	} else {
		p.spriteIdleAnim.SetFlip(false)
		p.spriteMoveAnim.SetFlip(false)
	}
	newIsMoving := p.Velocity.Len() > 0.1
	if newIsMoving != p.isMoving {
		p.isMoving = newIsMoving
		p.changeState()
	}
}

// 切换状态
func (p *Player) changeState() {
	if p.isMoving {
		p.spriteIdleAnim.SetActive(false)
		p.spriteMoveAnim.SetActive(true)
		p.spriteMoveAnim.SetFrameTimer(p.spriteIdleAnim.GetCurrentFrame())
		p.spriteMoveAnim.SetCurrentFrame(p.spriteIdleAnim.GetCurrentFrame())
	} else {
		p.spriteIdleAnim.SetActive(true)
		p.spriteMoveAnim.SetActive(false)
		p.spriteIdleAnim.SetFrameTimer(p.spriteMoveAnim.GetCurrentFrame())
		p.spriteIdleAnim.SetCurrentFrame(p.spriteMoveAnim.GetCurrentFrame())
	}
}

// 获取技能使用恢复百分比
func (p *Player) GetSkillPercent() float32 {
	if p.Weapon != nil {
		return p.Weapon.GetSkillPercent()
	}
	return 1.0
}

// 被伤害
func (p *Player) TakeDamage(damage float32) {
	if p.Stats == nil || p.Stats.GetInvincible() {
		return
	}
	p.Actor.TakeDamage(damage)
	// fmt.Printf("玩家受到伤害：%f\n", damage)
	p.Game().PlaySound("assets/sound/hit-flesh-02-266309.mp3", false)
}

// 检查是否死亡
func (p *Player) checkIsDead() {
	if !p.Stats.GetAlive() {
		// 玩家死亡
		p.Game().GetCurrentScene().SafeAddChild(p.deadEffect)
		p.deadEffect.SetPosition(p.GetPosition())
		p.SetActive(false)
		p.Game().PlaySound("assets/sound/female-scream-02-89290.mp3", false)
	}
}
