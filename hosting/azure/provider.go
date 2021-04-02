package azure

import "github.com/devrel-blox/drb/hosting"

func init() {
	vp := &AzureProvider{
		internalName:        "azure",
		internalDescription: "Azure static web apps provider",
	}
	hosting.Register(vp.Name(), vp)
}

type AzureProvider struct {
	internalName        string
	internalDescription string
}

func (p *AzureProvider) Name() string {
	return p.internalName
}
func (p *AzureProvider) Description() string {
	return p.internalDescription
}
func (p *AzureProvider) Install() error {
	return nil
}
