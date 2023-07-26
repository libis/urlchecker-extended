// Eventually I would like to create structs for each type (json, xml) and give
// them methods to crawl the urls they contain. Then I could reduce some of the
// duplicate code.
package checker

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"sync"
	"time"

	"github.com/libis/urlchecker-extended/pkg/client"
	"github.com/libis/urlchecker-extended/pkg/config"
	"github.com/libis/urlchecker-extended/pkg/slack"
)

type Messager interface {
	SendMessage(messages []Message)
}

type HealthCheck struct {
	// Path to check
	Path string `json:"path"`

	// Regex used to check the body of the response
	Regex string `json:"regex"`

	// Status code expected from URL
	Status int `json:"status"`

	// XMLSitemap indicates the Path to check is an XML Sitemap
	XMLSitemap bool `json:"xml_sitemap"`
}

type XMLSitemap struct {
	URL []struct {
		Location string `xml:"loc"`
	} `xml:"url"`
}

func Check(filename, protocol, hostname string, messager Messager, workers int, sleep time.Duration) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err.Error())
	}

	var urls []HealthCheck
	err = json.Unmarshal(bytes, &urls)
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var errorCount uint64 = 0
	messages := []slack.Message{} // initialize the slice to store messages

	for _, check := range urls {
		fmt.Printf(".")
		url := protocol + "://" + hostname + check.Path
		if config.Debug {
			log.Printf("Checking %s...", url)
		}

		status, body, err := client.Fetch(url)
		if err != nil {
			log.Printf("Error: %s\n", err.Error())
		}

		if status != check.Status {
			msg := fmt.Sprintf("Invalid HTTP Response Status %d", status)
			messages = append(messages, slack.Message{Status: status, Url: url, Message: msg})
			continue
		}

		if check.Regex != "" {
			if config.Debug {
				log.Println("Checking regex")
			}
			re, err := regexp.Compile(check.Regex)
			if err != nil {
				log.Fatal(err)
			}

			matches := re.MatchString(body)
			if !matches {
				messages = append(messages, slack.Message{Status: status, Url: url, Message: "HTTP Response Body Error"})
				continue
			}
		}

		if check.XMLSitemap {
			numberOfXMLWorkers := workers
			log.Println("Checking sitemap")
			var sitemapUrls XMLSitemap
			err := xml.Unmarshal([]byte(body), &sitemapUrls)
			if err != nil {
				log.Fatal(err)
			}

			xmlUrlsChan := make(chan string)
			go func() {
				for _, xmlUrl := range sitemapUrls.URL {
					xmlUrlsChan <- xmlUrl.Location
				}

				close(xmlUrlsChan)
			}()

			log.Printf("Starting %d workers...\n", numberOfXMLWorkers)
			var xmlWg sync.WaitGroup
			for x := 0; x < numberOfXMLWorkers; x++ {
				log.Printf("Starting worker # %d\n", x)
				xmlWg.Add(1)
				go XMLWorker(ctx, cancel, xmlUrlsChan, x, messager, &xmlWg, sleep, &errorCount)
			}

			xmlWg.Wait()
		}

		time.Sleep(sleep * time.Second)

		if config.Debug {
			log.Printf("%s Good\n", url)
		}
	}

	if len(messages) > 0 {
		messager.SendMessage(messages) // send collected messages
	}
}
