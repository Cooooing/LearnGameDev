package game

import (
	"github.com/veandco/go-sdl2/sdl"
	"learn-game-dev/global"
)

type TestGame struct {
	title string
}

func NewTestGame() TestGame {
	return TestGame{title: "Test Game"}
}

func (g TestGame) Input() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.KeyboardEvent:
			keyboardEvent := event.(*sdl.KeyboardEvent)
			if keyboardEvent.Type == sdl.KEYDOWN && keyboardEvent.Keysym.Sym == sdl.K_SPACE {
				global.IsPause = !global.IsPause
				println("按下空格")
			}
			break
		case *sdl.QuitEvent:
			println("Quit")
			global.Running = false
			break
		}
	}
}

func (g TestGame) Render() {
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

func (g TestGame) Update() {

}

func (g TestGame) Cleanup() {

}
