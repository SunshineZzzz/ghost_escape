package screen

import (
	"ghost_escape/game/affiliate"
	"ghost_escape/game/core"

	"github.com/SunshineZzzz/purego-sdl3/sdl"
	"github.com/go-gl/mathgl/mgl32"
)

// 技能使用HUD
type HudSkill struct {
	// 继承基础屏幕对象
	core.ObjectScreen
	// 目标角色
	target core.IActor
	// 技能图标
	icon affiliate.ISprite
	// 技能使用恢复百分比
	percent float32
}

// 创建技能使用HUD
func AddHudSkillChild(parent core.IObject, target core.IActor, iconFile string, pos mgl32.Vec2, scale float32, anchor core.AnchorType) *HudSkill {
	hs := &HudSkill{}
	hs.Init()
	hs.target = target
	hs.icon = affiliate.AddSpriteChild(hs, iconFile, scale, anchor)
	hs.SetRenderPosition(pos)
	if parent != nil {
		parent.AddChild(hs)
	}
	return hs
}

// 初始化
func (hs *HudSkill) Init() {
	hs.ObjectScreen.Init()
	hs.percent = 1.0
}

// 渲染
func (hs *HudSkill) Render() {
	// 先绘制浅色背景
	sdl.SetTextureColorModFloat(hs.icon.GetTexture().Texture, 0.3, 0.3, 0.3)
	pos := hs.GetRenderPosition().Add(hs.icon.GetOffset())
	hs.Game().RenderTexture(hs.icon.GetTexture(), pos, hs.icon.GetSize(), mgl32.Vec2{1.0, 1.0})
	sdl.SetTextureColorModFloat(hs.icon.GetTexture().Texture, 1.0, 1.0, 1.0)
	// 再进行正常绘制
	hs.ObjectScreen.Render()
}

// 更新
func (hs *HudSkill) Update(dt float32) {
	hs.ObjectScreen.Update(dt)
	// 更新技能图标
	if hs.target != nil {
		hs.SetPercent(hs.target.GetSkillPercent())
	}
}

// 非接口实现

// 设置百分比
func (hs *HudSkill) SetPercent(percent float32) {
	hs.percent = mgl32.Clamp(percent, 0.0, 1.0)
	if hs.icon != nil {
		hs.icon.SetPercent(mgl32.Vec2{1.0, hs.percent})
	}
}

// 获取百分比
func (hs *HudSkill) GetPercent() float32 {
	return hs.percent
}
