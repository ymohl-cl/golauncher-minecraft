package main

import (
	"fmt"

	"github.com/ymohl-cl/game-builder/objects"
	"github.com/ymohl-cl/game-builder/objects/block"
	"github.com/ymohl-cl/game-builder/objects/button"
	"github.com/ymohl-cl/game-builder/objects/image"
	"github.com/ymohl-cl/game-builder/objects/text"
)

// Build the scene
func (u *ui) Build() error {
	var err error

	if err = u.addBackground(); err != nil {
		return err
	}
	if err = u.addButtonInstall(); err != nil {
		return err
	}
	return nil
}

func (u *ui) addBackground() error {
	// background
	var i *image.Image
	var err error

	backgroundPath := fmt.Sprintf("%s/%s", u.config.ResourcePath, u.config.UI.Background)
	i = image.New(backgroundPath, 0, 0, u.config.Window.Width, u.config.Window.Height)
	if err = i.Init(u.renderer); err != nil {
		return err
	}
	u.layers[layerBackground] = append(u.layers[layerBackground], i)

	return nil
}

func (u *ui) addButtonInstall() error {
	var b *button.Button
	var err error

	if b, err = u.createDefaultButton("INSTALL"); err != nil {
		return err
	}
	if err = b.Init(u.renderer); err != nil {
		return err
	}
	u.layers[layer1] = append(u.layers[layer1], b)
	return nil
}

func (u *ui) createDefaultButton(str string) (*button.Button, error) {
	var bl *block.Block
	var txt *text.Text
	var b *button.Button
	var err error

	w := int32(500)
	h := int32(50)
	x := u.config.Window.Width/int32(2) - (w / int32(2))
	y := u.config.Window.Height/int32(3) - (h / int32(2))

	// color orange dark: rgb(255,140,0) 255 opacity
	// color grey to click: rgb(128,128,128) 255 opacity
	// color grey to over: rgb(42, 42, 42) 255 opacity
	if bl, err = block.New(block.Filled); err != nil {
		return nil, err
	}
	// Set style fix and basic
	if err = bl.SetVariantStyle(255, 140, 0, 255, objects.SFix, objects.SBasic); err != nil {
		return nil, err
	}
	// Set style click
	if err = bl.SetVariantStyle(128, 128, 128, 255, objects.SClick); err != nil {
		return nil, err
	}
	// set style over
	if err = bl.SetVariantStyle(42, 42, 42, 255, objects.SOver); err != nil {
		return nil, err
	}
	// Setters...
	bl.UpdatePosition(x, y)
	bl.UpdateSize(w, h)

	x = x + w/int32(2)
	y = y + h/int32(2)
	fontPath := fmt.Sprintf("%s/%s", u.config.ResourcePath, u.config.UI.Font)
	if txt, err = text.New(str, 42, fontPath, x, y); err != nil {
		return nil, err
	}
	// text white 255, 255, 255
	txt.SetVariantStyle(255, 255, 255, 255, objects.SFix)
	txt.SetVariantUnderStyle(255, 255, 255, 255, objects.SFix)
	//	txt.SetUnderPosition(x-10, y-10)

	b = button.New(bl, nil, txt)
	b.SetAction(u.actionInstall)
	return b, nil
}
