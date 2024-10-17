package home

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"learn-game-dev/global"
	"time"
)

type Home struct {
	title     string
	isRunning bool
}

func NewHome() *Home {
	return &Home{title: "Game Dev Home", isRunning: false}
}

func (g *Home) Start() {
	g.isRunning = true
	g.Run()
}

func (g *Home) Stop() {
	g.isRunning = false
}

func (g *Home) SetRunning(isRun bool) {
	g.isRunning = isRun
}

func (g *Home) Run() {
	frames := 0.0                                // 每秒帧数 fps
	var frameCounter int64 = 0                   // 帧计数器，记录累计展示的多个帧的时间，大于等于 1s 时刷新 fps 的值，并重置为 0
	lastTime := time.Now().UnixNano()            // 上次循环结束时间
	unprocessedTime := 0.0                       // 未通过的时间。累加多次循环的时间，包括一次 更新和渲染 及 多次等待。直至达到预设一帧的时间，才进入循环，允许下次 更新和渲染
	frameTime := 1.0 / float64(global.Framerate) // 每帧的时间
	fmt.Println(frameTime)
	for g.isRunning {
		isRender := false // 是否 更新和渲染
		startTIme := time.Now().UnixNano()
		passedTime := startTIme - lastTime // 上一次更新和渲染所用时间
		lastTime = startTIme
		unprocessedTime += float64(passedTime) / global.NANOSECOND // 累加上次循环消耗的时间
		frameCounter += passedTime                                 // 累加上次循环消耗的时间
		g.Input()
		for unprocessedTime > frameTime {
			isRender = true
			unprocessedTime = 0

			if frameCounter >= global.NANOSECOND {
				global.Fps = frames
				global.Window.SetTitle(fmt.Sprintf("%s | FPS: %.2f", g.title, global.Fps))
				frames = 0
				frameCounter = 0
			}
		}

		// 更新和渲染
		if isRender {
			g.Update()
			g.Render()
			frames++
		}

		if global.NextGame != "Home" {
			global.GameList[global.NextGame].Start()
		}
		//fmt.Println("home ")
	}
	g.Cleanup()
}

func (g *Home) Input() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.KeyboardEvent:
			keyboardEvent := event.(*sdl.KeyboardEvent)
			if keyboardEvent.Type == sdl.KEYDOWN && keyboardEvent.Keysym.Sym == sdl.K_SPACE {
				global.IsPause = !global.IsPause
				println("按下空格")
			} else if keyboardEvent.Type == sdl.KEYDOWN && keyboardEvent.Keysym.Sym == sdl.K_a {
				fmt.Println("按下a")
				global.ChangeGame("ConwayLifeGame")
			}
			break
		case *sdl.QuitEvent:
			println("Quit")
			g.isRunning = false
			break
		}
	}
}

func (g *Home) Update() {

}

func (g *Home) Render() {
	var err error
	err = global.Renderer.SetDrawColor(0, 0, 0, 0)
	err = global.Renderer.Clear()
	rect := &sdl.Rect{X: 10, Y: 10, W: 100, H: 100}
	err = global.Renderer.SetDrawColor(255, 255, 255, 0)
	err = global.Renderer.DrawRect(rect)
	global.Renderer.Present()

	if err != nil {
		panic(err)
	}
}

func (g *Home) Cleanup() {

}
