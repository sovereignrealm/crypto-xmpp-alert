package services

import (
	"io/ioutil"
	"net/http"
	"strings"
)

type MockHTTPClient struct {
	RespBody []byte
}

func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(strings.NewReader(string(m.RespBody))),
	}, nil
}
