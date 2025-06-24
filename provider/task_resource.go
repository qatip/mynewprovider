package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &taskResource{}

type taskResource struct{}

type taskModel struct {
	ID          types.String `tfsdk:"id"`
	Title       types.String `tfsdk:"title"`
	Description types.String `tfsdk:"description"`
	Completed   types.Bool   `tfsdk:"completed"`
}

func NewTaskResource() resource.Resource {
	return &taskResource{}
}

func (r *taskResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "mynewprovider_task"
}

func (r *taskResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"title": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"completed": schema.BoolAttribute{
				Optional: true,
			},
		},
	}
}

func (r *taskResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data taskModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	task := map[string]interface{}{
		"title":       data.Title.ValueString(),
		"description": data.Description.ValueString(),
		"completed":   data.Completed.ValueBool(),
	}

	body, _ := json.Marshal(task)
	res, err := http.Post("http://localhost:8080/tasks", "application/json", bytes.NewBuffer(body))
	if err != nil {
		resp.Diagnostics.AddError("API Error", err.Error())
		return
	}
	defer res.Body.Close()

	respData, _ := io.ReadAll(res.Body)
	var created map[string]interface{}
	json.Unmarshal(respData, &created)

	data.ID = types.StringValue(fmt.Sprintf("%v", created["id"]))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *taskResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data taskModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := http.Get("http://localhost:8080/tasks/" + data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read Failed", err.Error())
		return
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	var task map[string]interface{}
	json.NewDecoder(res.Body).Decode(&task)

	data.Title = types.StringValue(task["title"].(string))
	data.Description = types.StringValue(task["description"].(string))
	data.Completed = types.BoolValue(task["completed"].(bool))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *taskResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data taskModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	task := map[string]interface{}{
		"title":       data.Title.ValueString(),
		"description": data.Description.ValueString(),
		"completed":   data.Completed.ValueBool(),
	}

	body, _ := json.Marshal(task)
	reqURL := fmt.Sprintf("http://localhost:8080/tasks/%s", data.ID.ValueString())
	reqPut, _ := http.NewRequest(http.MethodPut, reqURL, bytes.NewBuffer(body))
	reqPut.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err := client.Do(reqPut)
	if err != nil {
		resp.Diagnostics.AddError("Update Failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *taskResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data taskModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqURL := fmt.Sprintf("http://localhost:8080/tasks/%s", data.ID.ValueString())
	reqDel, _ := http.NewRequest(http.MethodDelete, reqURL, nil)
	client := &http.Client{}
	_, err := client.Do(reqDel)
	if err != nil {
		resp.Diagnostics.AddError("Delete Failed", err.Error())
	}
}
