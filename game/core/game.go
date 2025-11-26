package core

import (
	"fmt"
	"sync"

	"github.com/SunshineZzzz/purego-sdl3/sdl"
	"github.com/SunshineZzzz/purego-sdl3/ttf"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	// 游戏帧率
	FPS = 60
)

var (
	instance *Game
	once     sync.Once
)

func GetInstance() *Game {
	once.Do(func() {
		instance = &Game{
			fps:         FPS,
			dt:          0.0,
			frameDelay:  1e9 / FPS,
			isRunning:   false,
			sdlWindow:   nil,
			sdlRenderer: nil,
		}
	})
	return instance
}

type Game struct {
	// 屏幕大小
	screenSize mgl32.Vec2
	// 游戏帧率
	fps uint64
	// 帧间隔，单位秒
	dt float32
	// 帧延迟，单位纳秒
	frameDelay float32
	// 是否运行中
	isRunning bool
	// SDL窗口
	sdlWindow *sdl.Window
	// SDL渲染器
	sdlRenderer *sdl.Renderer
	// 当前场景
	currentScene IScene
}

func (g *Game) Init(title string, width, height int32, scene IScene) error {
	g.screenSize = mgl32.Vec2{float32(width), float32(height)}

	// 初始化 SDL
	if !sdl.Init(sdl.InitVideo | sdl.InitAudio | sdl.InitEvents) {
		return fmt.Errorf("sdl init error,%s", sdl.GetError())
	}

	// 创建窗口与渲染器
	if !sdl.CreateWindowAndRenderer(title, width, height, sdl.WindowResizable, &g.sdlWindow, &g.sdlRenderer) {
		return fmt.Errorf("sdl create window and renderer error,%s", sdl.GetError())
	}

	// 设置渲染器的逻辑尺寸
	if !sdl.SetRenderLogicalPresentation(g.sdlRenderer, width, height, sdl.LogicalPresentationLetterbox) {
		return fmt.Errorf("sdl set render logical presentation error,%s", sdl.GetError())
	}

	// 初始化TTF
	if !ttf.Init() {
		return fmt.Errorf("ttf init error,%s", sdl.GetError())
	}

	g.currentScene = scene
	g.currentScene.Init()

	g.isRunning = true
	return nil
}

func (g *Game) Run() {
	// 主循环
	for g.isRunning {
		start := sdl.GetTicksNS()
		g.handleEvent()
		g.update(g.dt)
		g.render()
		end := sdl.GetTicksNS()
		elapsed := float32(end - start)
		if elapsed < g.frameDelay {
			sdl.DelayNS(uint64(g.frameDelay - elapsed))
			g.dt = g.frameDelay / 1e9
		} else {
			g.dt = elapsed / 1e9
		}
		// fmt.Printf("dt: %f\n", g.dt)
	}
}

// 处理事件
func (g *Game) handleEvent() {
	var event sdl.Event
	for sdl.PollEvent(&event) {
		if event.Type() == sdl.EventQuit {
			g.isRunning = false
			return
		}
		g.currentScene.HandleEvent(&event)
	}
}

// 更新状态
func (g *Game) update(dt float32) {
	g.currentScene.Update(dt)
}

// 渲染
func (g *Game) render() {
	// 清空渲染器
	sdl.RenderClear(g.sdlRenderer)

	// 渲染当前场景
	g.currentScene.Render()

	// 显示更新
	sdl.RenderPresent(g.sdlRenderer)
}

// 清理资源
func (g *Game) Clean() {
	if g.currentScene != nil {
		g.currentScene.Clean()
		g.currentScene = nil
	}

	// 清理SDL资源
	if g.sdlRenderer != nil {
		sdl.DestroyRenderer(g.sdlRenderer)
		g.sdlRenderer = nil
	}
	if g.sdlWindow != nil {
		sdl.DestroyWindow(g.sdlWindow)
		g.sdlWindow = nil
	}
	sdl.Quit()
}

// 获取屏幕大小
func (g *Game) GetScreenSize() mgl32.Vec2 {
	return g.screenSize
}

// 获取世界大小
func (g *Game) GetWorldSize() mgl32.Vec2 {
	return g.currentScene.GetWorldSize()
}

// 绘制网格
func (g *Game) DrawGrid(topLeft, bottomRight mgl32.Vec2, gridWidth float32, fcolor sdl.FColor) {
	sdl.SetRenderDrawColorFloat(g.sdlRenderer, fcolor.R, fcolor.G, fcolor.B, fcolor.A)
	screenRect := sdl.FRect{
		X: 0.0,
		Y: 0.0,
		W: g.screenSize.X(),
		H: g.screenSize.Y(),
	}
	for x := topLeft.X(); x < bottomRight.X(); x += gridWidth {
		if !sdl.PointInRectFloat(sdl.FPoint{X: x, Y: 0.0}, screenRect) {
			continue
		}
		sdl.RenderLine(g.sdlRenderer, x, topLeft.Y(), x, bottomRight.Y())
	}
	for y := topLeft.Y(); y < bottomRight.Y(); y += gridWidth {
		if !sdl.PointInRectFloat(sdl.FPoint{X: 0.0, Y: y}, screenRect) {
			continue
		}
		sdl.RenderLine(g.sdlRenderer, topLeft.X(), y, bottomRight.X(), y)
	}
	sdl.SetRenderDrawColorFloat(g.sdlRenderer, 0, 0, 0, 1)
}

// 绘制边界
func (g *Game) DrawBoundary(topLeft, bottomRight mgl32.Vec2, boundaryWidth float32, fcolor sdl.FColor) {
	sdl.SetRenderDrawColorFloat(g.sdlRenderer, fcolor.R, fcolor.G, fcolor.B, fcolor.A)
	screenRect := sdl.FRect{
		X: 0.0,
		Y: 0.0,
		W: g.screenSize.X(),
		H: g.screenSize.Y(),
	}
	for i := float32(0.0); i < boundaryWidth; i++ {
		rect := sdl.FRect{
			X: topLeft.X() - i,
			Y: topLeft.Y() - i,
			W: bottomRight.X() - topLeft.X() + 2*i,
			H: bottomRight.Y() - topLeft.Y() + 2*i,
		}
		intersectionRect, ok := sdl.GetRectIntersectionFloat(screenRect, rect)
		if !ok {
			continue
		}
		sdl.RenderRect(g.sdlRenderer, &intersectionRect)
	}
	sdl.SetRenderDrawColorFloat(g.sdlRenderer, 0, 0, 0, 1)
}

// 获取当前场景
func (g *Game) GetCurrentScene() IScene {
	return g.currentScene
}
