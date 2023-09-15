package client

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
)

// Fetch fetches a URL and returns information about the response.
func Fetch(url string) (int, string, error) {

	// Ignore tls verification
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	resp, err := http.Get(url)
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
