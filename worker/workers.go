package worker

import (
	"fmt"
	"sync"
	"worker-poll/models"
)

// StartWorker - запускает нового воркера
func StartWorker(channels *models.Channels, cancelChan chan struct{}, wg *sync.WaitGroup, id int) {
	for {
		select {
		case <-channels.Interrupt.Done():
			wg.Done()
			fmt.Printf("Worker %d received interrupt\n", id)
			return
		case <-cancelChan:
			close(cancelChan)
			fmt.Printf("Worker %d received cancel\n", id)
			return
		case value, ok := <-channels.InputData:
			if ok {
				fmt.Printf("Worker %d received message: %s\n", id, value)
			}
		}
	}
}
