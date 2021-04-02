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

//go:embed article.cue
var ArticleCue string

//go:embed article.md
var ArticleTemplate string

//go:embed category.cue
var CategoryCue string

//go:embed category.md
var CategoryTemplate string

//go:embed profile.cue
var ProfileCue string

//go:embed profile.md
var ProfileTemplate string

func FromYAML(path string, modelName string, cue string) (map[string]interface{}, error) {
	var model = make(map[string]interface{})

	cuego.DefaultContext = &cuego.Context{}

	err := cuego.Constrain(&model, cue)
	if err != nil {
		return nil, cueutils.UsefulError(err)
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, cueutils.UsefulError(err)
	}

	err = yaml.Unmarshal(bytes, &model)
	if err != nil {
		return nil, cueutils.UsefulError(err)
	}

	err = cuego.Complete(&model)
	if err != nil {
		return nil, cueutils.UsefulError(err)
	}

	ext := filepath.Ext(path)
	slug := strings.Replace(filepath.Base(path), ext, "", -1)

	model["ID"] = slug
	fmt.Printf("Model %s '%s' validated successfully\n", modelName, model["ID"])

	return model, nil
}
