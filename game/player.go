package game

import (
	"ghost_escape/game/affiliate"
	"ghost_escape/game/core"

	"github.com/SunshineZzzz/purego-sdl3/sdl"
	"github.com/go-gl/mathgl/mgl32"
)

// 玩家
type Player struct {
	// 继承基础角色
	core.Actor
	// 空闲精灵动画
	spriteIdleAnim *affiliate.SpriteAnim
	// 移动精灵动画
	spriteMoveAnim *affiliate.SpriteAnim
	// 是否在运动
	isMoving bool
}

var _ core.IObject = (*Player)(nil)
var _ core.IObjectScreen = (*Player)(nil)

// 初始化
func (p *Player) Init() {
	p.Actor.Init()
	p.MaxSpeed = 500.0
	p.spriteIdleAnim = affiliate.AddSpriteAnimChild(p, "assets/sprite/ghost-idle.png", 2.0)
	p.spriteMoveAnim = affiliate.AddSpriteAnimChild(p, "assets/sprite/ghost-move.png", 2.0)
	p.spriteIdleAnim.SetIsActive(true)
	p.spriteMoveAnim.SetIsActive(false)
	p.isMoving = false
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
	p.move(dt)
	p.checkState()
	p.syncCamera()
}

// 渲染
func (p *Player) Render() {
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

// 移动
func (p *Player) move(dt float32) {
	newPos := p.Position.Add(p.Velocity.Mul(dt))
	p.Position[0] = mgl32.Clamp(newPos.X(), 0.0, p.Game().GetWorldSize().X())
	p.Position[1] = mgl32.Clamp(newPos.Y(), 0.0, p.Game().GetWorldSize().Y())
	p.SetPosition(p.Position)
	// fmt.Printf("dt: %f, pos: %v, vel: %v\n", dt, p.Position, p.Velocity)
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
		p.spriteIdleAnim.SetIsActive(false)
		p.spriteMoveAnim.SetIsActive(true)
		p.spriteMoveAnim.SetFrameTimer(p.spriteIdleAnim.GetCurrentFrame())
		p.spriteMoveAnim.SetCurrentFrame(p.spriteIdleAnim.GetCurrentFrame())
	} else {
		p.spriteIdleAnim.SetIsActive(true)
		p.spriteMoveAnim.SetIsActive(false)
		p.spriteIdleAnim.SetFrameTimer(p.spriteMoveAnim.GetCurrentFrame())
		p.spriteIdleAnim.SetCurrentFrame(p.spriteMoveAnim.GetCurrentFrame())
	}
}
