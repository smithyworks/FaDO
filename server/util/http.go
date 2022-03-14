package util

import (
	"net/http"
)

func HttpDelete(url string) (resp *http.Response, err error) {
	// Request (DELETE http://www.example.com/bucket/sample)

    // Create client
    client := &http.Client{}

    // Create request
    req, err := http.NewRequest("DELETE", url, nil)
    if err != nil {
        return resp, ProcessErr(err)
    }

    // Fetch Request
    if resp, err = client.Do(req); err != nil {
        return resp, ProcessErr(err)
    }

	return
}
