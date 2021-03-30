package main

import (
	goui "github.com/ymohl-cl/go-ui"
)

func main() {
	var ui UI
	var driver goui.GoUI
	var s goui.Scene
	var err error
	var c Config

	if c, err = NewConfig(); err != nil {
		panic(err)
	}
	if ui, err = NewUI(c); err != nil {
		panic(err)
	}

	if driver, err = goui.New(goui.ConfigUI{
		Window: goui.Window{
			Title:  c.UI.Window.Title,
			Width:  c.UI.Window.Width,
			Height: c.UI.Window.Height,
		},
	}); err != nil {
		panic(err)
	}
	defer driver.Close()

	if s, err = goui.NewScene(); err != nil {
		panic(err)
	}
	defer s.Close()
	if err = ui.Build(s); err != nil {
		panic(err)
	}

	if err = driver.Run(s); err != nil {
		panic(err)
	}
}
