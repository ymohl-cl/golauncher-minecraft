package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/ymohl-cl/golauncher-minecraft/pkg/installer"
)

const (
	configFile = "config.json"
)

// Button conf
type Button struct {
	Logo string `json:"logo"`
	URL  string `json:"url"`
}

// Title conf
type Title struct {
	Content string `json:"content"`
	Font    string `json:"font"`
}

// ConfigElement json
type ConfigElement struct {
	Font       string   `json:"font"`
	Background string   `json:"background_image"`
	Music      string   `json:"background_music"`
	ClickAudio string   `json:"click_audio"`
	Buttons    []Button `json:"buttons"`
	Title      Title    `json:"title"`
}

// Window configuration
type Window struct {
	Width  int32  `json:"width"`
	Height int32  `json:"height"`
	Title  string `json:"title"`
}

// ConfigUI describe the scene ui
type ConfigUI struct {
	Window    Window        `json:"window"`
	Resources string        `json:"resources"`
	Elements  ConfigElement `json:"elements"`
}

// Config data about ui
type Config struct {
	UI        ConfigUI         `json:"ui`
	Installer installer.Config `json:"installer"`
}

// NewConfig read the default configuration (config.json) to setup a new ui instance
func NewConfig() (Config, error) {
	var buf []byte
	var err error
	var c Config

	if _, err = os.Stat(configFile); err != nil {
		return Config{}, err
	}
	if buf, err = ioutil.ReadFile(configFile); err != nil {
		return Config{}, err
	}

	if err = json.Unmarshal(buf, &c); err != nil {
		return Config{}, err
	}
	return c, nil
}
