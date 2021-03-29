package config

import (
	"errors"
	"io"
	"os"
	"path"

	"github.com/goccy/go-yaml"
)

type BloxConfig struct {
	Base        string `json:"base"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Templates   string `json:"templates"`
}

func (b *BloxConfig) Write(w io.Writer) error {
	bb, err := yaml.Marshal(b)
	if err != nil {
		return err
	}
	_, err = w.Write(bb)
	return err
}

func LoadConfig() (*BloxConfig, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, errors.New("could not get current directory")
	}
	bb, err := os.ReadFile(path.Join(cwd, "blox.yaml"))
	if err != nil {
		return nil, errors.New("could not load blox.yaml in current directory")
	}
	var cfg BloxConfig
	err = yaml.Unmarshal(bb, &cfg)
	if err != nil {
		return nil, errors.New("could not parse blox.yaml")
	}
	return &cfg, nil
}
