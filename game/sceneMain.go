package game

import (
	"ghost_escape/game/core"
	"ghost_escape/game/screen"

	"github.com/SunshineZzzz/purego-sdl3/sdl"
	"github.com/go-gl/mathgl/mgl32"
)

type SceneMain struct {
	// 继承基础场景
	core.Scene
	// 生成器
	spawner *Spawner
	// UI鼠标
	uimouse *screen.UIMouse
}

var _ core.IObject = (*SceneMain)(nil)
var _ core.IScene = (*SceneMain)(nil)

func (s *SceneMain) Init() {
	s.Scene.Init()
	s.WorldSize = s.Game().GetScreenSize().Mul(3.0)
	s.CameraPositon = s.WorldSize.Mul(0.5).Sub(s.Game().GetScreenSize().Mul(0.5))

	// 玩家
	player := &Player{}
	player.Init()
	player.SetPosition(s.WorldSize.Mul(0.5))
	s.AddChild(player)

	// 生成器
	spawner := &Spawner{}
	spawner.Init()
	spawner.SetTarget(player)
	s.spawner = spawner
	s.AddChild(spawner)

	// UI鼠标
	s.uimouse = screen.AddUIMouse(s, "assets/UI/29.png", "assets/UI/30.png", 1.0, core.AnchorTypeCenter)

	// // 敌人
	// enemy := &Enemy{}
	// enemy.Init()
	// enemy.SetTarget(player)
	// enemy.SetPosition(s.WorldSize.Mul(0.5).Add(mgl32.Vec2{200.0, 200.0}))
	// // 敌人产生是从特效精灵动画结束后产生，所以这里生成特效精灵动画
	// SpriteAnim := affiliate.CteateSpriteAnimChild("assets/effect/184_3.png", 1.0, core.AnchorTypeCenter)
	// // 上面的特效加载到场景中
	// core.AddEffect(s, SpriteAnim, enemy.GetPosition(), enemy)
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
