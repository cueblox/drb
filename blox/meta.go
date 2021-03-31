package blox

import (
	"errors"
	"os"
	"path"

	"github.com/devrel-blox/drb/config"
	"github.com/spf13/cobra"
)

// Models is used by various commands to determine
// how to perform certain actions based on arguments
// and flags provided. All new types must be
// represented in this slice.
var Models = []Model{
	{
		ID:         "profile",
		Name:       "Profile",
		Folder:     "profiles",
		ForeignKey: "profile_id",
	},
}

type Model struct {
	ID         string
	Name       string
	Folder     string
	ForeignKey string
}

// GetModel finds a Model definition and returns
// it to the caller.
func GetModel(id string) (Model, error) {
	for _, m := range Models {
		if m.ID == id {
			return m, nil
		}
	}
	return Model{}, errors.New("model not found")
}

func (m Model) SourceContentPath() string {
	cfg, err := config.Load()
	cobra.CheckErr(err)
	return path.Join(cfg.Base, cfg.Source, m.Folder)
}
func (m Model) DestinationContentPath() string {
	cfg, err := config.Load()
	cobra.CheckErr(err)
	return path.Join(cfg.Base, cfg.Destination, m.Folder)
}
func (m Model) SourceFilePath(slug string) string {
	cfg, err := config.Load()
	cobra.CheckErr(err)
	fileName := slug + cfg.DefaultExtension

	return path.Join(cfg.Base, cfg.Source, m.Folder, fileName)
}
func (m Model) DestinationFilePath(slug string) string {
	cfg, err := config.Load()
	cobra.CheckErr(err)
	fileName := slug + ".yaml"

	return path.Join(cfg.Base, cfg.Destination, m.Folder, fileName)
}
func (m Model) New(slug string) error {
	err := os.MkdirAll(m.SourceContentPath(), 0744)
	if err != nil {
		return err
	}
	f, err := os.Create(m.SourceFilePath(slug))
	if err != nil {
		return err
	}
	defer f.Close()
	f.WriteString("---")

	return nil
}

// baseModel defines fields used by all drb
// models
type BaseModel struct {
	ID      string `json:"id" yaml:"id"`
	Body    string `json:"body" yaml:"body"`         //filled in by processing
	BodyRaw string `json:"body_raw" yaml:"body_raw"` //filled in by processing
}