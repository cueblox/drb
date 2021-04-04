package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/goccy/go-yaml"
)

var Models []Model

func Init() {

	cfg, err := Load()
	if err != nil {
		panic(err)
	}
	for _, s := range cfg.SchemaList {
		fmt.Println(s.Name)
	}
}

type BloxConfig struct {
	Base             string    `json:"base"`
	Source           string    `json:"source"`
	Destination      string    `json:"destination"`
	Templates        string    `json:"templates"`
	Static           string    `json:"static"`
	Schemas          string    `json:"schemas"`
	DefaultExtension string    `json:"default_extension"`
	SchemaList       []*Schema `json:"remote_schemas"`
}

func (c *BloxConfig) AddSchema(s *Schema) error {
	fmt.Println(c.SchemaList)
	for _, m := range s.Models {
		m.Schema = s.Name
	}
	c.SchemaList = append(c.SchemaList, s)

	err := s.Save()
	if err != nil {
		return err
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	configFile := path.Join(wd, "blox.yaml")
	f, err := os.Create(configFile)

	if err != nil {
		return err
	}
	defer f.Close()
	return c.Write(f)
}

type Schema struct {
	Namespace string   `json:"namespace"`
	Name      string   `json:"name"`
	Sources   []string `json:"sources"`
	Models    []*Model `json:"models"`
}

func (s *Schema) Save() error {
	cfg, err := Load()
	if err != nil {
		return err
	}

	// save the schema file
	schemaPath := path.Join(cfg.Base, cfg.Schemas)
	schemaFile := path.Join(schemaPath, s.Name+".json")
	err = os.MkdirAll(schemaPath, 0755)
	if err != nil {
		return err
	}
	bb, err := json.Marshal(s)
	err = os.WriteFile(schemaFile, bb, 0755)
	if err != nil {
		return err
	}
	schemaDir := path.Join(schemaPath, s.Name)
	err = os.MkdirAll(schemaDir, 0755)
	if err != nil {
		return err
	}
	for _, m := range s.Models {
		schemaFilePath := path.Join(schemaDir, m.Name+".json")
		bb, err := json.Marshal(m)
		if err != nil {
			return err
		}
		err = os.WriteFile(schemaFilePath, bb, 0755)
		if err != nil {
			return err
		}
	}

	return err
}

type Model struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Folder     string `json:"folder"`
	ForeignKey string `json:"foreign_key"`
	Cue        string `json:"cue"`
	Schema     string `json:"schema"`
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
