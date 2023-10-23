package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"sync"
	"time"
)

type Message struct {
	Message   string
	Timestamp time.Time
}

func getContent(ctx context.Context, fileName string) (<-chan string, <-chan error) {
	errChan := make(chan error)

	resChan := make(chan string, 3)

	go func() {
		file, err := os.Open(fileName)
		if err != nil {
			errChan <- err
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
					close(resChan)
					return
				}
				if err != nil {
					errChan <- err
					continue
				}
				var message Message
				err = json.Unmarshal([]byte(line), &message)
				if err != nil {
					errChan <- err
					continue
				}
				timeToSleep := rand.Intn(5) * int(time.Second)
				time.Sleep(time.Duration(timeToSleep))
				resChan <- "[" + message.Timestamp.Format(time.DateTime) + "]: " + message.Message
			}

		}
	}()

	return resChan, errChan
}

func saveContent(ctx context.Context, contentChan <-chan string, fileName string) <-chan error {
	errChan := make(chan error)
	go func() {
		file, err := os.Create(fileName)
		if err != nil {
			errChan <- err
		}
		defer file.Close()
		defer close(errChan)
		w := bufio.NewWriter(file)
		for {
			select {
			case <-ctx.Done():
				w.Flush()
				return
			case chunk, ok := <-contentChan:
				if !ok {
					w.Flush()
					return
				}
				if _, err := w.WriteString(chunk); err != nil {
					errChan <- err
				}
			}
		}

	}()

	return errChan
}

func main() {
	fileNames := map[string]string{"file1.txt": "out1.txt", "file2.txt": "out2.txt", "file3.txt": "out3.txt"}

	ctx, cancelFunc := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	for inFileName, outFileName := range fileNames {
		inChan, errChan := getContent(ctx, inFileName)
		outChan := make(chan string)
		saveErrChan := saveContent(ctx, outChan, outFileName)
		wg.Add(1)
		go func() {
			outClosed := false
			for {
				select {
				case res, more := <-inChan:
					if more {
						fmt.Println(res)
						outChan <- res + "\n"
						continue
					}
					if !outClosed {
						outClosed = true
						close(outChan)
					}

				case err := <-errChan:
					fmt.Println("ERROR READING", err)
					cancelFunc()
					wg.Done()
					return
				case err, more := <-saveErrChan:
					if !more {
						wg.Done()
						return
					}
					fmt.Println("ERROR SAVING", err)
					cancelFunc()
				}
			}

		}()
	}
	wg.Wait()
}
