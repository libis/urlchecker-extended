package client

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"time"
)

// Fetch fetches a URL and returns information about the response.
func Fetch(url string) (int, string, error) {

	// This allows invalid certs
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 10,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, "", err
	}

	resp, err := client.Do(req)

	if err != nil {
		return 0, "", err
	}

	// Cleanup
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, "", err
	}

	client.CloseIdleConnections() //Cleanup connection
	return resp.StatusCode, string(body), nil
}
