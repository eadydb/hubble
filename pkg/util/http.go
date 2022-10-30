package util

import (
	"fmt"
	"io"
	"net/http"
)

func Download(url string) ([]byte, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating http request: %w", err)
	}
	// req.Header.Set("User-Agent", version.UserAgentWithClient())
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http %d, error %q", resp.StatusCode, resp.Status)
	}
	return io.ReadAll(resp.Body)
}
