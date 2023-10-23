package file

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"os"
	"time"
)

type Message struct {
	Message   string
	Timestamp time.Time
}

type MessageWithFilename struct {
	filename string
	message  Message
}

// in - filenames, out - filename and line of the file
func GetContent(ctx context.Context, in <-chan string, out chan<- MessageWithFilename) {
	//errChan := make(chan error)

	select {
	case <-ctx.Done(): // if cancel() execute
		return
	case filename := <-in:
		file, err := os.Open(filename)
		if err != nil {
			//errChan <- err
		}
		defer file.Close()
		reader := bufio.NewReader(file)
		for {
			select {
			case <-ctx.Done(): // if cancel() execute
				return
			default:
				line, err := reader.ReadString('\n')
				if err == io.EOF {
					return
				}
				if err != nil {
					//errChan <- err
					continue
				}
				var message Message
				err = json.Unmarshal([]byte(line), &message)
				if err != nil {
					//errChan <- err
					continue
				}
				out <- MessageWithFilename{
					filename: filename,
					message:  message,
				}
				//timeToSleep := rand.Intn(5) * int(time.Second)
				//time.Sleep(time.Duration(timeToSleep))
				//resChan <- "[" + message.Timestamp.Format(time.DateTime) + "]: " + message.Message
			}

		}
	}

	return resChan, errChan
}
