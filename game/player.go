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
}

var _ core.IObject = (*Player)(nil)
var _ core.IObjectScreen = (*Player)(nil)

// 初始化
func (p *Player) Init() {
	p.Actor.Init()
	p.MaxSpeed = 500.0
	sprite := &affiliate.Sprite{}
	sprite.Init()
	sprite.SetTexture(core.CreateTexture("assets/sprite/ghost-idle.png"))
	sprite.SetParent(p)
	p.AddChild(sprite)
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
	p.Actor.Update(dt)
	p.keyboardControl()
	p.move(dt)
	p.syncCamera()
}

// 渲染
func (p *Player) Render() {
	p.Actor.Render()
	p.Game().DrawBoundary(p.RenderPosition, p.RenderPosition.Add(mgl32.Vec2{20.0, 20.0}), 5.0, sdl.FColor{R: 1.0, G: 0.0, B: 0.0, A: 1.0})
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
