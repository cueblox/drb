package blox

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/devrel-blox/drb/config"
	"github.com/goccy/go-yaml"
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
		Cue:        ProfileCue,
	},
	{
		ID:         "article",
		Name:       "Article",
		Folder:     "articles",
		ForeignKey: "article_id",
		Cue:        ArticleCue,
	},
	{
		ID:         "category",
		Name:       "Category",
		Folder:     "categories",
		ForeignKey: "category_id",
		Cue:        CategoryCue,
	},
}

type Model struct {
	ID         string
	Name       string
	Folder     string
	ForeignKey string
	Cue        string
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

func (m Model) StaticContentPath() string {
	cfg, err := config.Load()
	cobra.CheckErr(err)
	return path.Join(cfg.Base, cfg.Static)
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

	switch m.ID {
	case "article":
		exampleArticle := Article{
			Title: "My Title",
		}

		bytes, err := yaml.Marshal(exampleArticle)
		if err != nil {
			return err
		}

		f.WriteString("---\n")
		f.Write(bytes)
		f.WriteString("---\n")
	case "category":
		exampleCategory := Category{
			Title: "My Title",
		}

		bytes, err := yaml.Marshal(exampleCategory)
		if err != nil {
			return err
		}

		f.WriteString("---\n")
		f.Write(bytes)
		f.WriteString("---\n")
	case "profile":
		exampleProfile := Profile{
			FirstName: "FirstName",
			LastName:  "LastName",
			Company:   "Company",
			Title:     "Title",
			SocialAccounts: []SocialAccount{
				{
					Network:  "twitter",
					Username: "username",
				},
			},
		}

		bytes, err := yaml.Marshal(exampleProfile)
		if err != nil {
			return err
		}

		f.WriteString("---\n")
		f.Write(bytes)
		f.WriteString("---\n")

	default:
		return fmt.Errorf("generator doesn't support %s yet", m.ID)
	}

	return nil
}

// baseModel defines fields used by all drb
// models
type BaseModel struct {
	ID      string
	Body    string `json:"body"`
	BodyRaw string `json:"body_raw"`
}
