package raw

import (
	"ghost_escape/game/core"
	"math"

	"github.com/SunshineZzzz/purego-sdl3/sdl"
	"github.com/go-gl/mathgl/mgl32"
)

// 视差星空背景
type BgStar struct {
	// 继承基础对象
	core.Object
	// 远星
	starFar []mgl32.Vec2
	// 中星
	starMid []mgl32.Vec2
	// 近星
	starNear []mgl32.Vec2
	// 视差系数，用于计算视差效果，0.0表示没有视差(固定背景)，1.0表示背景跟着相机移动，0~1之间视差
	// 移动得越快，看起来离玩家越近, 视差系数越接近1.0。移动得越慢，看起来离玩家越远，视差系数越接近0.0。
	// 远星视察系数
	parallaxFar float32
	// 中星视察系数
	parallaxMid float32
	// 近星视察系数
	parallaxNear float32
	// 定时器，用于更新颜色平滑变化
	timer float32
	// 每一层星星数量
	num int
	// 远星颜色，随着timer平滑变化
	colorFar sdl.FColor
	// 中星颜色，随着timer平滑变化
	colorMid sdl.FColor
	// 近星颜色，随着timer平滑变化
	colorNear sdl.FColor
}

// 添加视差星空背景
func AddBgStarChild(parent core.IObject, num int, far, mid, near float32) *BgStar {
	bgStar := &BgStar{}
	bgStar.Init()
	bgStar.num = num
	bgStar.parallaxFar = far
	bgStar.parallaxMid = mid
	bgStar.parallaxNear = near
	bgStar.starFar = make([]mgl32.Vec2, 0, num)
	bgStar.starMid = make([]mgl32.Vec2, 0, num)
	bgStar.starNear = make([]mgl32.Vec2, 0, num)

	// 计算出，摄像机从左到右，从上到下的移动的像素值
	extra := bgStar.Game().GetCurrentScene().GetWorldSize().Add(mgl32.Vec2{-bgStar.Game().GetScreenSize().X(), -bgStar.Game().GetScreenSize().Y()})
	for range num {
		// 依次随机生成远星、中星、近星的图片内容
		// 比如，世界大小是100x100，屏幕大小是10x10，那么extra就是90x90。
		// 远星的视差系数是0.1，那么就随机生成一个位置范围是[0.0， 19.0]^2，当摄像机从世界的(0, 0)移动到(100, 100)时，
		// 就会看到屏幕上的星星从(0, 0)移动到(19.0, 19.0)。
		// 其他同理
		bgStar.starFar = append(bgStar.starFar, core.GetInstance().RandVec2(mgl32.Vec2{0.0, 0.0}, bgStar.Game().GetScreenSize().Add(extra.Mul(bgStar.parallaxFar))))
		bgStar.starMid = append(bgStar.starMid, core.GetInstance().RandVec2(mgl32.Vec2{0.0, 0.0}, bgStar.Game().GetScreenSize().Add(extra.Mul(bgStar.parallaxMid))))
		bgStar.starNear = append(bgStar.starNear, core.GetInstance().RandVec2(mgl32.Vec2{0.0, 0.0}, bgStar.Game().GetScreenSize().Add(extra.Mul(bgStar.parallaxNear))))
	}
	if parent != nil {
		parent.AddChild(bgStar)
	}
	return bgStar
}

// 初始化
func (b *BgStar) Init() {
	b.Object.Init()
	b.colorFar = sdl.FColor{R: 0.0, G: 0.0, B: 0.0, A: 1.0}
	b.colorMid = sdl.FColor{R: 0.0, G: 0.0, B: 0.0, A: 1.0}
	b.colorNear = sdl.FColor{R: 0.0, G: 0.0, B: 0.0, A: 1.0}
	b.parallaxFar = 0.2
	b.parallaxMid = 0.5
	b.parallaxNear = 0.7
	b.num = 200
}

// 更新
func (b *BgStar) Update(dt float32) {
	b.timer += dt
	// 颜色平滑变化
	b.colorFar.R = 0.5 + 0.5*float32(math.Sin(float64(b.timer*0.9)))
	b.colorMid.R = 0.5 + 0.5*float32(math.Sin(float64(b.timer*0.8)))
	b.colorNear.R = 0.5 + 0.5*float32(math.Sin(float64(b.timer*0.7)))

	b.colorFar.G = 0.5 + 0.5*float32(math.Sin(float64(b.timer*0.8)))
	b.colorMid.B = 0.5 + 0.5*float32(math.Sin(float64(b.timer*0.7)))
	b.colorNear.G = 0.5 + 0.5*float32(math.Sin(float64(b.timer*0.6)))

	b.colorFar.B = 0.5 + 0.5*float32(math.Sin(float64(b.timer*0.7)))
	b.colorMid.B = 0.5 + 0.5*float32(math.Sin(float64(b.timer*0.6)))
	b.colorNear.B = 0.5 + 0.5*float32(math.Sin(float64(b.timer*0.5)))

	b.colorFar.A = 1.0
	b.colorMid.A = 1.0
	b.colorNear.A = 1.0
}

// 渲染
func (b *BgStar) Render() {
	// 负数原因是，b.starFar中每一个星星的位置要计算到渲染坐标系的位置，可以理解为 远星位置 - 相机位置 * 视差系数 = 渲染坐标系的远星绘制位置
	b.Game().DrawPoints(&b.starFar, b.Game().GetCurrentScene().GetCameraPosition().Mul(b.parallaxFar).Mul(-1.0), b.colorFar)
	b.Game().DrawPoints(&b.starMid, b.Game().GetCurrentScene().GetCameraPosition().Mul(b.parallaxMid).Mul(-1.0), b.colorMid)
	b.Game().DrawPoints(&b.starNear, b.Game().GetCurrentScene().GetCameraPosition().Mul(b.parallaxNear).Mul(-1.0), b.colorNear)
}
