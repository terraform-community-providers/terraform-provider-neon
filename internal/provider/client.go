package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type authedTransport struct {
	token   string
	wrapped http.RoundTripper
}

func (t *authedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = "https"
	req.URL.Host = "console.neon.tech"
	req.URL.Path = "/api/v2" + req.URL.Path

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+t.token)

	return t.wrapped.RoundTrip(req)
}

func do[O interface{}](client *http.Client, req *http.Request, output *O) (err error) {
	res, err := client.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	if res.StatusCode >= 400 {
		return fmt.Errorf(string(responseBody))
	}

	return json.Unmarshal(responseBody, output)
}

func get[O interface{}](client *http.Client, diagnostics diag.Diagnostics, url string, output *O) (err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		diagnostics.AddError("Provider Error", fmt.Sprintf("Unable to form request, got error: %s", err))
		return nil
	}

	return do(client, req, output)
}

func call[I interface{}, O interface{}](client *http.Client, diagnostics diag.Diagnostics, method string, url string, input I, output *O) (err error) {
	requestBody, err := json.Marshal(input)

	if err != nil {
		diagnostics.AddError("Provider Error", fmt.Sprintf("Unable to marshal JSON, got error: %s", err))
		return nil
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))

	if err != nil {
		diagnostics.AddError("Provider Error", fmt.Sprintf("Unable to form request, got error: %s", err))
		return nil
	}

	return do(client, req, output)
}

func delete(client *http.Client, diagnostics diag.Diagnostics, url string) (res *http.Response, err error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)

	if err != nil {
		diagnostics.AddError("Provider Error", fmt.Sprintf("Unable to form request, got error: %s", err))
		return nil, nil
	}

	return client.Do(req)
}
