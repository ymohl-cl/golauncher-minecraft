package main

import (
	"fmt"

	"github.com/sqweek/dialog"
	gamebuilder "github.com/ymohl-cl/game-builder"
)

func main() {
	var ui UI
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

	r := driver.Renderer().Driver()
	if ui, err = NewUI(c, r); err != nil {
		panic(err)
	}

	s := driver.Script()
	if err = s.AddScene("launcher", ui); err != nil {
		panic(err)
	}
	if err = driver.Run("launcher"); err != nil {
		panic(err)
	}

	directory, err := dialog.Directory().Title("Load images").Browse()
	if err != nil {
		panic(err)
	}
	fmt.Printf("directory find: %s\n", directory)
}
