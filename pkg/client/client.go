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

	//Set timeout
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 100,
	}

	resp, err := client.Get(url)

	if err != nil {
		return 0, "", err
	}

	// Cleanup
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, "", err
	}

	return resp.StatusCode, string(body), nil
}
