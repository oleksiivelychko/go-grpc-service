package request

import (
	"fmt"
	"io"
	"net/http"
)

func GET(url string) ([]byte, error) {
	response, err := http.DefaultClient.Get(url)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to get successful request, got status code %d", response.StatusCode)
	}

	if err != nil {
		return []byte{}, err
	}

	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(response.Body)

	bytesArr, err := io.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}

	return bytesArr, nil
}
