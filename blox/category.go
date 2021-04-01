package blox

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"cuelang.org/go/cuego"
	"github.com/devrel-blox/drb/cueutils"
	"github.com/goccy/go-yaml"
)

type Category struct {
	BaseModel `json:",omitempty"`

	Title       string `json:"title"`
	Description string `json:"description"`

	CoverImage string `json:"cover_image"`
	ShareImage string `json:"share_image"`
}

//go:embed category.cue
var CategoryCue string

func CategoryFromYAML(path string) (Category, error) {
	err := cuego.Constrain(&Category{}, CategoryCue)
	if err != nil {
		return Category{}, cueutils.UsefulError(err)
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return Category{}, cueutils.UsefulError(err)
	}

	var category Category

	err = yaml.Unmarshal(bytes, &category)
	if err != nil {
		return Category{}, cueutils.UsefulError(err)
	}

	err = cuego.Complete(&category)
	if err != nil {
		return Category{}, cueutils.UsefulError(err)
	}

	ext := filepath.Ext(path)
	slug := strings.Replace(filepath.Base(path), ext, "", -1)

	category.ID = slug
	fmt.Printf("Category '%s' validated successfully\n", category.ID)

	return category, nil
}
