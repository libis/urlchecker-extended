package client

import (
	"io/ioutil"
	"net/http"
)

// Fetch fetches a URL and returns information about the response.
func Fetch(url string) (int, string, error) {

	resp, err := http.Get(url)

	//var netClient = &http.Client{
	//	Timeout: time.Second * 100,
	//}
	//resp, err := netClient.Get(url)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return 1, "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 2, "", err
	}

	return resp.StatusCode, string(body), nil
}
