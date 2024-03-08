package global

type Game interface {

	// 属性

	// 方法

	Start()
	Stop()
	SetRunning(bool)
	Run()
	Input()
	Render()
	Update()
	Cleanup()
}
