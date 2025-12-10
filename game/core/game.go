package core

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

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
	// 随机数生成器
	rand *rand.Rand
	// 资源管理器
	assetStore *AssetStore
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
	// 字体引擎
	ttfEngine *ttf.TextEngine
	// 当前场景
	currentScene IScene
	// 鼠标位置
	// 物理窗口中鼠标坐标，它的范围是从 (0, 0) 到 (WindowWidth, WindowHeight)。如果您把窗口从 1000x800 拉伸到 2000x1600，这些坐标的最大值也会随之变化。
	// 需要转化到 游戏逻辑坐标
	// SDL 会负责将 800x600 逻辑画面，根据 LETTERBOX 规则，缩放、居中到玩家的物理窗口上。
	// 这就是为什么鼠标的物理坐标（例如 1600x900 范围）必须被 “翻译”（去缩放、去偏移）回 800x600 的逻辑坐标，您的游戏代码才能正确识别点击位置。
	// 放大以后或者缩小以后，逻辑分辨率始终是800X600
	mousePosition mgl32.Vec2
	// 鼠标按钮状态
	mouseButtons sdl.MouseButtonFlags
	// 分数
	score int
	// 最高分
	highScore int
	// 下一个场景
	nextScene IScene
}

func (g *Game) Init(title string, width, height int32, scene IScene) error {
	g.screenSize = mgl32.Vec2{float32(width), float32(height)}
	g.rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	// 初始化 SDL
	if !sdl.Init(sdl.InitVideo | sdl.InitAudio | sdl.InitEvents) {
		return fmt.Errorf("sdl init error,%s", sdl.GetError())
	}

	// 创建窗口与渲染器
	if !sdl.CreateWindowAndRenderer(title, width, height, sdl.WindowResizable, &g.sdlWindow, &g.sdlRenderer) {
		return fmt.Errorf("sdl create window and renderer error,%s", sdl.GetError())
	}

	// 设置渲染器的逻辑尺寸
	// sdl.LogicalPresentationLetterbox
	// 它会把游戏画面放大到窗口允许的最大尺寸，同时不改变画面的比例。
	// 如果窗口比逻辑画面宽，您会看到左右两侧有黑边(Letterbox)
	// 如果窗口比逻辑画面高，您会看到顶部和底部有黑边(Letterbox)
	if !sdl.SetRenderLogicalPresentation(g.sdlRenderer, width, height, sdl.LogicalPresentationLetterbox) {
		return fmt.Errorf("sdl set render logical presentation error,%s", sdl.GetError())
	}

	// 初始化TTF
	if !ttf.Init() {
		return fmt.Errorf("ttf init error,%s", sdl.GetError())
	}
	// 创建字体引擎
	g.ttfEngine = ttf.CreateRendererTextEngine(g.sdlRenderer)

	// 创建资源管理器
	g.assetStore = CreateAssetStore(g.sdlRenderer)

	g.currentScene = scene
	g.currentScene.Init()

	g.isRunning = true
	return nil
}

func (g *Game) Run() {
	// 主循环
	for g.isRunning {
		start := sdl.GetTicksNS()
		if g.nextScene != nil {
			g.ChangeScene(g.nextScene)
			g.nextScene = nil
		}
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
	g.updateMouse()
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
	if g.ttfEngine != nil {
		ttf.DestroyRendererTextEngine(g.ttfEngine)
		g.ttfEngine = nil
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

// 获取资源管理器
func (g *Game) GetAssetStore() *AssetStore {
	return g.assetStore
}

// 渲染纹理
func (g *Game) RenderTexture(texture *Texture, pos mgl32.Vec2, size mgl32.Vec2, percent mgl32.Vec2) {
	srcRect := sdl.FRect{
		X: texture.SrcRect.X,
		Y: texture.SrcRect.Y,
		W: texture.SrcRect.W * percent.X(),
		H: texture.SrcRect.H * percent.Y(),
	}
	dstRect := sdl.FRect{
		X: pos.X(),
		Y: pos.Y(),
		W: size.X() * percent.X(),
		H: size.Y() * percent.Y(),
	}
	screenRect := sdl.FRect{
		X: 0.0,
		Y: 0.0,
		W: g.screenSize.X(),
		H: g.screenSize.Y(),
	}
	intersectionRect, ok := sdl.GetRectIntersectionFloat(screenRect, dstRect)
	if !ok {
		return
	}
	flipMode := sdl.FlipNone
	if texture.IsFlip {
		flipMode = sdl.FlipHorizontal
	}
	sdl.RenderTextureRotated(g.sdlRenderer, texture.Texture, &srcRect, &intersectionRect, texture.Angle, nil, flipMode)
}

// 绘制填充圆，并不是画圆，而是用绘制圆形纹理，目的是可视化碰撞器
func (g *Game) RenderFillCircle(pos mgl32.Vec2, size mgl32.Vec2, alpha float32) {
	dstRect := sdl.FRect{
		X: pos.X(),
		Y: pos.Y(),
		W: size.X(),
		H: size.Y(),
	}
	screenRect := sdl.FRect{
		X: 0.0,
		Y: 0.0,
		W: g.screenSize.X(),
		H: g.screenSize.Y(),
	}
	intersectionRect, ok := sdl.GetRectIntersectionFloat(screenRect, dstRect)
	if !ok {
		return
	}
	texture, err := g.assetStore.GetImage("assets/UI/circle.png")
	if err != nil {
		return
	}
	sdl.SetTextureAlphaModFloat(texture, alpha)
	sdl.RenderTexture(g.sdlRenderer, texture, nil, &intersectionRect)
}

// 随机min和max范围的浮点数
func (g *Game) RandFloat32(min, max float32) float32 {
	return min + g.rand.Float32()*(max-min)
}

// 随机min和max范围的整数
func (g *Game) RandInt(min, max int) int {
	return min + g.rand.Int()*(max-min)
}

// 随机min和max范围的Vec2
func (g *Game) RandVec2(min, max mgl32.Vec2) mgl32.Vec2 {
	return mgl32.Vec2{
		g.RandFloat32(min.X(), max.X()),
		g.RandFloat32(min.Y(), max.Y()),
	}
}

// 获取鼠标位置
func (g *Game) GetMousePosition() mgl32.Vec2 {
	return g.mousePosition
}

// 获取鼠标按钮状态
func (g *Game) GetMouseButtons() sdl.MouseButtonFlags {
	return g.mouseButtons
}

// 绘制水平进度条
func (g *Game) RenderHBar(pos mgl32.Vec2, size mgl32.Vec2, percent mgl32.Vec2, color sdl.FColor) {
	boundaryRect := sdl.FRect{
		X: pos.X(),
		Y: pos.Y(),
		W: size.X(),
		H: size.Y(),
	}
	fillRect := sdl.FRect{
		X: pos.X(),
		Y: pos.Y(),
		W: size.X() * percent.X(),
		H: size.Y() * percent.Y(),
	}
	screenRect := sdl.FRect{
		X: 0.0,
		Y: 0.0,
		W: g.screenSize.X(),
		H: g.screenSize.Y(),
	}
	intersectionRect1, ok1 := sdl.GetRectIntersectionFloat(screenRect, boundaryRect)
	intersectionRect2, ok2 := sdl.GetRectIntersectionFloat(screenRect, fillRect)
	if !ok1 && !ok2 {
		return
	}
	sdl.SetRenderDrawColorFloat(g.sdlRenderer, color.R, color.G, color.B, color.A)
	if ok1 {
		sdl.RenderRect(g.sdlRenderer, &intersectionRect1)
	}
	if ok2 {
		sdl.RenderFillRect(g.sdlRenderer, &intersectionRect2)
	}
	sdl.SetRenderDrawColorFloat(g.sdlRenderer, 0.0, 0.0, 0.0, 1.0)
}

// 创建TTF文本
func (g *Game) CreateTTFText(text string, fontPath string, fontSize float32) *ttf.Text {
	font, err := g.assetStore.GetFont(fontPath, fontSize)
	if err != nil {
		return nil
	}
	return ttf.CreateText(g.ttfEngine, font, text, 0)
}

// 设置分数
func (g *Game) SetScore(score int) {
	g.score = score
	if score > g.highScore {
		g.highScore = score
	}
}

// 获取分数
func (g *Game) GetScore() int {
	return g.score
}

// 增加分数
func (g *Game) AddScore(score int) {
	g.SetScore(g.score + score)
}

// 获取最高分
func (g *Game) GetHighScore() int {
	return g.highScore
}

// 设置最高分
func (g *Game) SetHighScore(highScore int) {
	g.highScore = highScore
}

// 播放音乐
func (g *Game) PlayMusic(musicPath string, loop bool) {
	music, err := g.assetStore.GetSound(musicPath, SoundTypeMusic)
	if err != nil {
		return
	}
	for _, m := range music {
		m.SetLoop(loop)
		if m.Play() {
			return
		}
	}
	// 走到这里肯定是所有池子都失败了
	if err = g.assetStore.loadSound(musicPath, SoundTypeMusic); err != nil {
		fmt.Printf("load music error,%v\n", err)
		return
	}
	music, err = g.assetStore.GetSound(musicPath, SoundTypeMusic)
	if err != nil {
		fmt.Printf("load music error,%v\n", err)
		return
	}
	// fmt.Printf("play new music %s\n", musicPath)
	music[len(music)-1].Play()
}

// 播放音效
func (g *Game) PlaySound(soundPath string, loop bool) {
	sound, err := g.assetStore.GetSound(soundPath, SoundTypeEffect)
	if err != nil {
		return
	}
	for _, s := range sound {
		s.SetLoop(loop)
		if s.Play() {
			return
		}
	}
	// 走到这里肯定是所有池子都失败了
	if err = g.assetStore.loadSound(soundPath, SoundTypeEffect); err != nil {
		fmt.Printf("load sound error,%v\n", err)
		return
	}
	sound, err = g.assetStore.GetSound(soundPath, SoundTypeEffect)
	if err != nil {
		fmt.Printf("load sound error,%v\n", err)
		return
	}
	// fmt.Printf("play new sound %s\n", soundPath)
	sound[len(sound)-1].Play()
}

// 停止所有音乐
func (g *Game) StopAllMusic() {
	g.assetStore.StopAllSound(SoundTypeMusic)
}

// 停止所有音效
func (g *Game) StopAllEffects() {
	g.assetStore.StopAllSound(SoundTypeEffect)
}

// 暂停所有音效
func (g *Game) PauseAllEffects() {
	g.assetStore.PauseAllSound(SoundTypeEffect)
}

// 恢复所有音效
func (g *Game) ResumeAllEffects() {
	g.assetStore.ResumeAllSound(SoundTypeEffect)
}

// 暂停所有音乐
func (g *Game) PauseAllMusic() {
	g.assetStore.PauseAllSound(SoundTypeMusic)
}

// 恢复所有音乐
func (g *Game) ResumeAllMusic() {
	g.assetStore.ResumeAllSound(SoundTypeMusic)
}

// 离开游戏
func (g *Game) Quit() {
	g.isRunning = false
}

// 鼠标位置是否在矩形内
func (g *Game) IsMouseInRect(topLeft, bottomRight mgl32.Vec2) bool {
	if g.mousePosition.X() >= topLeft.X() && g.mousePosition.X() <= bottomRight.X() &&
		g.mousePosition.Y() >= topLeft.Y() && g.mousePosition.Y() <= bottomRight.Y() {
		return true
	}
	return false
}

// 从文件中加载文本
func (g *Game) LoadTextFromFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	text := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		text += line + "\n"
	}
	return text, nil
}

// 安全切换场景
func (g *Game) SafeChangeScene(scene IScene) {
	g.nextScene = scene
}

// 切换场景
func (g *Game) ChangeScene(scene IScene) {
	if g.currentScene != nil {
		g.currentScene.Clean()
	}
	g.currentScene = scene
	g.currentScene.Init()
}

// 绘制点们
func (g *Game) DrawPoints(points *[]mgl32.Vec2, renderPos mgl32.Vec2, color sdl.FColor) {
	sdl.SetRenderDrawColorFloat(g.sdlRenderer, color.R, color.G, color.B, color.A)
	screenRect := sdl.FRect{
		X: 0.0,
		Y: 0.0,
		W: g.screenSize.X(),
		H: g.screenSize.Y(),
	}
	for _, p := range *points {
		pos := sdl.FPoint{X: p.X() + renderPos.X(), Y: p.Y() + renderPos.Y()}
		if !sdl.PointInRectFloat(pos, screenRect) {
			continue
		}
		sdl.RenderPoint(g.sdlRenderer, pos.X, pos.Y)
	}
	sdl.SetRenderDrawColorFloat(g.sdlRenderer, 0, 0, 0, 1)
}

// 鼠标从物理坐标转换到游戏逻辑坐标
func (g *Game) updateMouse() {
	// 限制比例，不要出现黑边(letterbox)
	g.mouseButtons = sdl.GetMouseState(&g.mousePosition[0], &g.mousePosition[1])
	// 获取当前物理窗口大小
	w, h := int32(0), int32(0)
	sdl.GetWindowSize(g.sdlWindow, &w, &h)
	sdl.RenderCoordinatesFromWindow(g.sdlRenderer, g.mousePosition[0], g.mousePosition[1], &g.mousePosition[0], &g.mousePosition[1])
}
