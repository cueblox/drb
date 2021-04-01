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

//go:embed article.cue
var ArticleCue string

//go:embed first-post.md
var ArticleTemplate string

func ArticleFromYAML(path string) (Article, error) {
	err := cuego.Constrain(&Article{}, ArticleCue)
	if err != nil {
		return Article{}, cueutils.UsefulError(err)
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return Article{}, cueutils.UsefulError(err)
	}

	var article Article

	err = yaml.Unmarshal(bytes, &article)
	if err != nil {
		return Article{}, cueutils.UsefulError(err)
	}

	err = cuego.Complete(&article)
	if err != nil {
		return Article{}, cueutils.UsefulError(err)
	}
	ext := filepath.Ext(path)
	slug := strings.Replace(filepath.Base(path), ext, "", -1)

	article.ID = slug
	fmt.Printf("Article '%s' validated successfully\n", article.ID)

	return article, nil
}
