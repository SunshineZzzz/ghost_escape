package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unsafe"

	"github.com/SunshineZzzz/purego-sdl3/sdl"
	"github.com/hajimehoshi/go-mp3"
	"github.com/jfreymuth/oggvorbis"
)

// 声音抽象
type ISound interface {
	// 播放声音
	Play()
	// 暂停声音
	Pause()
	// 停止声音
	Stop()
	// 设置是否循环播放
	SetLoop(loop bool)
	// 关闭声音
	Close()
}

// 创建声音
func NewSound(soundFilePath string) (ISound, error) {
	extWithDot := filepath.Ext(soundFilePath)
	ext := strings.ToLower(extWithDot[1:])

	switch ext {
	case "ogg":
		return newOggSound(soundFilePath)
	case "wav":
		return newWavSound(soundFilePath)
	case "mp3":
		return newMp3Sound(soundFilePath)
	default:
		return nil, fmt.Errorf("unsupported audio file format: %s", extWithDot)
	}
}

// 全局声音句柄管理
var soundHandles = struct {
	handles map[uint32]ISound
	nextID  uint32
}{
	handles: make(map[uint32]ISound),
	nextID:  1,
}

// 注册声音
func registerSound(sound ISound) uint32 {
	id := soundHandles.nextID
	soundHandles.handles[id] = sound
	soundHandles.nextID++
	return id
}

// 获取声音
func getSound(id uint32) ISound {
	return soundHandles.handles[id]
}

// 注销声音
func unregisterSound(id uint32) {
	delete(soundHandles.handles, id)
}

// OGG格式声音
type oggSound struct {
	// SDL音频流
	stream *sdl.AudioStream
	// ogg音频数据
	audioData []byte
	// 当前播放位置
	dataPos int
	// 正在播放
	isPlaying bool
	// 是否循环播放
	loop bool
	// id
	id uint32
	// 音频规格
	sampleRate int32
	channels   int32
}

var _ ISound = (*oggSound)(nil)

func newOggSound(soundFilePath string) (*oggSound, error) {
	file, err := os.Open(soundFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sound file, %v, %v", soundFilePath, err)
	}
	defer file.Close()

	oggReader, err := oggvorbis.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to create oggvorbis reader, %v, %v", soundFilePath, err)
	}

	pcmData := make([]float32, 1024*1024)
	totalSamples := 0
	for {
		n, err := oggReader.Read(pcmData[totalSamples:])
		if err != nil && err.Error() != "EOF" {
			return nil, fmt.Errorf("failed to read oggvorbis data, %v, %v", soundFilePath, err)
		}
		if n == 0 {
			break
		}
		totalSamples += n
	}

	spec := &sdl.AudioSpec{
		Freq:     int32(oggReader.SampleRate()),
		Channels: int32(oggReader.Channels()),
		Format:   sdl.AudioF32,
	}

	callback := sdl.NewAudioStreamCallback(oggAudioCallback)
	ogg := &oggSound{
		audioData:  Float32ToBytes(pcmData[:totalSamples]),
		dataPos:    0,
		isPlaying:  false,
		loop:       false,
		sampleRate: int32(oggReader.SampleRate()),
		channels:   int32(oggReader.Channels()),
	}

	ogg.id = registerSound(ogg)

	ogg.stream = sdl.OpenAudioDeviceStream(
		sdl.AudioDeviceDefaultPlayback,
		spec,
		callback,
		unsafe.Pointer(uintptr(ogg.id)),
	)

	if ogg.stream == nil {
		return nil, fmt.Errorf("failed to open audio stream: %s", sdl.GetError())
	}

	return ogg, nil
}

// 音频回调函数
func oggAudioCallback(userdata unsafe.Pointer, stream *sdl.AudioStream, additionalAmount, totalAmount int32) {
	id := uint32(uintptr(userdata))
	ogg := getSound(id).(*oggSound)

	// 安全检查
	if ogg == nil || ogg.id != id || !ogg.isPlaying || len(ogg.audioData) == 0 {
		return
	}

	// 计算剩余数据量
	remaining := len(ogg.audioData) - ogg.dataPos
	if remaining <= 0 {
		if ogg.loop {
			ogg.dataPos = 0
			remaining = len(ogg.audioData)
		} else {
			ogg.isPlaying = false
			sdl.AudioStreamDevicePaused(stream)
			return
		}
	}

	// 推送数据到音频流
	neededBytes := int(additionalAmount)
	dataToSend := min(neededBytes, remaining)
	if dataToSend > 0 {
		data := ogg.audioData[ogg.dataPos : ogg.dataPos+dataToSend]
		sdl.PutAudioStreamData(stream, (*uint8)(unsafe.Pointer(&data[0])), int32(dataToSend))
		ogg.dataPos += dataToSend
	}

	// 再次检查循环（如果刚好发送完所有数据）
	if ogg.loop && ogg.dataPos >= len(ogg.audioData) {
		ogg.dataPos = 0
	}
}

// 播放
func (o *oggSound) Play() {
	if o.stream == nil || o.id == 0 {
		return
	}
	o.isPlaying = true
	o.dataPos = 0
	sdl.ResumeAudioStreamDevice(o.stream)
}

// 暂停
func (o *oggSound) Pause() {
	if o.stream == nil || o.id == 0 {
		return
	}
	o.isPlaying = false
	sdl.AudioStreamDevicePaused(o.stream)
}

// 停止
func (o *oggSound) Stop() {
	if o.stream == nil || o.id == 0 {
		return
	}
	o.isPlaying = false
	o.dataPos = 0
	sdl.AudioStreamDevicePaused(o.stream)
	sdl.ClearAudioStream(o.stream)
}

// 设置循环播放
func (o *oggSound) SetLoop(loop bool) {
	o.loop = loop
}

// 关闭播放器，释放资源
func (o *oggSound) Close() {
	if o.stream != nil {
		o.Stop()
		sdl.DestroyAudioStream(o.stream)
		o.stream = nil
	}
	unregisterSound(o.id)
}

// wav格式声音
type wavSound struct {
	// SDL音频流
	stream *sdl.AudioStream
	// wav音频数据
	audioBuf *uint8
	// wav音频数据长度
	audioLen int
	// 当前播放位置
	dataPos int
	// 正在播放
	isPlaying bool
	// 是否循环播放
	loop bool
	// 音频规格
	spec *sdl.AudioSpec
	// id
	id uint32
}

func newWavSound(soundFilePath string) (*wavSound, error) {
	// 打开文件IO流
	ioStream := sdl.IOFromFile(soundFilePath, "rb")
	if ioStream == nil {
		return nil, fmt.Errorf("failed to open WAV file: %s", sdl.GetError())
	}
	// 自动释放了
	// defer sdl.CloseIO(ioStream)

	// 使用SDL直接加载WAV文件
	var audioBuf *uint8
	var audioLen uint32
	spec := &sdl.AudioSpec{}
	// 加载WAV数据
	success := sdl.LoadWAVIO(ioStream, true, spec, &audioBuf, &audioLen)
	if !success {
		return nil, fmt.Errorf("failed to load WAV data: %s", sdl.GetError())
	}

	wav := &wavSound{
		audioBuf:  audioBuf,
		audioLen:  int(audioLen),
		spec:      spec,
		dataPos:   0,
		isPlaying: false,
		loop:      false,
	}

	// 注册WAV播放器
	wav.id = registerSound(wav)

	// 创建音频流
	callback := sdl.NewAudioStreamCallback(wavAudioCallback)
	wav.stream = sdl.OpenAudioDeviceStream(
		sdl.AudioDeviceDefaultPlayback,
		spec,
		callback,
		unsafe.Pointer(uintptr(wav.id)),
	)

	if wav.stream == nil {
		sdl.Free(unsafe.Pointer(audioBuf))
		return nil, fmt.Errorf("failed to open audio stream: %s", sdl.GetError())
	}

	return wav, nil
}

// wav音频回调函数
func wavAudioCallback(userdata unsafe.Pointer, stream *sdl.AudioStream, additionalAmount, totalAmount int32) {
	id := uint32(uintptr(userdata))
	wav := getSound(id).(*wavSound)

	// 安全检查
	if wav == nil || wav.id != id || !wav.isPlaying || wav.audioLen == 0 {
		return
	}

	// 计算剩余数据量
	remaining := wav.audioLen - wav.dataPos
	if remaining <= 0 {
		if wav.loop {
			wav.dataPos = 0
			remaining = wav.audioLen
		} else {
			wav.isPlaying = false
			sdl.AudioStreamDevicePaused(stream)
			return
		}
	}

	// 推送数据到音频流
	neededBytes := int(additionalAmount)
	dataToSend := min(neededBytes, remaining)
	if dataToSend > 0 {
		data := (*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(wav.audioBuf)) + uintptr(wav.dataPos)))
		// wav.audioBuf+wav.dataPos
		sdl.PutAudioStreamData(stream, data, int32(dataToSend))
		wav.dataPos += dataToSend
	}

	// 再次检查循环
	if wav.loop && wav.dataPos >= wav.audioLen {
		wav.dataPos = 0
	}
}

// 播放控制方法（保持不变）
func (w *wavSound) Play() {
	if w.stream == nil || w.id == 0 {
		return
	}
	w.isPlaying = true
	w.dataPos = 0
	sdl.ResumeAudioStreamDevice(w.stream)
}

func (w *wavSound) Pause() {
	if w.stream == nil || w.id == 0 {
		return
	}
	w.isPlaying = false
	sdl.AudioStreamDevicePaused(w.stream)
}

func (w *wavSound) Stop() {
	if w.stream == nil || w.id == 0 {
		return
	}
	w.isPlaying = false
	w.dataPos = 0
	sdl.AudioStreamDevicePaused(w.stream)
	sdl.ClearAudioStream(w.stream)
}

func (w *wavSound) SetLoop(loop bool) {
	w.loop = loop
}

// 关闭播放器，释放资源
func (w *wavSound) Close() {
	if w.stream != nil {
		w.Stop()
		sdl.DestroyAudioStream(w.stream)
		w.stream = nil
	}
	if w.audioLen > 0 {
		sdl.Free(unsafe.Pointer(w.audioBuf))
		w.audioBuf = nil
		w.audioLen = 0
	}
	unregisterSound(w.id)
}

// mp3格式声音
type mp3Sound struct {
	// SDL音频流
	stream *sdl.AudioStream
	// MP3音频数据
	audioData []byte
	// 当前播放位置
	dataPos int
	// 正在播放
	isPlaying bool
	// 是否循环播放
	loop bool
	// id
	id uint32
	// 音频规格
	sampleRate int32
	channels   int32
}

var _ ISound = (*mp3Sound)(nil)

func newMp3Sound(soundFilePath string) (*mp3Sound, error) {
	file, err := os.Open(soundFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sound file, %v, %v", soundFilePath, err)
	}
	defer file.Close()

	d, err := mp3.NewDecoder(file)
	if err != nil {
		return nil, fmt.Errorf("failed to create oggvorbis reader, %v, %v", soundFilePath, err)
	}

	pcmData := make([]byte, 1024*1024)
	totalSamples := 0
	for {
		n, err := d.Read(pcmData[totalSamples:])
		if err != nil && err.Error() != "EOF" {
			return nil, fmt.Errorf("failed to read mp3 data, %v, %v", soundFilePath, err)
		}
		if n == 0 {
			break
		}
		totalSamples += n
	}

	spec := &sdl.AudioSpec{
		Freq:     int32(d.SampleRate()),
		Channels: 2,
		Format:   sdl.AudioS16,
	}

	callback := sdl.NewAudioStreamCallback(mp3AudioCallback)
	mp3 := &mp3Sound{
		audioData:  pcmData,
		dataPos:    0,
		isPlaying:  false,
		loop:       false,
		sampleRate: int32(d.SampleRate()),
		channels:   2,
	}

	mp3.id = registerSound(mp3)

	mp3.stream = sdl.OpenAudioDeviceStream(
		sdl.AudioDeviceDefaultPlayback,
		spec,
		callback,
		unsafe.Pointer(uintptr(mp3.id)),
	)

	if mp3.stream == nil {
		return nil, fmt.Errorf("failed to open audio stream: %s", sdl.GetError())
	}

	return mp3, nil
}

// 音频回调函数
func mp3AudioCallback(userdata unsafe.Pointer, stream *sdl.AudioStream, additionalAmount, totalAmount int32) {
	id := uint32(uintptr(userdata))
	mp3 := getSound(id).(*mp3Sound)

	// 安全检查
	if mp3 == nil || mp3.id != id || !mp3.isPlaying || len(mp3.audioData) == 0 {
		return
	}

	// 计算剩余数据量
	remaining := len(mp3.audioData) - mp3.dataPos
	if remaining <= 0 {
		if mp3.loop {
			mp3.dataPos = 0
			remaining = len(mp3.audioData)
		} else {
			mp3.isPlaying = false
			sdl.AudioStreamDevicePaused(stream)
			return
		}
	}

	// 推送数据到音频流
	neededBytes := int(additionalAmount)
	dataToSend := min(neededBytes, remaining)
	if dataToSend > 0 {
		data := mp3.audioData[mp3.dataPos : mp3.dataPos+dataToSend]
		sdl.PutAudioStreamData(stream, (*uint8)(unsafe.Pointer(&data[0])), int32(dataToSend))
		mp3.dataPos += dataToSend
	}

	// 再次检查循环（如果刚好发送完所有数据）
	if mp3.loop && mp3.dataPos >= len(mp3.audioData) {
		mp3.dataPos = 0
	}
}

// 播放
func (o *mp3Sound) Play() {
	if o.stream == nil || o.id == 0 {
		return
	}
	o.isPlaying = true
	o.dataPos = 0
	sdl.ResumeAudioStreamDevice(o.stream)
}

// 暂停
func (o *mp3Sound) Pause() {
	if o.stream == nil || o.id == 0 {
		return
	}
	o.isPlaying = false
	sdl.AudioStreamDevicePaused(o.stream)
}

// 停止
func (o *mp3Sound) Stop() {
	if o.stream == nil || o.id == 0 {
		return
	}
	o.isPlaying = false
	o.dataPos = 0
	sdl.AudioStreamDevicePaused(o.stream)
	sdl.ClearAudioStream(o.stream)
}

// 设置循环播放
func (o *mp3Sound) SetLoop(loop bool) {
	o.loop = loop
}

// 关闭播放器，释放资源
func (o *mp3Sound) Close() {
	if o.stream != nil {
		o.Stop()
		sdl.DestroyAudioStream(o.stream)
		o.stream = nil
	}
	unregisterSound(o.id)
}
