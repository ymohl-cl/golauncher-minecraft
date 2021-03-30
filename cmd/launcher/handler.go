package main

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/sqweek/dialog"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/ymohl-cl/go-ui/widget"
)

// LoadPathInstallation open a dialog to select a new path to locate minecraft installation path
func (u *ui) LoadPathInstallation() {
	var err error
	var dir string

	if dir, err = dialog.Directory().Title("minecraft folder").Browse(); err != nil {
		fmt.Printf("error to update the installation path: %s\n", err.Error())
	}
	u.pathInstallation = dir
	u.ReloadButton()
}

// MuteUnmuteMusic stop the music if it is playing or run it else
func (u *ui) MuteUnmuteMusic() {
	if u.musicUI.mute.State() != widget.StateOff {
		u.musicUI.mute.SetState(widget.StateOff)
		u.musicUI.unmute.SetState(widget.StateBase)
		u.musicUI.audio.SetStateToPlay(widget.StateOff)
	} else {
		u.musicUI.mute.SetState(widget.StateBase)
		u.musicUI.unmute.SetState(widget.StateOff)
		u.musicUI.audio.SetStateToPlay(widget.StateBase)
	}
}

// PlayButton to run Minecraft
func (u *ui) PlayButton() {
	var err error

	if err = u.installer.Install(); err != nil {
		fmt.Printf("error to install mycube: %s\n", err.Error())
	}
	cmd := exec.Command("open", "/Applications/Minecraft.app")
	err = cmd.Start()
	if err != nil {
		fmt.Printf("error to run minecraft: %s\n", err.Error())
	}
	sdl.PushEvent(&sdl.QuitEvent{Type: sdl.QUIT, Timestamp: uint32(time.Now().Unix())})
}

// ReloadButton : _
func (u *ui) ReloadButton() {
	var err error

	if err = u.installer.UpdateLocation(u.pathInstallation); err != nil {
		fmt.Printf("error to update the installation path: %s\n", err.Error())
	}
	u.updateInfo()

	// context main button
	if err = u.updateTextButton(); err != nil {
		fmt.Printf("error to update the main button: %s", err.Error())
	}
}
