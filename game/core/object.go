package core

import (
	"container/list"

	"github.com/SunshineZzzz/purego-sdl3/sdl"
)

// 对象抽象
type IObject interface {
	// 初始化
	Init()
	// 处理事件
	HandleEvent(*sdl.Event)
	// 更新
	Update(dt float32)
	// 渲染
	Render()
	// 清理
	Clean()
	// 增加孩子
	AddChild(child IObject)
	// 移除孩子
	RemoveChild(child IObject)
	// 获取对象类型
	GetType() ObjectType
	// 设置对象类型
	SetType(ObjectType)
	// 获取是否活跃
	GetIsActive() bool
	// 设置是否活跃
	SetActive(bool)
	// 获取是否需要移除
	GetNeedRemove() bool
	// 设置是否需要移除
	SetNeedRemove(bool)
	// 安全加入孩子
	SafeAddChild(child IObject)
}

// 基础对象
type Object struct {
	// 自身引用，实现多态，要不然始终调用Object的方法
	Self IObject
	// 对象类型
	ObjectType ObjectType
	// 子对象列表
	Children list.List
	// 等待加入场景的子对象列表
	ChildrenToAdd list.List
	// 是否活跃
	IsActive bool
	// 是否需要移除
	NeedRemove bool
}

var _ IObject = (*Object)(nil)

// 初始化
func (o *Object) Init() {
	o.ObjectType = ObjectTypeNone
	o.Children.Init()
	o.ChildrenToAdd.Init()
	o.IsActive = true
	o.NeedRemove = false
}

// 处理事件
func (o *Object) HandleEvent(event *sdl.Event) {
	for e := o.Children.Front(); e != nil; e = e.Next() {
		if e.Value.(IObject).GetIsActive() {
			e.Value.(IObject).HandleEvent(event)
		}
	}
}

// 更新
func (o *Object) Update(dt float32) {
	for e := o.ChildrenToAdd.Front(); e != nil; {
		next := e.Next()
		if o.Self != nil {
			o.Self.AddChild(e.Value.(IObject))
		} else {
			o.AddChild(e.Value.(IObject))
		}
		o.ChildrenToAdd.Remove(e)
		e = next
	}
	o.ChildrenToAdd.Init()
	for e := o.Children.Front(); e != nil; {
		next := e.Next()
		if e.Value.(IObject).GetNeedRemove() {
			o.Children.Remove(e)
			e.Value.(IObject).SetActive(false)
		}
		if e.Value.(IObject).GetIsActive() {
			e.Value.(IObject).Update(dt)
		}
		e = next
	}
}

// 渲染
func (o *Object) Render() {
	// 渲染子对象
	for e := o.Children.Front(); e != nil; e = e.Next() {
		if e.Value.(IObject).GetIsActive() {
			e.Value.(IObject).Render()
		}
	}
}

// 清理
func (o *Object) Clean() {
	// 清理子对象
	for e := o.Children.Front(); e != nil; e = e.Next() {
		e.Value.(IObject).Clean()
	}
	// 清理等待加入场景的子对象
	for e := o.ChildrenToAdd.Front(); e != nil; e = e.Next() {
		e.Value.(IObject).Clean()
	}
	// 清理自身
	o.ChildrenToAdd.Init()
	o.Children.Init()
}

// 增加孩子
func (o *Object) AddChild(child IObject) {
	o.Children.PushBack(child)
}

// 移除孩子
func (o *Object) RemoveChild(child IObject) {
	for e := o.Children.Front(); e != nil; e = e.Next() {
		if e.Value.(IObject) == child {
			o.Children.Remove(e)
		}
	}
}

// 获取对象类型
func (o *Object) GetType() ObjectType {
	return o.ObjectType
}

// 设置对象类型
func (o *Object) SetType(t ObjectType) {
	o.ObjectType = t
}

// 获取是否活跃
func (o *Object) GetIsActive() bool {
	return o.IsActive
}

// 设置是否活跃
func (o *Object) SetActive(active bool) {
	o.IsActive = active
}

// 获取是否需要移除
func (o *Object) GetNeedRemove() bool {
	return o.NeedRemove
}

// 设置是否需要移除
func (o *Object) SetNeedRemove(needRemove bool) {
	o.NeedRemove = needRemove
}

// 安全加入孩子
func (o *Object) SafeAddChild(child IObject) {
	o.ChildrenToAdd.PushBack(child)
}

// 非接口实现

// 获取游戏实例
func (o *Object) Game() *Game {
	return GetInstance()
}
