package main

import (
	gamebuilder "github.com/ymohl-cl/game-builder"
	"github.com/ymohl-cl/golauncher-minecraft/cmd/launcher/ui"
)

func main() {
	var u ui.UI
	var driver gamebuilder.GameBuilder
	var err error
	var c Config

	if c, err = NewConfig(); err != nil {
		panic(err)
	}
	if driver, err = gamebuilder.New(gamebuilder.ConfigUI{
		Window: gamebuilder.Window{
			Title:  c.Window.Title,
			Width:  c.Window.Width,
			Height: c.Window.Height,
		},
	}); err != nil {
		panic(err)
	}

	if u, err = ui.New(); err != nil {
		panic(err)
	}

	s := driver.Script()
	if err = s.AddScene("launcher", u); err != nil {
		panic(err)
	}
	if err = driver.Run("launcher"); err != nil {
		panic(err)
	}
}
