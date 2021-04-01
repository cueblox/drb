package blox

import (
	"fmt"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/encoding/gocode/gocodec"
	"cuelang.org/go/encoding/yaml"
	"github.com/davecgh/go-spew/spew"
)

type Article struct {
	BaseModel `json:",omitempty"`

	Title       string `json:"title"`
	Excerpt     string `json:"excerpt"`
	PublishDate string `json:"publish_date"`
	EditDate    string `json:"edit_date"`

	Draft    bool `json:"draft"`
	Featured bool `json:"featured"`

	CoverImage string `json:"cover_image"`
	ShareImage string `json:"share_image"`

	CategoryID string   `json:"category_id"`
	Tags       []string `json:"tags"`

	ProfileID string `json:"profile_id"`
}

const ArticleCUE = `title: string
excerpt?: string
draft?: bool
featured?: bool
cover_image?: string
share_image?: string
profile_id?: string
category_id?: string
body_raw?: string
body?: string
id?: string

`

/*

excerpt?: string
draft?: bool
featured?: bool
cover_image?: string
share_image?: string
category_id? string
profile_id?: string
*/
func ArticleFromYAML(path string) (Article, error) {
	var cueRuntime cue.Runtime
	articleInstance, err := cueRuntime.Compile("article", ArticleCUE)

	if err != nil {
		fmt.Println("Error compiling CUE ")
		return Article{}, err
	}

	valueInstance, err := yaml.Decode(&cueRuntime, path, nil)
	if err != nil {

		fmt.Println("Error decoding yaml ")
		return Article{}, fmt.Errorf("parse YAML file error: %w", err)
	}

	merged := cue.Merge(articleInstance, valueInstance)
	err = merged.Value().Validate()
	if err != nil {

		fmt.Println("Error merging ")
		return Article{}, fmt.Errorf("validation error: %w", err)
	}

	var article Article
	codec := gocodec.New(&cueRuntime, &gocodec.Config{})

	err = codec.Encode(merged.Value(), &article)

	if err != nil {
		return Article{}, fmt.Errorf("encoding error: %w", err)
	}
	spew.Println(article)
	ext := filepath.Ext(path)
	slug := strings.Replace(filepath.Base(path), ext, "", -1)

	article.ID = slug

	fmt.Printf("Article '%s' validated successfully\n", article.ID)

	return article, nil
}
