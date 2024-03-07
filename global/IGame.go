package global

type Game interface {
	Input()
	Render()
	Update()
	Cleanup()
}
