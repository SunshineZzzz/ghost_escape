package core

// 物体动画组件抽象
type IObjectAnima interface {
	// 继承基础依附对象接口
	IObjectAffiliate
	// 设置纹理
	SetTexture(*Texture)
	// 获取当前帧
	GetCurrentFrame() float32
	// 设置当前帧
	SetCurrentFrame(float32)
	// 获取总帧数
	GetTotalFrame() float32
	// 设置总帧数
	SetTotalFrame(float32)
	// 获取帧率
	GetFps() float32
	// 设置帧率
	SetFps(float32)
	// 获取动画帧计时器
	GetFrameTimer() float32
	// 设置动画帧计时器
	SetFrameTimer(float32)
	// 获取是否循环播放
	GetLoop() bool
	// 设置是否循环播放
	SetLoop(bool)
	// 获取是否播放完毕
	GetFinish() bool
	// 设置是否播放完毕
	SetFinish(bool)
}
