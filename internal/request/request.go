package request

import (
	"context"
	"io"
	"net/http"
)

func MakeHTTPGet(u string, headers map[string]string) (int, []byte, error) {
	client := new(http.Client)

	req, _ := http.NewRequest(http.MethodGet, u, nil)

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req.WithContext(context.Background()))
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, err
	}

	return resp.StatusCode, body, nil
}
