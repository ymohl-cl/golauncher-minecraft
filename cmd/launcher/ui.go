package main

import (
	goui "github.com/ymohl-cl/go-ui"
	"github.com/ymohl-cl/go-ui/widget"
	"github.com/ymohl-cl/golauncher-minecraft/pkg/installer"
)

// Theme color
var (
	colorGreyClick      = widget.Color{Red: 42, Green: 42, Blue: 42, Alpha: 255}
	colorBackgroundGrey = widget.Color{Red: 42, Green: 42, Blue: 42, Alpha: 75}
	colorGreyHover      = widget.Color{Red: 128, Green: 128, Blue: 128, Alpha: 255}
	colorBackgroundText = widget.Color{Red: 255, Green: 255, Blue: 255, Alpha: 75}
)

const (
	// order layers of scene
	layerBackground uint8 = iota
	layer1
	layer2
	layer3
	layer4
)

const (
	defaultIndexUI = 1
)

// UI for the golauncher-minecraft (implement the scene builder from game-builder)
type UI interface {
	Build(s goui.Scene) error
}

type ui struct {
	config    ConfigUI
	installer installer.Installer

	minecraftUI      infoUI
	forgeUI          infoUI
	musicUI          musicUI
	buttonUI         buttonUI
	pathInstallation string
}

type musicUI struct {
	mute   *widget.Picture
	unmute *widget.Picture
	audio  *widget.Audio
}

type infoUI struct {
	path    *widget.Text
	valid   *widget.Picture
	invalid *widget.Picture
}

type buttonUI struct {
	text  *widget.Text
	block *widget.Rect
}

// NewUI return a new instance of ui
func NewUI(c Config) (UI, error) {
	var err error

	u := ui{config: c.UI}
	if u.installer, err = installer.NewWithConfig(c.Installer); err != nil {
		return nil, err
	}

	u.pathInstallation = u.installer.MinecraftPath().Value
	return &u, nil
}
