package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"math/rand"
	"time"
)

/*
生命游戏中，对于任意细胞，规则如下：

每个细胞有两种状态 - 存活或死亡，每个细胞与以自身为中心的周围八格细胞产生互动（如图，黑色为存活，白色为死亡）
当前细胞为存活状态时，当周围的存活细胞低于2个时（不包含2个），该细胞变成死亡状态。（模拟生命数量稀少）
当前细胞为存活状态时，当周围有2个或3个存活细胞时，该细胞保持原样。
当前细胞为存活状态时，当周围有超过3个存活细胞时，该细胞变成死亡状态。（模拟生命数量过多）
当前细胞为死亡状态时，当周围有3个存活细胞时，该细胞变成存活状态。（模拟繁殖）
可以把最初的细胞结构定义为种子，当所有在种子中的细胞同时被以上规则处理后，可以得到第一代细胞图。按规则继续处理当前的细胞图，可以得到下一代的细胞图，周而复始。
*/

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
const TITLE = "Conway's Game of Life"

var window *sdl.Window

var renderer *sdl.Renderer

func main() {
	var err error
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err = sdl.CreateWindow(TITLE, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer func(window *sdl.Window) {
		err := window.Destroy()
		if err != nil {
			panic(err)
		}
	}(window)

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_SOFTWARE)
	if err != nil {
		panic(err)
	}

	start()
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
	initCell()
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
			render()
			update()
			frames++
		}
	}
	cleanup()
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

func render() {
	var err error
	err = renderer.SetDrawColor(0, 0, 0, 0)
	err = renderer.Clear()
	rects := make([]sdl.Rect, len(cellStatus))
	var i int32
	for i = 0; i < int32(len(cellStatus)); i++ {
		if cellStatus[i] == 1 {
			rects = append(rects, sdl.Rect{X: i * cellWidth % width, Y: i * cellWidth / width * cellHeight, W: cellWidth, H: cellHeight})
		}
	}
	err = renderer.SetDrawColor(255, 255, 255, 0)
	err = renderer.FillRects(rects)
	renderer.Present()

	if err != nil {
		panic(err)
	}
}

func update() {
	if isPause {
		return
	} else {
		updateCell()
	}
}

func cleanup() {

}

var totalNumber int32
var number int32
var widthNumber = width / cellWidth

func initCell() {
	totalNumber = width * height / cellHeight / cellWidth
	number = totalNumber / 10
	cellStatus = make([]int, totalNumber)

	// 初始化 cellStatus
	var i int32
	for i = 0; i < totalNumber; i++ {
		cellStatus[i] = 0
	}
	// 随机生成初始细胞
	tmpNumber := number
	for i = 0; i < tmpNumber; i++ {
		item := rand.Int31n(totalNumber)
		if cellStatus[item] == 1 {
			tmpNumber++
			continue
		} else {
			cellStatus[item] = 1
		}
	}
	fmt.Println("inited")
}

func updateCell() {
	var i int32
	for i = 0; i < totalNumber; i++ {
		count := 0 // 周围存活细胞数量
		if cellStatus[(i-1+totalNumber)%totalNumber] == 1 {
			count = count + 1
		}
		if cellStatus[(i+1+totalNumber)%totalNumber] == 1 {
			count = count + 1
		}
		if cellStatus[(i-widthNumber+totalNumber)%totalNumber] == 1 {
			count = count + 1
		}
		if cellStatus[(i+widthNumber+totalNumber)%totalNumber] == 1 {
			count = count + 1
		}
		if cellStatus[(i-1+widthNumber+totalNumber)%totalNumber] == 1 {
			count = count + 1
		}
		if cellStatus[(i-1-widthNumber+totalNumber)%totalNumber] == 1 {
			count = count + 1
		}
		if cellStatus[(i+1+widthNumber+totalNumber)%totalNumber] == 1 {
			count = count + 1
		}
		if cellStatus[(i+1-widthNumber+totalNumber)%totalNumber] == 1 {
			count = count + 1
		}
		// 更新细胞状态
		newCellMap := make([]int, totalNumber)
		copy(newCellMap, cellStatus)
		if cellStatus[i] == 1 && (count < 2 || count > 3) {
			newCellMap[i] = 0
		}
		if cellStatus[i] == 0 && count == 3 {
			newCellMap[i] = 1
		}
		cellStatus = newCellMap

	}
}
