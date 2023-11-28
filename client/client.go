package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type HttpGetter struct {
	BaseURL string
}

func (h *HttpGetter) Get(group string, key string) ([]byte, error) {
	u := fmt.Sprintf("%v%v/%v",
		h.BaseURL,
		url.QueryEscape(group),
		url.QueryEscape(key))

	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned %v", res.Status)
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}

	return bytes, nil
}
