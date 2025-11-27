package game

import (
	"ghost_escape/game/core"

	"github.com/SunshineZzzz/purego-sdl3/sdl"
	"github.com/go-gl/mathgl/mgl32"
)

type SceneMain struct {
	// 继承基础场景
	core.Scene
}

var _ core.IObject = (*SceneMain)(nil)
var _ core.IScene = (*SceneMain)(nil)

func (s *SceneMain) Init() {
	s.Scene.Init()
	s.WorldSize = s.Game().GetScreenSize().Mul(3.0)
	s.CameraPositon = s.WorldSize.Mul(0.5).Sub(s.Game().GetScreenSize().Mul(0.5))
	player := &Player{}
	player.Init()
	player.SetPosition(s.WorldSize.Mul(0.5))
	s.AddChild(player)

	enemy := &Enemy{}
	enemy.Init()
	enemy.SetTarget(player)
	enemy.SetPosition(s.WorldSize.Mul(0.5).Add(mgl32.Vec2{200.0, 200.0}))
	s.AddChild(enemy)
}

func (s *SceneMain) HandleEvent(event *sdl.Event) {
	s.Scene.HandleEvent(event)
}

func (s *SceneMain) Update(dt float32) {
	s.Scene.Update(dt)
}

func (s *SceneMain) Render() {
	s.renderBackground()
	s.Scene.Render()
}

func (s *SceneMain) Clean() {
	s.Scene.Clean()
}

// 非接口实现

// 渲染背景
func (s *SceneMain) renderBackground() {
	// 背景绘制起始
	start := s.WorldToScreen(mgl32.Vec2{0.0, 0.0})
	// 背景绘制结束
	end := s.WorldToScreen(s.WorldSize)
	s.Game().DrawGrid(start, end, 80.0, sdl.FColor{R: 0.5, G: 0.5, B: 0.5, A: 1.0})
	s.Game().DrawBoundary(start, end, 5.0, sdl.FColor{R: 1.0, G: 1.0, B: 1.0, A: 1.0})
}
