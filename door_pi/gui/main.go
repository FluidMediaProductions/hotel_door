package gui

import (
	"github.com/zserge/webview"
	"time"
)

type App struct {
	Name string `json:"name"`
	Opening bool `json:"opening"`
}

type GUI struct {
	w webview.WebView
	app *App
	appSync func()
}

func NewGui() (*GUI) {
	gui := &GUI{}

	gui.w = webview.New(webview.Settings{
		Resizable: true,
		Width: 500,
		Height: 500,
		Debug: true,
	})

	gui.app = &App{}

	gui.w.Dispatch(func() {
		//w.SetFullscreen(true)
		appSync, err := gui.w.Bind("app", gui.app)
		if err != nil {
			panic(err)
		}
		gui.appSync = appSync

		gui.loadRes()
	})
	return gui
}

func (gui *GUI) Start()  {
	gui.w.Run()
}

func (gui *GUI) SetRoomName(name string) {
	gui.app.Name = name
	gui.w.Dispatch(gui.appSync)
}

func (gui *GUI) SetDoorOpening() {
	gui.app.Opening = true
	gui.w.Dispatch(gui.appSync)

	go func() {
		<-time.After(time.Second * 3)
		gui.app.Opening = false
		gui.w.Dispatch(gui.appSync)
	}()
}