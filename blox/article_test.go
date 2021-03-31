package blox

import (
	"testing"

	"github.com/go-test/deep"
	"gopkg.in/yaml.v2"
)

func TestArticleDecoding(t *testing.T) {
	expected := Article{
		Title: "Test",
	}

	data := `
title: Test
`
	var article Article

	err := yaml.Unmarshal([]byte(data), &article)

	// spew.Println(profile)

	if nil != err {
		t.Error("Failed to unmarshall YAML")
	}

	if diff := deep.Equal(article, expected); diff != nil {
		t.Error(diff)
	}
}
