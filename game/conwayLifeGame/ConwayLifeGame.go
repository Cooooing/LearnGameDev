package conwayLifeGame

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"learn-game-dev/global"
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

type ConwayLifeGame struct {
	title     string
	isRunning bool

	cellWidth     int32
	cellHeight    int32
	cellStatus    []int
	newCellStatus []int

	totalNumber int32
	number      int32
	widthNumber int32
}

func NewConwayLifeGame() *ConwayLifeGame {
	return &ConwayLifeGame{title: "ConwayLifeGame", isRunning: false, cellWidth: 4, cellHeight: 4}
}

func (g *ConwayLifeGame) Start() {
	g.isRunning = true
	g.Run()
}

func (g *ConwayLifeGame) Stop() {
	g.isRunning = false
}

func (g *ConwayLifeGame) SetRunning(isRun bool) {
	g.isRunning = isRun
}

func (g *ConwayLifeGame) Run() {
	frames := 0.0                                // 每秒帧数 fps
	var frameCounter int64 = 0                   // 帧计数器，记录累计展示的多个帧的时间，大于等于 1s 时刷新 fps 的值，并重置为 0
	lastTime := time.Now().UnixNano()            // 上次循环结束时间
	unprocessedTime := 0.0                       // 未通过的时间。累加多次循环的时间，包括一次 更新和渲染 及 多次等待。直至达到预设一帧的时间，才进入循环，允许下次 更新和渲染
	frameTime := 1.0 / float64(global.Framerate) // 每帧的时间
	g.initCell()
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
	}
	g.Cleanup()
}

func (g *ConwayLifeGame) Input() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.KeyboardEvent:
			keyboardEvent := event.(*sdl.KeyboardEvent)
			if keyboardEvent.Type == sdl.KEYDOWN && keyboardEvent.Keysym.Sym == sdl.K_SPACE {
				global.IsPause = !global.IsPause
				println("按下空格")
			} else if keyboardEvent.Type == sdl.KEYDOWN && keyboardEvent.Keysym.Sym == sdl.K_ESCAPE {
				global.BackHome()
			}
			break
		case *sdl.QuitEvent:
			println("Quit")
			g.isRunning = false
			break
		}
	}
}

func (g *ConwayLifeGame) Render() {
	var err error
	err = global.Renderer.SetDrawColor(0, 0, 0, 0)
	err = global.Renderer.Clear()
	rects := make([]sdl.Rect, len(g.cellStatus))
	var i int32
	for i = 0; i < int32(len(g.cellStatus)); i++ {
		if g.cellStatus[i] == 1 {
			rects = append(rects, sdl.Rect{X: i * g.cellWidth % global.Width, Y: i * g.cellWidth / global.Width * g.cellHeight, W: g.cellWidth, H: g.cellHeight})
		}
	}
	err = global.Renderer.SetDrawColor(255, 255, 255, 0)
	err = global.Renderer.FillRects(rects)
	global.Renderer.Present()

	if err != nil {
		panic(err)
	}
}

func (g *ConwayLifeGame) Update() {
	if global.IsPause {
		return
	} else {
		g.updateCell()
	}
}

func (g *ConwayLifeGame) Cleanup() {

}

func (g *ConwayLifeGame) initCell() {
	global.Width = 960
	global.Height = 640
	g.widthNumber = global.Width / g.cellWidth
	g.totalNumber = global.Width * global.Height / g.cellHeight / g.cellWidth
	g.number = g.totalNumber / 16
	g.cellStatus = make([]int, g.totalNumber)
	g.newCellStatus = make([]int, g.totalNumber)
	// 初始化 g.cellStatus
	var i int32
	for i = 0; i < g.totalNumber; i++ {
		g.cellStatus[i] = 0
	}
	// 随机生成初始细胞
	tmpNumber := g.number
	for i = 0; i < tmpNumber; i++ {
		item := rand.Int31n(g.totalNumber)
		if g.cellStatus[item] == 1 {
			tmpNumber++
			continue
		} else {
			g.cellStatus[item] = 1
		}
	}
	fmt.Println("inited")
}

func (g *ConwayLifeGame) updateCell() {
	var i int32
	for i = 0; i < g.totalNumber; i++ {
		count := 0 // 周围存活细胞数量
		if g.cellStatus[(i-1+g.totalNumber)%g.totalNumber] == 1 {
			count = count + 1
		}
		if g.cellStatus[(i+1+g.totalNumber)%g.totalNumber] == 1 {
			count = count + 1
		}
		if g.cellStatus[(i-g.widthNumber+g.totalNumber)%g.totalNumber] == 1 {
			count = count + 1
		}
		if g.cellStatus[(i+g.widthNumber+g.totalNumber)%g.totalNumber] == 1 {
			count = count + 1
		}
		if g.cellStatus[(i-1+g.widthNumber+g.totalNumber)%g.totalNumber] == 1 {
			count = count + 1
		}
		if g.cellStatus[(i-1-g.widthNumber+g.totalNumber)%g.totalNumber] == 1 {
			count = count + 1
		}
		if g.cellStatus[(i+1+g.widthNumber+g.totalNumber)%g.totalNumber] == 1 {
			count = count + 1
		}
		if g.cellStatus[(i+1-g.widthNumber+g.totalNumber)%g.totalNumber] == 1 {
			count = count + 1
		}
		// 更新细胞状态

		if g.cellStatus[i] == 1 {
			if count < 2 || count > 3 {
				g.newCellStatus[i] = 0
			} else {
				g.newCellStatus[i] = 1
			}
		} else if g.cellStatus[i] == 0 {
			if count == 3 {
				g.newCellStatus[i] = 1
			} else {
				g.newCellStatus[i] = 0
			}
		}
		g.cellStatus[i] = g.newCellStatus[i]

	}
	g.cellStatus, g.newCellStatus = g.newCellStatus, g.cellStatus
}
