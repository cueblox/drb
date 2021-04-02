package vercel

import (
	_ "embed"
	"fmt"
	"os"
	"path"

	"github.com/devrel-blox/drb/hosting"
)

func init() {
	vp := &VercelProvider{
		internalName:        "vercel",
		internalDescription: "Vercel express api-only provider",
	}
	hosting.Register(vp.Name(), vp)
}

type VercelProvider struct {
	internalName        string
	internalDescription string
}

func (p *VercelProvider) Name() string {
	return p.internalName
}
func (p *VercelProvider) Description() string {
	return p.internalDescription
}
func (p *VercelProvider) Install() error {
	// make api directory
	root, err := os.Getwd()
	if err != nil {
		return err
	}
	api := path.Join(root, "api")
	err = os.MkdirAll(api, 0755)
	if err != nil {
		return err
	}
	// create index.js in api directory

	index := path.Join(api, "index.js")
	err = hosting.CreateFileWithContents(index, indexjs)
	if err != nil {
		return err
	}

	// create package.json
	pkg := path.Join(root, "package.json")
	err = hosting.CreateFileWithContents(pkg, packagejson)
	if err != nil {
		return err
	}
	// create vercel.json
	vc := path.Join(root, "vercel.json")
	err = hosting.CreateFileWithContents(vc, verceljson)
	if err != nil {
		return err
	}
	fmt.Println("Vercel provider installed.")
	fmt.Println("Run `npm install` to install dependencies.")

	return nil
}

//go:embed index.js
var indexjs string

//go:embed vercel.json
var verceljson string

//go:embed package.json
var packagejson string
