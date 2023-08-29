package client

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"time"
)

// Fetch fetches a URL and returns information about the response.
func Fetch(url string) (int, string, error) {

	//resp, err := http.Get(url)

	//var netClient = &http.Client{
	//	Timeout: time.Second * 100,
	//}
	//resp, err := netClient.Get(url)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 10,
	}
	defer client.CloseIdleConnections() // Cleanup any old connections
	resp, err := client.Get(url)

	if err != nil {
		return 0, "", err
	}

	defer resp.Body.Close()
	tr.CloseIdleConnections()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, "", err
	}

	return resp.StatusCode, string(body), nil
}
