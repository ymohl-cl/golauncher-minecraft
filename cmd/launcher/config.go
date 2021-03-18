package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const (
	configFile = "config.json"
)

// ConfigUI describe the scene ui
type ConfigUI struct {
	Background string `json:"background"`
	Font       string `json:"font"`
}

// Window configuration
type Window struct {
	Width  int32  `json:"width"`
	Height int32  `json:"height"`
	Title  string `json:"title"`
}

// Config data about ui
type Config struct {
	UI           ConfigUI `json:"ui`
	Window       Window   `json:"window"`
	ResourcePath string   `json:"resources_path"`
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
