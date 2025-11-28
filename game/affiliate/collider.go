package affiliate

import (
	"ghost_escape/game/core"

	"github.com/go-gl/mathgl/mgl32"
)

// 碰撞器组件
type Collider struct {
	// 继承基础碰撞器组件
	core.ObjectCollider
}

var _ core.IObject = (*Collider)(nil)
var _ core.IObjectCollider = (*Collider)(nil)
var _ core.IObjectAffiliate = (*Collider)(nil)

// 添加碰撞器子对象
func AddColliderChild(parent core.IObjectScreen, size mgl32.Vec2, colliderType core.ColliderType, anchorType core.AnchorType) *Collider {
	if colliderType != core.ColliderTypeCircle {
		return nil
	}

	child := &Collider{}
	child.Init()
	child.SetSize(size)
	child.SetParent(parent)
	child.SetColliderType(colliderType)
	child.SetAnchorType(anchorType)
	parent.AddChild(child)
	return child
}

// 初始化
func (s *Collider) Init() {
	s.ObjectAffiliate.Init()
	s.Type = core.ColliderTypeCircle
}

// 渲染
// 渲染碰撞器
func (s *Collider) Render() {
	s.ObjectAffiliate.Render()
	if s.GetColliderType() == core.ColliderTypeCircle {
		pos := s.Parent.GetRenderPosition().Add(s.Offset)
		// 圆形碰撞器渲染
		s.Game().RenderFillCircle(pos, s.Size, 0.3)
	}
}

// 是否发生碰撞
func (s *Collider) IsColliding(other core.IObjectCollider) bool {
	if other == nil {
		return false
	}
	if s.GetColliderType() == core.ColliderTypeCircle && other.GetColliderType() == core.ColliderTypeCircle {
		// 圆形碰撞器检测
		// s.Parent.GetPosition().Add(s.Offset) -> 圆形碰撞器的绘制起点，一般是左上角
		// s.Parent.GetPosition().Add(s.Offset).Add(s.Size.Mul(0.5)) -> 圆形碰撞器的圆心
		point1 := s.Parent.GetPosition().Add(s.Offset).Add(s.Size.Mul(0.5))
		point2 := other.GetParent().GetPosition().Add(other.GetOffset()).Add(other.GetSize().Mul(0.5))
		// 圆形碰撞器检测，判断圆心距离是否小于半径和，说明发生碰撞
		if point1.Sub(point2).Len() < (s.Size.X()+other.GetSize().X())*0.5 {
			return true
		}
		return false
	}
	// TODO: 其他碰撞器类型
	return false
}

// 非接口实现
