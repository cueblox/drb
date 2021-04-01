package blox

import (
	"fmt"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/encoding/gocode/gocodec"
	"cuelang.org/go/encoding/yaml"
)

type Category struct {
	BaseModel `json:",omitempty"`

	Title       string `json:"title"`
	Description string `json:"description"`

	CoverImage string `json:"cover_image"`
	ShareImage string `json:"share_image"`
}

const CategoryCUE = `title: string
description?: string

cover_image?: string
share_image?: string
`

func CategoryFromYAML(path string) (Category, error) {
	var cueRuntime cue.Runtime
	categoryInstance, err := cueRuntime.Compile("category", CategoryCUE)

	if err != nil {
		return Category{}, err
	}

	valueInstance, err := yaml.Decode(&cueRuntime, path, nil)
	if err != nil {
		return Category{}, fmt.Errorf("parse YAML file error: %w", err)
	}

	merged := cue.Merge(categoryInstance, valueInstance)
	err = merged.Value().Validate()
	if err != nil {
		return Category{}, fmt.Errorf("validation error: %w", err)
	}

	var category Category
	codec := gocodec.New(&cueRuntime, &gocodec.Config{})

	err = codec.Encode(merged.Value(), &category)

	if err != nil {
		return Category{}, fmt.Errorf("encoding error: %w", err)
	}
	ext := filepath.Ext(path)
	slug := strings.Replace(filepath.Base(path), ext, "", -1)

	category.ID = slug

	fmt.Printf("Category '%s' validated successfully\n", category.ID)

	return category, nil
}
