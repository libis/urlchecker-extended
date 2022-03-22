package checker

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/bythepixel/urlchecker/pkg/client"
)

func XMLWorker(ctx context.Context, urlChan chan string, id int, messager Messager, wg *sync.WaitGroup, sleep time.Duration) {
	defer wg.Done()
	for {
		select {
		case url, ok := <-urlChan:
			if !ok {
				return
			}
			// log.Printf("Worker [%d] %s...\n", id, url)
			status, _, err := client.Fetch(url)
			if err != nil {
				log.Printf("Error: %s\n", err.Error())
			}

			if status != 200 {
				log.Println(status)
				msg := fmt.Sprintf("Invalid HTTP Response Status %d", status)
				messager.SendMessage(status, url, msg)
			}

			time.Sleep(sleep)

			// log.Printf("Worker [%d] %s good!\n", id, url)
		case <-ctx.Done():
			return
		}

	}

}
