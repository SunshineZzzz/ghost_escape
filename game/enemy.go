package game

import (
	"ghost_escape/game/affiliate"
	"ghost_escape/game/core"
)

// 敌人状态
type EnemyState int

const (
	// 正常状态
	EnemyStateNormal EnemyState = iota
	// 受伤状态
	EnemyStateHurt
	// 死亡状态
	EnemyStateDead
)

// 敌人
type Enemy struct {
	// 继承基础角色
	core.Actor
	// 目标玩家
	target *Player
	// 敌人状态
	currentState EnemyState
	// 正常状态精灵动画
	spriteAnimNormal *affiliate.SpriteAnim
	// 受伤状态精灵动画
	spriteAnimHurt *affiliate.SpriteAnim
	// 死亡状态精灵动画
	spriteAnimDead *affiliate.SpriteAnim
	// 当前精灵动画
	currentSpriteAnim *affiliate.SpriteAnim
	// 暂时测试用的
	timer float32
}

var _ core.IObject = (*Enemy)(nil)
var _ core.IObjectScreen = (*Enemy)(nil)

// 初始化
func (e *Enemy) Init() {
	e.Actor.Init()
	e.currentState = EnemyStateNormal
	e.spriteAnimNormal = affiliate.AddSpriteAnimChild(e, "assets/sprite/ghost-Sheet.png", 2.0)
	e.spriteAnimHurt = affiliate.AddSpriteAnimChild(e, "assets/sprite/ghostHurt-Sheet.png", 2.0)
	e.spriteAnimDead = affiliate.AddSpriteAnimChild(e, "assets/sprite/ghostDead-Sheet.png", 2.0)
	e.spriteAnimNormal.SetActive(true)
	e.spriteAnimHurt.SetActive(false)
	e.spriteAnimDead.SetActive(false)
	e.spriteAnimDead.SetLoop(false)
	e.currentSpriteAnim = e.spriteAnimNormal
}

// 更新
func (e *Enemy) Update(dt float32) {
	e.Actor.Update(dt)
	e.aimTarget(e.target)
	e.Move(dt)
	e.timer += dt
	if e.timer > 1.0 && e.timer < 2.0 {
		e.changeState(EnemyStateHurt)
	} else if e.timer > 2.0 && e.timer < 3.0 {
		e.changeState(EnemyStateDead)
	} else {
		e.Remove()
	}
}

// 非接口实现

// 设置目标玩家
func (e *Enemy) SetTarget(target *Player) {
	e.target = target
}

// 瞄准目标
func (e *Enemy) aimTarget(target *Player) {
	if target == nil {
		return
	}
	// 计算目标方向，并且归一化
	direction := target.GetPosition().Sub(e.GetPosition()).Normalize()
	// 设置速度
	e.SetVelocity(direction.Mul(e.GetMaxSpeed()))
}

// 改变状态
func (e *Enemy) changeState(newState EnemyState) {
	if e.currentState == newState {
		return
	}
	e.currentSpriteAnim.SetActive(false)

	switch newState {
	case EnemyStateNormal:
		e.currentSpriteAnim = e.spriteAnimNormal
		e.currentSpriteAnim.SetActive(true)
	case EnemyStateHurt:
		e.currentSpriteAnim = e.spriteAnimHurt
		e.currentSpriteAnim.SetActive(true)
	case EnemyStateDead:
		e.currentSpriteAnim = e.spriteAnimDead
		e.currentSpriteAnim.SetActive(true)
	}

	e.currentState = newState
}

// 移除
func (e *Enemy) Remove() {
	if e.currentSpriteAnim.GetFinish() {
		e.SetNeedRemove(true)
	}
}
