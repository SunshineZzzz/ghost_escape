package game

import (
	"encoding/binary"
	"ghost_escape/game/core"
	"ghost_escape/game/screen"
	"math"
	"os"
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
	creditsButton *screen.HudButton
	// 贡献者名单文本
	creditsText *screen.HudText
	// UI鼠标
	uimouse *screen.UIMouse
}

var _ core.IObject = (*SceneTitle)(nil)
var _ core.IScene = (*SceneTitle)(nil)

// 初始化
func (s *SceneTitle) Init() {
	s.Scene.Init()
	s.LoadData("assets/score.dat")
	sdl.HideCursor()
	s.Game().StopAllMusic()
	s.Game().StopAllEffects()
	s.Game().PlayMusic("assets/bgm/Spooky music.mp3", true)
	size := mgl32.Vec2{s.Game().GetScreenSize().X() / 2.0, s.Game().GetScreenSize().Y() / 3.0}
	screen.AddHudTextChild(s, "幽 灵 逃 生", s.Game().GetScreenSize().Mul(0.5).Add(mgl32.Vec2{0.0, -100}),
		size, "assets/font/VonwaonBitmap-16px.ttf", 64, "assets/UI/Textfield_01.png", core.AnchorTypeCenter)
	scoreText := "最高分: " + strconv.Itoa(s.Game().GetHighScore())
	screen.AddHudTextChild(s, scoreText, s.Game().GetScreenSize().Mul(0.5).Add(mgl32.Vec2{0.0, 100}),
		mgl32.Vec2{200, 50}, "assets/font/VonwaonBitmap-16px.ttf", 32, "assets/UI/Textfield_01.png", core.AnchorTypeCenter)

	// 开始按钮
	s.startButton = screen.AddHudButtonChild(s, s.Game().GetScreenSize().Mul(0.5).Add(mgl32.Vec2{-200.0, 200.0}),
		"assets/UI/A_Start1.png", "assets/UI/A_Start2.png", "assets/UI/A_Start3.png", 2.0, core.AnchorTypeCenter)
	// 贡献者名单按钮
	s.creditsButton = screen.AddHudButtonChild(s, s.Game().GetScreenSize().Mul(0.5).Add(mgl32.Vec2{0, 200.0}),
		"assets/UI/A_Credits1.png", "assets/UI/A_Credits2.png", "assets/UI/A_Credits3.png", 2.0, core.AnchorTypeCenter)
	// 退出按钮
	s.quitButton = screen.AddHudButtonChild(s, s.Game().GetScreenSize().Mul(0.5).Add(mgl32.Vec2{200.0, 200.0}),
		"assets/UI/A_Quit1.png", "assets/UI/A_Quit2.png", "assets/UI/A_Quit3.png", 2.0, core.AnchorTypeCenter)

	text, err := s.Game().LoadTextFromFile("assets/credits.txt")
	if err != nil {
		return
	}
	s.creditsText = screen.AddHudTextChild(s, text, s.Game().GetScreenSize().Mul(0.5),
		mgl32.Vec2{500, 500}, "assets/font/VonwaonBitmap-16px.ttf", 16, "assets/UI/Textfield_01.png", core.AnchorTypeCenter)
	s.creditsText.SetActive(false)
	s.creditsText.SetBgSizeByText(50.0)

	// UI鼠标
	s.uimouse = screen.AddUIMouseChild(s, "assets/UI/pointer_c_shaded.png", "assets/UI/pointer_c_shaded.png", 1.0, core.AnchorTypeTopLeft)
}

// 处理事件
func (s *SceneTitle) HandleEvent(event *sdl.Event) {
	if s.creditsText.GetActive() {
		if event.Type() == sdl.EventMouseButtonUp {
			s.creditsText.SetActive(false)
		}
		return
	}
	s.Scene.HandleEvent(event)
}

// 更新
func (s *SceneTitle) Update(dt float32) {
	s.colorTimer += dt
	s.updateColor()
	if s.creditsText.GetActive() {
		s.uimouse.Update(dt)
		return
	}
	s.Scene.Update(dt)
	s.checkButtonQuit()
	s.checkButtonStart()
	s.checkButtonCredits()
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

// 检查开始按钮是否触发
func (s *SceneTitle) checkButtonStart() {
	if s.startButton.GetIsTrigger() {
		s.Game().SafeChangeScene(&SceneMain{})
	}
}

// 检查贡献者名单按钮是否触发
func (s *SceneTitle) checkButtonCredits() {
	if s.creditsButton.GetIsTrigger() {
		s.creditsText.SetActive(true)
	}
}

// 加载数据
func (s *SceneTitle) LoadData(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	var highScore int32
	err = binary.Read(file, binary.LittleEndian, &highScore)
	if err != nil {
		return
	}
	s.Game().SetHighScore(int(highScore))
}
