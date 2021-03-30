package main

import (
	"fmt"
	"path/filepath"

	goui "github.com/ymohl-cl/go-ui"
	"github.com/ymohl-cl/go-ui/widget"
)

func (u *ui) Build(s goui.Scene) error {
	var err error

	if err = u.initComponents(); err != nil {
		return err
	}
	if err = u.setBackground(s); err != nil {
		return err
	}
	if err = u.setInfoBlocks(s); err != nil {
		return err
	}
	if err = u.setTitle(s); err != nil {
		return err
	}
	for i, b := range u.config.Elements.Buttons {
		if err = u.setButton(b, int32(i), s); err != nil {
			return err
		}
	}
	if err = u.setPlayButton(s); err != nil {
		return err
	}
	if err = u.setReloadButton(s); err != nil {
		return err
	}
	return nil
}

func (u *ui) initComponents() error {
	var err error
	var font widget.Font

	fontPath := filepath.Join(fmt.Sprintf("%s/%s", u.config.Resources, u.config.Elements.Font))
	if font, err = widget.NewFont(fontPath, 14); err != nil {
		return err
	}
	if u.minecraftUI.path, err = widget.NewText("", font); err != nil {
		return err
	}
	if u.forgeUI.path, err = widget.NewText("", font); err != nil {
		return err
	}

	iconValid := filepath.Join(fmt.Sprintf("%s/%s", u.config.Resources, "check.png"))
	iconInvalid := filepath.Join(fmt.Sprintf("%s/%s", u.config.Resources, "close.png"))
	if u.minecraftUI.valid, err = widget.NewPicture(iconValid); err != nil {
		return err
	}
	if u.minecraftUI.invalid, err = widget.NewPicture(iconInvalid); err != nil {
		return err
	}
	if u.forgeUI.valid, err = widget.NewPicture(iconValid); err != nil {
		return err
	}
	if u.forgeUI.invalid, err = widget.NewPicture(iconInvalid); err != nil {
		return err
	}

	// create mute button
	pathMute := filepath.Join(fmt.Sprintf("%s/%s", u.config.Resources, "mute.png"))
	pathUmute := filepath.Join(fmt.Sprintf("%s/%s", u.config.Resources, "unmute.png"))
	if u.musicUI.mute, err = widget.NewPicture(pathMute); err != nil {
		return err
	}
	if u.musicUI.unmute, err = widget.NewPicture(pathUmute); err != nil {
		return err
	}
	pathMusic := filepath.Join(fmt.Sprintf("%s/%s", u.config.Resources, u.config.Elements.Music))
	if u.musicUI.audio, err = widget.NewAudio(pathMusic, true); err != nil {
		panic(err)
	}

	// create play button
	if font, err = widget.NewFont(fontPath, 38); err != nil {
		return err
	}
	if u.buttonUI.text, err = widget.NewText("", font); err != nil {
		return err
	}
	u.buttonUI.block = widget.NewRect(0, 0)
	return nil
}

func (u *ui) setBackground(s goui.Scene) error {
	var err error
	var p *widget.Picture

	// background image
	path := filepath.Join(fmt.Sprintf("%s/%s", u.config.Resources, u.config.Elements.Background))
	if p, err = widget.NewPicture(path); err != nil {
		return err
	}
	p.SetSize(u.config.Window.Width, u.config.Window.Height)
	s.AddWidget(p, goui.Layer(layerBackground))

	// background music
	s.AddWidget(u.musicUI.audio, goui.Layer(layerBackground))

	// button music on / off
	ellipse := widget.NewEllipse(60, 60)
	ellipse.SetColor(colorBackgroundText)
	ellipse.SetHoverColor(widget.ColorWhite)
	ellipse.SetPosition(50, 50)
	ellipse.SetAction(widget.NewFuncAction(u.MuteUnmuteMusic))
	u.musicUI.mute.SetPosition(30, 30)
	u.musicUI.mute.SetSize(40, 40)
	u.musicUI.unmute.SetPosition(30, 30)
	u.musicUI.unmute.SetSize(40, 40)
	u.musicUI.unmute.SetState(widget.StateOff)

	s.AddWidget(ellipse, goui.Layer(layer1))
	s.AddWidget(u.musicUI.mute, goui.Layer(layer2))
	s.AddWidget(u.musicUI.unmute, goui.Layer(layer2))
	return nil
}

func (u *ui) setInfoBlocks(s goui.Scene) error {
	var err error

	w := int32(1000)
	h := int32(20)
	x := int32(100)
	y := int32(650)

	// block info to path minecraft
	block, border := u.infoBlock(w, h, x, y)
	s.AddWidget(block, goui.Layer(layer1))
	s.AddWidget(border, goui.Layer(layer2))

	// text and icons configuration about path minecraft
	u.setInfoMinecraftPictures(h-2, h-2, x+w-h, y+1)
	u.minecraftUI.path.SetSize(w, h-2)
	u.minecraftUI.path.SetPosition(x+h, y+1)
	u.minecraftUI.path.SetColor(widget.ColorWhite)
	s.AddWidget(u.minecraftUI.path, goui.Layer(layer3))
	s.AddWidget(u.minecraftUI.valid, goui.Layer(layer3))
	s.AddWidget(u.minecraftUI.invalid, goui.Layer(layer3))

	// button to load a new minecraft path installation
	widthButton := int32(60)
	var iconMore *widget.Picture
	path := filepath.Join(fmt.Sprintf("%s/%s", u.config.Resources, "more.png"))
	if iconMore, err = widget.NewPicture(path); err != nil {
		return err
	}
	iconMore.SetHeightRatioWidth(h - 2)
	widthIcon, _ := iconMore.Size()
	iconX := x + w + 20 + ((widthButton - widthIcon) / 2)
	iconMore.SetPosition(iconX, y+1)

	block, border = u.infoBlock(widthButton, h, x+w+20, y)
	block.SetHoverColor(colorGreyHover)
	block.SetActionColor(colorGreyClick)
	block.SetAction(widget.NewFuncAction(u.LoadPathInstallation))

	s.AddWidget(block, goui.Layer(layer1))
	s.AddWidget(border, goui.Layer(layer2))
	s.AddWidget(iconMore, goui.Layer(layer3))

	// block info to path forge
	y += 25
	block, border = u.infoBlock(w, h, x, y)
	s.AddWidget(block, goui.Layer(layer1))
	s.AddWidget(border, goui.Layer(layer2))

	// text and icons configuration about path forge
	u.setInfoForgePictures(h-2, h-2, x+w-h, y+1)
	u.forgeUI.path.SetSize(w, h-2)
	u.forgeUI.path.SetPosition(x+h, y+1)
	u.forgeUI.path.SetColor(widget.ColorWhite)
	s.AddWidget(u.forgeUI.path, goui.Layer(layer3))
	s.AddWidget(u.forgeUI.valid, goui.Layer(layer3))
	s.AddWidget(u.forgeUI.invalid, goui.Layer(layer3))

	// update content info block
	u.updateInfo()
	return nil
}

func (u *ui) infoBlock(w, h, x, y int32) (*widget.Rect, *widget.Rect) {
	var block *widget.Rect
	var border *widget.Rect

	block = widget.NewRect(w, h)
	block.SetPosition(x, y)
	block.SetColor(colorBackgroundText)

	border = widget.NewRect(w, h)
	border.SetPosition(x, y)
	border.SetStyle(widget.RectStyleBorder)
	border.SetColor(widget.ColorWhite)

	return block, border
}

func (u *ui) setInfoMinecraftPictures(w, h, x, y int32) {
	u.minecraftUI.valid.SetSize(w, h)
	u.minecraftUI.valid.SetPosition(x, y)
	u.minecraftUI.invalid.SetSize(w, h)
	u.minecraftUI.invalid.SetPosition(x, y)
}

func (u *ui) setInfoForgePictures(w, h, x, y int32) {
	u.forgeUI.valid.SetSize(w, h)
	u.forgeUI.valid.SetPosition(x, y)
	u.forgeUI.invalid.SetSize(w, h)
	u.forgeUI.invalid.SetPosition(x, y)
}

func (u *ui) updateInfo() {
	if u.installer.MinecraftPath().Valid {
		u.minecraftUI.path.Set(u.installer.MinecraftPath().Value)
		u.minecraftUI.valid.SetState(widget.StateBase)
		u.minecraftUI.invalid.SetState(widget.StateOff)
	} else {
		u.minecraftUI.path.Set(fmt.Sprintf("minecraft version %s not found, please install it", u.installer.GameVersion()))
		u.minecraftUI.valid.SetState(widget.StateOff)
		u.minecraftUI.invalid.SetState(widget.StateBase)
	}

	// info about forge installation
	if u.installer.ForgePath().Valid {
		u.forgeUI.path.Set(u.installer.ForgePath().Value)
		u.forgeUI.valid.SetState(widget.StateBase)
		u.forgeUI.invalid.SetState(widget.StateOff)
	} else {
		u.forgeUI.path.Set(fmt.Sprintf("forge version %s not found, please install it", u.installer.ForgeVersion()))
		u.forgeUI.valid.SetState(widget.StateOff)
		u.forgeUI.invalid.SetState(widget.StateBase)
	}
}

func (u *ui) setTitle(s goui.Scene) error {
	var err error
	var font widget.Font
	var title *widget.Text
	var clickAudio *widget.Audio

	if u.config.Elements.Title.Content == "" {
		fmt.Printf("Warning: no title")
		return nil
	}
	path := filepath.Join(fmt.Sprintf("%s/%s", u.config.Resources, u.config.Elements.Title.Font))
	if font, err = widget.NewFont(path, 42); err != nil {
		return err
	}
	if title, err = widget.NewText(u.config.Elements.Title.Content, font); err != nil {
		return err
	}
	titleW, titleH := title.Size()

	x := u.config.Window.Width/2 - titleW/2
	y := u.config.Window.Height/3 - titleH/2
	title.SetPosition(x, y)
	title.SetSize(titleW, titleH)
	title.SetColor(widget.ColorWhite)
	title.SetHoverColor(colorGreyClick)

	block := widget.NewRect(u.config.Window.Width, titleH*2)
	block.SetColor(colorBackgroundText)
	block.SetPosition(0, y-titleH/2)

	s.AddWidget(block, goui.Layer(layer1))
	s.AddWidget(title, goui.Layer(layer2))

	path = filepath.Join(fmt.Sprintf("%s/%s", u.config.Resources, u.config.Elements.ClickAudio))
	if clickAudio, err = widget.NewAudio(path, false); err != nil {
		return err
	}
	clickAudio.SetStateToPlay(widget.StateHover)
	clickAudio.SetPosition(x, y)
	clickAudio.SetSize(titleW, titleH)
	s.AddWidget(clickAudio, goui.Layer(0))

	return nil
}

func (u *ui) setButton(b Button, n int32, s goui.Scene) error {
	var err error
	var block *widget.Rect
	var logo *widget.Picture
	var clickAudio *widget.Audio

	path := filepath.Join(fmt.Sprintf("%s/%s", u.config.Resources, b.Logo))
	if logo, err = widget.NewPicture(path); err != nil {
		return err
	}
	logo.SetWidthRatioHeight(50)
	w, h := logo.Size()
	x := u.config.Window.Width - w - 10
	y := int32(10) + (50 * n) + (10 * n)
	logo.SetPosition(x, y)

	block = widget.NewRect(w, h)
	block.SetPosition(x, y)
	block.SetColor(widget.ColorNoColor)
	block.SetHoverColor(colorBackgroundText)
	block.SetAction(widget.NewLinkAction(b.URL))

	s.AddWidget(logo, goui.Layer(layer2))
	s.AddWidget(block, goui.Layer(layer3))

	path = filepath.Join(fmt.Sprintf("%s/%s", u.config.Resources, u.config.Elements.ClickAudio))
	if clickAudio, err = widget.NewAudio(path, false); err != nil {
		return err
	}
	clickAudio.SetStateToPlay(widget.StateAction)
	clickAudio.SetPosition(x, y)
	clickAudio.SetSize(w, h)
	s.AddWidget(clickAudio, goui.Layer(0))

	return nil
}

func (u *ui) setPlayButton(s goui.Scene) error {
	var border *widget.Rect
	var clickAudio *widget.Audio
	var err error

	w := int32(600)
	h := int32(60)
	x := u.config.Window.Width/2 - w/2
	y := int32(550)
	u.buttonUI.block.SetSize(w, h)
	u.buttonUI.block.SetColor(colorBackgroundGrey)
	u.buttonUI.block.SetHoverColor(colorGreyHover)
	u.buttonUI.block.SetPosition(x, y)
	s.AddWidget(u.buttonUI.block, goui.Layer(layer1))

	border = widget.NewRect(w, h)
	border.SetColor(widget.ColorWhite)
	border.SetStyle(widget.RectStyleBorder)
	border.SetPosition(x, y)
	s.AddWidget(border, goui.Layer(layer2))

	u.buttonUI.text.SetSize(w, h)
	u.buttonUI.text.SetColor(widget.ColorWhite)
	if err = u.updateTextButton(); err != nil {
		return err
	}
	s.AddWidget(u.buttonUI.text, goui.Layer(layer3))

	path := filepath.Join(fmt.Sprintf("%s/%s", u.config.Resources, u.config.Elements.ClickAudio))
	if clickAudio, err = widget.NewAudio(path, false); err != nil {
		return err
	}
	clickAudio.SetStateToPlay(widget.StateAction)
	clickAudio.SetPosition(x, y)
	clickAudio.SetSize(w, h)
	s.AddWidget(clickAudio, goui.Layer(layerBackground))

	return nil
}

func (u *ui) updateTextButton() error {
	if !u.installer.MinecraftPath().Valid {
		u.buttonUI.text.Set("Telecharger Minecraft")
		u.buttonUI.block.SetAction(widget.NewLinkAction("https://www.minecraft.net/en-us/get-minecraft/"))
	} else if !u.installer.ForgePath().Valid {
		u.buttonUI.text.Set("Telecharger Forge")
		link := fmt.Sprintf("http://files.minecraftforge.net/maven/net/minecraftforge/forge/index_%s.html", u.installer.ForgeVersion())
		u.buttonUI.block.SetAction(widget.NewLinkAction(link))
	} else {
		u.buttonUI.text.Set("Jouer !")
		u.buttonUI.block.SetAction(widget.NewFuncAction(u.PlayButton))
	}

	// block size
	w, h := u.buttonUI.text.Size()
	// str size
	textW, textH, err := u.buttonUI.text.SizeSTR()
	if err != nil {
		return err
	}
	buttonPos := u.buttonUI.block.Position()
	x := buttonPos.X + w/2 - textW/2
	y := buttonPos.Y + h/2 - textH/2 - 1
	u.buttonUI.text.SetPosition(x, y)
	return nil
}

func (u *ui) setReloadButton(s goui.Scene) error {
	var err error
	var img *widget.Picture
	var block *widget.Rect
	var border *widget.Rect
	var clickAudio *widget.Audio

	mainButtonPos := u.buttonUI.block.Position()
	mainButtonW, mainButtonH := u.buttonUI.block.Size()
	x := mainButtonPos.X + 20 + mainButtonW
	y := mainButtonPos.Y

	path := filepath.Join(fmt.Sprintf("%s/%s", u.config.Resources, "reload.png"))
	if img, err = widget.NewPicture(path); err != nil {
		return err
	}
	img.SetSize(mainButtonH-10, mainButtonH-10)
	img.SetPosition(x+5, y+5)
	s.AddWidget(img, goui.Layer(layer2))

	block = widget.NewRect(mainButtonH, mainButtonH)
	block.SetPosition(x, y)
	block.SetColor(colorBackgroundText)
	block.SetHoverColor(colorGreyHover)
	block.SetActionColor(colorGreyClick)
	block.SetAction(widget.NewFuncAction(u.ReloadButton))
	s.AddWidget(block, goui.Layer(layer1))

	border = widget.NewRect(mainButtonH, mainButtonH)
	border.SetPosition(x, y)
	border.SetColor(widget.ColorWhite)
	border.SetStyle(widget.RectStyleBorder)
	s.AddWidget(border, goui.Layer(layer2))

	path = filepath.Join(fmt.Sprintf("%s/%s", u.config.Resources, u.config.Elements.ClickAudio))
	if clickAudio, err = widget.NewAudio(path, false); err != nil {
		return err
	}
	clickAudio.SetStateToPlay(widget.StateAction)
	clickAudio.SetPosition(x, y)
	clickAudio.SetSize(mainButtonH, mainButtonH)
	s.AddWidget(clickAudio, goui.Layer(layerBackground))
	return nil
}
