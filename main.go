package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"learn-game-dev/game"
	"learn-game-dev/global"
	"time"
)

func main() {
	start()
	destroy()

}

func start() {
	global.Running = true
	run()
}

func addGameList() {
	global.GameList["Test"] = game.NewTestGame()
	global.ShowGame = global.GameList["Test"]
}

func init() {
	var err error
	addGameList()
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	global.Window, err = sdl.CreateWindow(global.Title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, global.Width, global.Height, sdl.WINDOW_SHOWN)
	global.Renderer, err = sdl.CreateRenderer(global.Window, -1, sdl.RENDERER_SOFTWARE)
	if err != nil {
		panic(err)
	}
}

func destroy() {
	err := global.Window.Destroy()
	if err != nil {
		panic(err)
	}
	sdl.Quit()
}

func run() {
	frames := 0.0                                // 每秒帧数 fps
	var frameCounter int64 = 0                   // 帧计数器，记录累计展示的多个帧的时间，大于等于 1s 时刷新 fps 的值，并重置为 0
	lastTime := time.Now().UnixNano()            // 上次循环结束时间
	unprocessedTime := 0.0                       // 未通过的时间。累加多次循环的时间，包括一次 更新和渲染 及 多次等待。直至达到预设一帧的时间，才进入循环，允许下次 更新和渲染
	frameTime := float64(1.0 / global.Framerate) // 每帧的时间
	for global.Running {
		isRender := false // 是否 更新和渲染
		startTIme := time.Now().UnixNano()
		passedTime := startTIme - lastTime // 上一次更新和渲染所用时间
		lastTime = startTIme

		unprocessedTime += float64(passedTime) / global.NANOSECOND // 累加上次循环消耗的时间
		frameCounter += passedTime                                 // 累加上次循环消耗的时间
		global.ShowGame.Input()
		for unprocessedTime > frameTime {
			isRender = true
			unprocessedTime = 0

			if frameCounter >= global.NANOSECOND {
				global.Fps = frames
				global.Window.SetTitle(fmt.Sprintf("%s | FPS: %.2f", global.Title, global.Fps))
				frames = 0
				frameCounter = 0
			}
		}

		// 更新和渲染
		if isRender {
			global.ShowGame.Render()
			global.ShowGame.Update()
			frames++
		}
	}
	global.ShowGame.Cleanup()
}
