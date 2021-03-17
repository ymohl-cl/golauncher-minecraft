package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const (
	configFile = "config.json"
)

// Window configuration
type Window struct {
	Title  string `json:"title"`
	Width  int32  `json:"width"`
	Height int32  `json:"height"`
}

// Config data about ui
type Config struct {
	Window Window `json:"window"`
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
