package provider

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func branchList(client *http.Client, diagnostics diag.Diagnostics, projectId string) (BranchListOutput, error) {
	var branches BranchListOutput

	err := get(client, diagnostics, fmt.Sprintf("/projects/%s/branches", projectId), &branches)

	return branches, err
}

func branchUpdate(client *http.Client, diagnostics diag.Diagnostics, projectId string, branchId string, input BranchUpdateInput) (BranchOutput, error) {
	var branch BranchOutput

	err := call(client, diagnostics, http.MethodPatch, fmt.Sprintf("/projects/%s/branches/%s", projectId, branchId), input, &branch)

	return branch, err
}

func endpointList(client *http.Client, diagnostics diag.Diagnostics, projectId string) (EndpointListOutput, error) {
	var endpoints EndpointListOutput

	err := get(client, diagnostics, fmt.Sprintf("/projects/%s/endpoints", projectId), &endpoints)

	return endpoints, err
}

func endpointUpdate(client *http.Client, diagnostics diag.Diagnostics, projectId string, endpointId string, input EndpointUpdateInput) (EndpointOutput, error) {
	var endpoint EndpointOutput

	for {
		err := get(client, diagnostics, fmt.Sprintf("/projects/%s/endpoints/%s", projectId, endpointId), &endpoint)

		if err != nil {
			return endpoint, err
		}

		if endpoint.Endpoint.CurrentState != "init" {
			break
		}

		time.Sleep(time.Second)
	}

	err := call(client, diagnostics, http.MethodPatch, fmt.Sprintf("/projects/%s/endpoints/%s", projectId, endpointId), input, &endpoint)

	return endpoint, err
}
