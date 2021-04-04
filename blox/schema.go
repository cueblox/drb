package blox

import (
	_ "embed"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
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

//go:embed page.cue
var PageCue string

//go:embed page.md
var PageTemplate string

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

	model["id"] = slug

	return model, nil
}

var GET_VERSION_CUE = `
{
	_schema: {
		version: string
	}
}
`

// GetSchemaVersion will attempt to pull a "version" out of the
// schema's metadata, returning an error if it can't
func GetSchemaVersion(schema cue.Value) (string, error) {
	var cueRuntime cue.Runtime
	cueInstance, err := cueRuntime.Compile("validateVersionCue", GET_VERSION_CUE)

	versionedSchema := schema.Unify(cueInstance.Value())
	if err = versionedSchema.Validate(); err != nil {
		return "", err
	}

	fields, err := versionedSchema.Fields(cue.All())
	if err != nil {
		return "", err
	}

	for fields.Next() {
		if fields.Label() == "_schema" {
			schemaValue := fields.Value()

			versionField, err := schemaValue.LookupField("version")
			if err != nil {
				return "", err
			}

			stringVersion, err := versionField.Value.String()
			if err != nil {
				return "", err
			}

			return stringVersion, nil
		}
	}

	return "", nil
}

// GetSchemaMetadata
type SchemaV1Metadata struct {
	Namespace string
	Name      string
}

var SCHEMA_V1_METADATA = `
{
	_schema: {
		namespace: string
		name: string
	}
}
`

func GetSchemaV1Metadata(schema cue.Value) (SchemaV1Metadata, error) {
	var cueRuntime cue.Runtime
	cueInstance, err := cueRuntime.Compile("schemav1", SCHEMA_V1_METADATA)

	versionedSchema := schema.Unify(cueInstance.Value())
	if err = versionedSchema.Validate(); err != nil {
		return SchemaV1Metadata{}, err
	}

	fields, err := versionedSchema.Fields(cue.All())
	if err != nil {
		return SchemaV1Metadata{}, err
	}

	schemaV1 := SchemaV1Metadata{}

	for fields.Next() {
		if fields.Label() == "_schema" {
			schemaValue := fields.Value()

			err := schemaValue.Decode(&schemaV1)
			if err != nil {
				return SchemaV1Metadata{}, err
			}

			return schemaV1, nil
		}
	}

	return SchemaV1Metadata{}, errors.New("Couldn't get SchemaV1 Metadata")
}

// Given a compiled Cue document, we extract all the models
func GetModels(schema cue.Value) ([]Model, error) {
	version, err := GetSchemaVersion(schema)
	if nil != err {
		return []Model{}, err
	}

	switch version {
	case "v1":
		metadata, err := GetSchemaV1Metadata(schema)
		if nil != err {
			return []Model{}, err
		}
		return getV1Models(metadata, schema)

	default:
		return []Model{}, errors.New("Unexpected Schema Version")
	}

	return []Model{}, errors.New("Failed to extract Models")
}

func getV1Models(metadata SchemaV1Metadata, schema cue.Value) ([]Model, error) {
	fields, err := schema.Fields(cue.Definitions(true))
	if nil != err {
		return []Model{}, err
	}

	var models []Model

	for fields.Next() {
		if !fields.IsDefinition() {
			continue
		}

		model := Model{
			ID:   strings.ToLower(metadata.Name),
			Name: metadata.Name,
			// 🤣
			Folder:     strings.ToLower(metadata.Name) + "s",
			ForeignKey: strings.ToLower(metadata.Name) + "_id",
			Cue:        "",
		}

		models = append(models, model)
	}

	return []Model{}, errors.New("Failed to extract Models")
}
