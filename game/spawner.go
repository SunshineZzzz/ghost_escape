package game

import (
	"ghost_escape/game/core"
	"ghost_escape/game/world"
)

// 生成器
type Spawner struct {
	// 继承基础对象
	core.Object
	// 生成数量
	num int
	// 生成间隔
	interval float32
	// 计时器
	timer float32
	// 目标
	target *Player
}

var _ core.IObject = (*Spawner)(nil)

// 初始化
func (s *Spawner) Init() {
	s.Object.Init()
	s.num = 20
	s.timer = 0.0
	s.interval = 3.0
}

// 更新
func (s *Spawner) Update(dt float32) {
	s.timer += dt
	if s.timer >= s.interval {
		s.timer = 0.0
		// s.num = 1
		if s.num > 0 {
			s.Game().PlaySound("assets/sound/silly-ghost-sound-242342.mp3", false)
		}
		for i := 0; i < s.num; i++ {
			pos := core.GetInstance().RandVec2(
				core.GetInstance().GetCurrentScene().GetCameraPosition(),
				core.GetInstance().GetCurrentScene().GetCameraPosition().
					Add(core.GetInstance().GetScreenSize()),
			)
			enemy := CreateEnemy(nil, pos, s.target)
			// 敌人产生是从特效精灵动画结束后产生，所以这里生成特效
			world.AddEffectChild(core.GetInstance().GetCurrentScene(), "assets/effect/184_3.png", enemy.GetPosition(), 1.0, core.AnchorTypeCenter, enemy)
		}
		// s.interval = 1000.0
	}
}

// 非接口实现

// 获取生成数量
func (s *Spawner) GetNum() int {
	return s.num
}

// 设置生成数量
func (s *Spawner) SetNum(num int) {
	s.num = num
}

// 获取生成间隔
func (s *Spawner) GetInterval() float32 {
	return s.interval
}

// 设置生成间隔
func (s *Spawner) SetInterval(interval float32) {
	s.interval = interval
}

// 获取定时器
func (s *Spawner) GetTimer() float32 {
	return s.timer
}

// 设置定时器
func (s *Spawner) SetTimer(timer float32) {
	s.timer = timer
}

// 获取目标
func (s *Spawner) GetTarget() *Player {
	return s.target
}

// 设置目标
func (s *Spawner) SetTarget(target *Player) {
	s.target = target
}
