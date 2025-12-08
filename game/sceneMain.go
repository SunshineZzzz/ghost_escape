package game

import (
	"encoding/binary"
	"ghost_escape/game/core"
	"ghost_escape/game/raw"
	"ghost_escape/game/screen"
	"os"
	"strconv"

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
	// HUD状态
	hudStats *screen.HudStats
	// HUD技能
	hudSkills *screen.HudSkill
	// HUD分数
	hudScore *screen.HudText
	// 暂停按钮
	buttonPause *screen.HudButton
	// 重新游戏按钮
	buttonRestart *screen.HudButton
	// 退到标题场景按钮
	buttonBack *screen.HudButton
	// 游戏结束timer
	endTimer *core.Timer
	// 玩家
	player *Player
}

var _ core.IObject = (*SceneMain)(nil)
var _ core.IScene = (*SceneMain)(nil)

func (s *SceneMain) Init() {
	s.Scene.Init()
	sdl.HideCursor()
	s.Game().StopAllMusic()
	s.Game().StopAllEffects()
	s.Game().PlayMusic("assets/bgm/OhMyGhost.ogg", true)
	s.WorldSize = s.Game().GetScreenSize().Mul(3.0)
	s.CameraPositon = s.WorldSize.Mul(0.5).Sub(s.Game().GetScreenSize().Mul(0.5))

	// 玩家
	s.player = &Player{}
	s.player.Init()
	s.player.SetPosition(s.WorldSize.Mul(0.5))
	s.AddChild(s.player)

	// 增加视差星空背景
	raw.AddBgStarChild(s, 200, 0.2, 0.5, 0.7)

	// 生成器
	spawner := &Spawner{}
	spawner.Init()
	spawner.SetTarget(s.player)
	s.spawner = spawner
	s.AddChild(spawner)

	// HUD状态
	s.hudStats = screen.AddHudStatsChild(s, &s.player.Actor, mgl32.Vec2{30.0, 30.0})
	// HUD技能
	s.hudSkills = screen.AddHudSkillChild(s, s.player, "assets/UI/Electric-Icon.png", mgl32.Vec2{s.player.Game().GetScreenSize().X() - 300.0, 30.0}, 0.14, core.AnchorTypeCenter)
	// HUD分数
	s.hudScore = screen.AddHudTextChild(s, "Score: 0", mgl32.Vec2{s.player.Game().GetScreenSize().X() - 120.0, 30.0}, mgl32.Vec2{200.0, 50.0},
		"assets/font/VonwaonBitmap-16px.ttf", 32.0, "assets/UI/Textfield_01.png", core.AnchorTypeCenter)

	s.buttonPause = screen.AddHudButtonChild(s, s.Game().GetScreenSize().Add(mgl32.Vec2{-230.0, -30.0}), "assets/UI/A_Pause1.png", "assets/UI/A_Pause2.png", "assets/UI/A_Pause3.png", 1.0, core.AnchorTypeCenter)
	s.buttonRestart = screen.AddHudButtonChild(s, s.Game().GetScreenSize().Add(mgl32.Vec2{-140.0, -30.0}), "assets/UI/A_Restart1.png", "assets/UI/A_Restart2.png", "assets/UI/A_Restart3.png", 1.0, core.AnchorTypeCenter)
	s.buttonBack = screen.AddHudButtonChild(s, s.Game().GetScreenSize().Add(mgl32.Vec2{-50.0, -30.0}), "assets/UI/A_Back1.png", "assets/UI/A_Back2.png", "assets/UI/A_Back3.png", 1.0, core.AnchorTypeCenter)

	// UI鼠标
	s.uimouse = screen.AddUIMouseChild(s, "assets/UI/29.png", "assets/UI/30.png", 1.0, core.AnchorTypeCenter)

	// 游戏结束timer
	s.endTimer = core.AddTimerChild(s, 3.0)

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
	s.updateScore()
	s.checkButtonRestart()
	s.checkButtonBack()
	s.checkButtonPause()
	if s.player != nil && !s.player.GetActive() {
		s.endTimer.Start()
		s.SaveData("assets/score.dat")
	}
	s.checkEndTimer()
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

// 更新分数
func (s *SceneMain) updateScore() {
	if s.hudScore == nil {
		return
	}
	s.hudScore.SetText("Score: " + strconv.Itoa(s.Game().GetScore()))
}

// 检查重新游戏按钮
func (s *SceneMain) checkButtonRestart() {
	if s.buttonRestart == nil {
		return
	}
	if !s.buttonRestart.GetIsTrigger() {
		return
	}
	s.SaveData("assets/score.dat")
	s.Game().SetScore(0)
	s.Game().SafeChangeScene(s)
}

// 检查退到标题场景按钮
func (s *SceneMain) checkButtonBack() {
	if s.buttonBack == nil {
		return
	}
	if !s.buttonBack.GetIsTrigger() {
		return
	}
	s.SaveData("assets/score.dat")
	s.Game().SetScore(0)
	s.Game().SafeChangeScene(&SceneTitle{})
}

// 检查暂停按钮
func (s *SceneMain) checkButtonPause() {
	if s.buttonPause == nil {
		return
	}
	if !s.buttonPause.GetIsTrigger() {
		return
	}
	if s.IsPause {
		s.Resume()
		return
	}
	s.Pause()
}

// 检查游戏结束timer
func (s *SceneMain) checkEndTimer() {
	if s.endTimer != nil && !s.endTimer.TimeOut() {
		return
	}
	s.Pause()
	s.buttonRestart.SetRenderPosition(s.Game().GetScreenSize().Mul(0.5).Add(mgl32.Vec2{-200.0, 0.0}))
	s.buttonRestart.SetScale(4.0)
	s.buttonBack.SetRenderPosition(s.Game().GetScreenSize().Mul(0.5).Add(mgl32.Vec2{200.0, 0.0}))
	s.buttonBack.SetScale(4.0)
	s.buttonPause.SetActive(false)
	s.endTimer.Stop()
}

// 保存数据
func (s *SceneMain) SaveData(filePath string) {
	highScore := s.Game().GetHighScore()
	file, err := os.Create(filePath)
	if err != nil {
		return
	}
	defer file.Close()
	err = binary.Write(file, binary.LittleEndian, int32(highScore))
	if err != nil {
		return
	}
}
