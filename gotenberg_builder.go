package gotenberg

import (
	"io"
	"net/http"
)

type ClientBuilder struct {
	baseURL    string
	httpClient *http.Client
}

func NewClientBuilder(client *http.Client, baseURL string) *ClientBuilder {
	return &ClientBuilder{
		baseURL:    baseURL,
		httpClient: client,
	}
}

func (cb *ClientBuilder) Build() *Client {
	return NewClient(cb.httpClient, cb.baseURL)
}

func (cb *ClientBuilder) ConvertHTML() *HTMLConversionBuilder {
	client := cb.Build()
	return &HTMLConversionBuilder{
		client: client,
		config: &clientOptions{
			Files: make(map[string]io.Reader),
		},
	}
}

func (cb *ClientBuilder) ConvertURL() *URLConversionBuilder {
	client := cb.Build()
	return &URLConversionBuilder{
		client: client,
		config: &clientOptions{
			Files: make(map[string]io.Reader),
		},
	}
}
