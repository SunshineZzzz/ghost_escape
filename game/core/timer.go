package core

// 定时器
type Timer struct {
	// 继承基础对象
	Object
	// 计数器
	timer float32
	// 时间间隔
	Interval float32
	// 是否超时
	IsTimeOut bool
}

// 创建timer
func AddTimerChild(parent IObject, interval float32) *Timer {
	t := &Timer{}
	t.SetInterval(interval)
	t.SetActive(false)
	if parent != nil {
		parent.AddChild(t)
	}
	return t
}

// 更新
func (t *Timer) Update(dt float32) {
	t.timer += dt
	if t.timer >= t.Interval {
		t.timer = 0.0
		t.IsTimeOut = true
	}
}

// 非接口实现

// 开始定时器
func (t *Timer) Start() {
	t.SetActive(true)
}

// 结束
func (t *Timer) Stop() {
	t.SetActive(false)
}

// 是否超时
func (t *Timer) TimeOut() bool {
	if t.IsTimeOut {
		t.IsTimeOut = false
		return true
	}
	return false
}

// 获取当前进程
func (t *Timer) GetProcess() float32 {
	return t.timer / t.Interval
}

// 设置时间间隔
func (t *Timer) SetInterval(interval float32) {
	t.Interval = interval
}

// 获取时间间隔
func (t *Timer) GetInterval() float32 {
	return t.Interval
}
