package file

import (
	"context"
	"math/rand"
	"time"
)

// in - read messages from certain filename, out - processed messages with certain filename
func ProcessContent(ctx context.Context, in <-chan string, out chan<- MessageWithFilename) {
	//errChan := make(chan error)

	select {
	case <-ctx.Done(): // if cancel() execute
		return
	case msgWithFile := <-in:
		// simulate hard work
		timeToSleep := rand.Intn(5) * int(time.Second)
		time.Sleep(time.Duration(timeToSleep))
		out <- msgWithFile
	}

	return resChan, errChan
}
