package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type myNewProvider struct{}

func (p *myNewProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "mynewprovider"
	resp.Version = "0.1.0"
}

func (p *myNewProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

func (p *myNewProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

func (p *myNewProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewTaskResource,
	}
}

func (p *myNewProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

func New() provider.Provider {
	return &myNewProvider{}
}
