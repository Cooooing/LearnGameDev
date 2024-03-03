package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

var running = false
var isPause = false
var fps float64
var width int32 = 960
var height int32 = 600
var cellWidth int32 = 8
var cellHeight int32 = 8
var cellStatus []int

const NANOSECOND = 1000000000 // 纳秒 10^9
const FRAMERATE = 30          // 帧速率 每秒钟刷新的图片的帧数
const TITLE = "Game Dev"

var window *sdl.Window

var renderer *sdl.Renderer

func main() {

	a := make([]int, 10)
	for i := 0; i < 10; i++ {
		fmt.Print(a[i], ",")
	}
	fmt.Println()
	b := make([]int, 10)
	a = append(a, 1)
	a = append(a, 2)
	a = append(a, 2)
	for i := 0; i < 10; i++ {
		b[i] = i
		fmt.Print(b[i], ",")
	}
	a = b
	for i := 0; i < 10; i++ {
		fmt.Print(a[i], ",")
	}
	//var err error
	//if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
	//	panic(err)
	//}
	//defer sdl.Quit()
	//
	//window, err = sdl.CreateWindow(TITLE, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_SHOWN)
	//if err != nil {
	//	panic(err)
	//}
	//defer func(window *sdl.Window) {
	//	err := window.Destroy()
	//	if err != nil {
	//		panic(err)
	//	}
	//}(window)
	//
	//renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_SOFTWARE)
	//err = renderer.SetDrawColor(0, 0, 0, 255)
	//err = renderer.Clear()
	//if err != nil {
	//	panic(err)
	//}
	//start()
}

func start() {
	running = true
	run()
}

func run() {
	frames := 0.0                     // 每秒帧数 fps
	var frameCounter int64 = 0        // 帧计数器，记录累计展示的多个帧的时间，大于等于 1s 时刷新 fps 的值，并重置为 0
	lastTime := time.Now().UnixNano() // 上次循环结束时间
	unprocessedTime := 0.0            // 未通过的时间。累加多次循环的时间，包括一次 更新和渲染 及 多次等待。直至达到预设一帧的时间，才进入循环，允许下次 更新和渲染
	frameTime := 1.0 / FRAMERATE      // 每帧的时间
	//initCell()
	for running {
		isRender := false // 是否 更新和渲染
		startTIme := time.Now().UnixNano()
		passedTime := startTIme - lastTime // 上一次更新和渲染所用时间
		lastTime = startTIme

		unprocessedTime += float64(passedTime) / NANOSECOND // 累加上次循环消耗的时间
		frameCounter += passedTime                          // 累加上次循环消耗的时间
		input()
		for unprocessedTime > frameTime {
			isRender = true
			unprocessedTime = 0

			if frameCounter >= NANOSECOND {
				fps = frames
				window.SetTitle(fmt.Sprintf("%s | FPS: %.2f", TITLE, fps))
				frames = 0
				frameCounter = 0
			}
		}

		// 更新和渲染
		if isRender {
			//render()
			//update()
			frames++
		}
	}
	//cleanup()
}
func input() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.KeyboardEvent:
			keyboardEvent := event.(*sdl.KeyboardEvent)
			if keyboardEvent.Type == sdl.KEYDOWN && keyboardEvent.Keysym.Sym == sdl.K_SPACE {
				isPause = !isPause
				println("按下空格")
			}
			break
		case *sdl.QuitEvent:
			println("Quit")
			running = false
			break
		}
	}
}
