package checker

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/libis/urlchecker-extended/pkg/client"
	"github.com/libis/urlchecker-extended/pkg/config"
	"github.com/libis/urlchecker-extended/pkg/slack"
)

var maxErrors uint64 = 5

func XMLWorker(ctx context.Context, cancel context.CancelFunc, urlChan chan string, id int, messager Messager, wg *sync.WaitGroup, sleep time.Duration, errorCount *uint64) {
	defer wg.Done()

	messages := []slack.Message{} // initialize the slice to store messages
	for {
		select {
		case url, ok := <-urlChan:
			if !ok {
				return
			}

			if config.Debug {
				log.Printf("Checking %s...", url)
			}

			status, _, err := client.Fetch(url)
			if err != nil {
				log.Printf("Error: %s\n", err.Error())
			}

			if status != 200 {
				log.Println(status)

				msg := fmt.Sprintf("Invalid HTTP Response Status %d", status)
				messages = append(messages, slack.Message{Status: status, Url: url, Message: msg})
				atomic.AddUint64(errorCount, 1)
			}

			if *errorCount > maxErrors {
				log.Printf("Aborting... error count [%d] is greater than max error count [%d]", *errorCount, maxErrors)
				cancel()
			}

			time.Sleep(sleep * time.Second)
		case <-ctx.Done():
			if len(messages) > 0 {
				messager.SendMessage(messages) // send collected messages
			}
			return
		}

	}
}
