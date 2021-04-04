package blox

import (
	"testing"

	"cuelang.org/go/cue"
)

func TestGetSchemaVersion(t *testing.T) {
	// Can we get the version from the schema's metadata?
	cueWithVersion := `{
		_schema: {
			version: "v1"
		}
}`

	var cueRuntime cue.Runtime
	cueInstance, err := cueRuntime.Compile("schemaWithVersion", cueWithVersion)

	version, err := GetSchemaVersion(cueInstance.Value())
	if nil != err {
		t.FailNow()
	}

	if version != "v1" {
		t.FailNow()
	}

	// Do we error if there is no eversion?
	cueWithoutVersion := `{
		_schema: {
			name: "v1"
		}
}`

	cueInstance, err = cueRuntime.Compile("schemaWithoutVersion", cueWithoutVersion)

	version, err = GetSchemaVersion(cueInstance.Value())
	if nil == err {
		t.FailNow()
	}
}

func TestGetSchemaV1Metadata(t *testing.T) {
	// Can we get the version from the schema's metadata?
	schemaV1Metadata := `{
		_schema: {
			namespace: "devrel-blox.com"
			name: "profile"
		}
}`

	var cueRuntime cue.Runtime
	cueInstance, err := cueRuntime.Compile("schemaWithVersion", schemaV1Metadata)

	metadata, err := GetSchemaV1Metadata(cueInstance.Value())
	if nil != err {
		t.FailNow()
	}

	if "devrel-blox.com" != metadata.Namespace {
		t.FailNow()
	}

	if "profile" != metadata.Name {
		t.FailNow()
	}

	invalidSchemaV1Metadata := `{
	_schema: {
		namespace: 123
	}
}`

	cueInstance, err = cueRuntime.Compile("invalidSchema", invalidSchemaV1Metadata)

	metadata, err = GetSchemaV1Metadata(cueInstance.Value())
	if nil == err {
		t.FailNow()
	}
}
