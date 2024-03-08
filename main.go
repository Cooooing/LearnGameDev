package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"learn-game-dev/game/conwayLifeGame"
	"learn-game-dev/game/home"
	"learn-game-dev/global"
)

func main() {
	global.ShowGame.Start()
	destroy()
}

func addGameList() {
	global.GameList["Home"] = home.NewHome()
	global.GameList["ConwayLifeGame"] = conwayLifeGame.NewConwayLifeGame()
	global.ShowGame = global.GameList["Home"]
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
