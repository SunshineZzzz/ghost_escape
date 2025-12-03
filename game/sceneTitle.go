package game

import (
	"ghost_escape/game/core"
	"ghost_escape/game/screen"
	"math"
	"strconv"

	"github.com/SunshineZzzz/purego-sdl3/sdl"
	"github.com/go-gl/mathgl/mgl32"
)

// 标题场景
type SceneTitle struct {
	// 继承基础场景
	core.Scene
	// 边界颜色
	boundaryColor sdl.FColor
	// 颜色计时器
	colorTimer float32
	// 开始按钮
	startButton *screen.HudButton
	// 退出按钮
	quitButton *screen.HudButton
	// 贡献者名单按钮
	creditButton *screen.HudButton
}

var _ core.IObject = (*SceneTitle)(nil)
var _ core.IScene = (*SceneTitle)(nil)

// 初始化
func (s *SceneTitle) Init() {
	s.Scene.Init()
	size := mgl32.Vec2{s.Game().GetScreenSize().X() / 2.0, s.Game().GetScreenSize().Y() / 3.0}
	screen.AddHudTextChild(s, "幽 灵 逃 生", s.Game().GetScreenSize().Mul(0.5).Add(mgl32.Vec2{0.0, -100}),
		size, "assets/font/VonwaonBitmap-16px.ttf", 64, "assets/UI/Textfield_01.png", core.AnchorTypeCenter)
	scoreText := "最高分: " + strconv.Itoa(s.Game().GetHighScore())
	screen.AddHudTextChild(s, scoreText, s.Game().GetScreenSize().Mul(0.5).Add(mgl32.Vec2{0.0, 100}),
		mgl32.Vec2{200, 50}, "assets/font/VonwaonBitmap-16px.ttf", 32, "assets/UI/Textfield_01.png", core.AnchorTypeCenter)

	// 退出按钮
	s.quitButton = screen.AddHudButtonChild(s, s.Game().GetScreenSize().Mul(0.5).Add(mgl32.Vec2{200.0, 200.0}),
		"assets/UI/A_Quit1.png", "assets/UI/A_Quit2.png", "assets/UI/A_Quit3.png", 2.0, core.AnchorTypeCenter)

}

// 处理事件
func (s *SceneTitle) HandleEvent(event *sdl.Event) {
	s.Scene.HandleEvent(event)
}

// 更新
func (s *SceneTitle) Update(dt float32) {
	s.Scene.Update(dt)
	s.colorTimer += dt
	s.updateColor()
	s.checkButtonQuit()
}

func (s *SceneTitle) Render() {
	s.renderBackground()
	s.Scene.Render()
}

func (s *SceneTitle) Clean() {
	s.Scene.Clean()
}

// 非接口实现

// 渲染背景
func (s *SceneTitle) renderBackground() {
	s.Game().DrawBoundary(mgl32.Vec2{30.0, 30.0}, s.Game().GetScreenSize().Sub(mgl32.Vec2{30.0, 30.0}), 10.0, s.boundaryColor)
}

// 更新颜色
func (s *SceneTitle) updateColor() {
	s.boundaryColor.R = 0.5 + 0.5*float32(math.Sin(float64(s.colorTimer*0.9)))
	s.boundaryColor.G = 0.5 + 0.5*float32(math.Sin(float64(s.colorTimer*0.8)))
	s.boundaryColor.B = 0.5 + 0.5*float32(math.Sin(float64(s.colorTimer*0.7)))
}

// 检查退出按钮是否触发
func (s *SceneTitle) checkButtonQuit() {
	if s.quitButton.GetIsTrigger() {
		s.Game().Quit()
	}
}
