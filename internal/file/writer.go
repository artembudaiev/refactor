package file

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"os"
)

// in - processed messages with filename, out - filename and line of the file
func WriteContent(ctx context.Context, in <-chan string, out chan<- MessageWithFilename) {
	//errChan := make(chan error)

	// todo collect and sort out by timestamp, write
	return resChan, errChan
