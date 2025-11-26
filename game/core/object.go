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
}

// 基础对象
type Object struct {
	// 对象类型
	ObjectType ObjectType
	// 子对象列表
	Children list.List
}

var _ IObject = (*Object)(nil)

// 初始化
func (o *Object) Init() {
	o.ObjectType = ObjectTypeNone
	o.Children.Init()
}

// 处理事件
func (o *Object) HandleEvent(*sdl.Event) {
}

// 更新
func (o *Object) Update(float32) {
}

// 渲染
func (o *Object) Render() {
}

// 清理
func (o *Object) Clean() {
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

// 非接口实现

// 获取游戏实例
func (o *Object) Game() *Game {
	return GetInstance()
}
