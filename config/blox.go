package config

import (
	"errors"
	"io"
	"os"
	"path"

	"github.com/goccy/go-yaml"
)

// I hate that I had to redefine this, but it's easier than refactoring the entire blox/model/config
// loop atm
type Model struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Folder     string `json:"folder"`
	ForeignKey string `json:"foreignkey"`
	Cue        string `json:"cue"`
}

type BloxConfig struct {
	Base             string  `json:"base"`
	Source           string  `json:"source"`
	Destination      string  `json:"destination"`
	Templates        string  `json:"templates"`
	Static           string  `json:"static"`
	DefaultExtension string  `json:"default_extension"`
	Models           []Model `json:"models"`
}

func (b *BloxConfig) Write(w io.Writer) error {
	bb, err := yaml.Marshal(b)
	if err != nil {
		return err
	}
	_, err = w.Write(bb)
	return err
}

func Load() (*BloxConfig, error) {
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
