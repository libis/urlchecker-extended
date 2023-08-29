package client

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"time"
)

// Fetch fetches a URL and returns information about the response.
func Fetch(url string) (int, string, error) {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	//client := &http.Client{
	//	Transport: tr,
	//	Timeout:   time.Second * 30,
	//}

	//defer client.CloseIdleConnections() //Cleanup old connections
	//resp, err := client.Get(url)

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 30,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, "", err
	}

	req.Header.Add("Connection", "close")
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

	return resp.StatusCode, string(body), nil
}
