package blox

import (
	"fmt"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/encoding/gocode/gocodec"
	"cuelang.org/go/encoding/yaml"
)

type Article struct {
	baseModel `json:",omitempty"`

	Title     string  `json:"title"`
	ProfileID string  `json:"profile_id"`
	Profile   Profile `json:"author"`
}

const CUE = `title: string
profile_id?: string

`

func ArticleFromYAML(path string) (Article, error) {
	var cueRuntime cue.Runtime
	articleInstance, err := cueRuntime.Compile("article", CUE)

	if err != nil {
		return Article{}, err
	}

	valueInstance, err := yaml.Decode(&cueRuntime, path, nil)
	if err != nil {
		return Article{}, fmt.Errorf("parse YAML file error: %w", err)
	}

	merged := cue.Merge(articleInstance, valueInstance)
	err = merged.Value().Validate()
	if err != nil {
		return Article{}, fmt.Errorf("validation error: %w", err)
	}

	var article Article
	codec := gocodec.New(&cueRuntime, &gocodec.Config{})

	err = codec.Encode(merged.Value(), &article)

	if err != nil {
		return Article{}, fmt.Errorf("encoding error: %w", err)
	}
	ext := filepath.Ext(path)
	slug := strings.Replace(filepath.Base(path), ext, "", -1)

	article.ID = slug

	fmt.Printf("Article '%s' validated successfully\n", article.ID)

	return article, nil
}
