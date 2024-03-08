package global

import (
	"github.com/veandco/go-sdl2/sdl"
)

var Window *sdl.Window
var Renderer *sdl.Renderer

const NANOSECOND = 1000000000 // 纳秒 10^9
var Framerate = 60            // 帧速率 每秒钟刷新的图片的帧数

var Title = "Game Dev"
var Fps float64
var Width int32 = 960
var Height int32 = 640

var IsPause = false
var NextGame = "Home"

var ShowGame Game
var GameList = make(map[string]Game, 10)
