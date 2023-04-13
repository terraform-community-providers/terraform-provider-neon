package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func delete(client *http.Client, url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)

	return do(client, req, err)
}

func do(client *http.Client, req *http.Request, e error) ([]byte, error) {
	if e != nil {
		return nil, fmt.Errorf("unable to form request, got error: %s", e)
	}

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf(string(responseBody))
	}

	return responseBody, nil
}

func doOut[O interface{}](client *http.Client, req *http.Request, e error, output *O) error {
	body, err := do(client, req, e)

	if err != nil {
		return err
	}

	return json.Unmarshal(body, output)
}

func get[O interface{}](client *http.Client, url string, output *O) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)

	return doOut(client, req, err, output)
}

func call[I interface{}, O interface{}](client *http.Client, method string, url string, input I, output *O) error {
	requestBody, err := json.Marshal(input)

	if err != nil {
		return fmt.Errorf("unable to marshal JSON, got error: %s", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))

	return doOut(client, req, err, output)
}
