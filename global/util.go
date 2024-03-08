package global

func ChangeGame(gameName string) {
	ShowGame.SetRunning(false)
	NextGame = gameName
	ShowGame = GameList[gameName]
	ShowGame.SetRunning(true)
}

func BackHome() {
	ChangeGame("Home")
}
