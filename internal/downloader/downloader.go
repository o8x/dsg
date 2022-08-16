package downloader

import (
	"bytes"
	"encoding/base64"
	"io"
	"net/http"
)

func Download(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	bs, err := base64.StdEncoding.DecodeString(string(all))
	if err != nil {
		return nil, err
	}

	return bs, nil
}

func DownAsReader(url string) (io.Reader, error) {
	bs, err := Download(url)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(bs), nil
}
