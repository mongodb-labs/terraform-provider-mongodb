package ssh

import (
	"bufio"
	"bytes"
	"io"
	"log"
)

// OutputBuffer used to send SSH command output to log.Printf
type OutputBuffer struct {
	buffer *bytes.Buffer
}

// WithLogging initializes a new output buffer
func WithLogging() *OutputBuffer {
	buffer := new(bytes.Buffer)
	go LogBufferContents(buffer)
	return &OutputBuffer{buffer: buffer}
}

// Output writes the passed line into a byte buffer for later consumption; respects the terraform.UIOutput interface
func (o *OutputBuffer) Output(line string) {
	if _, err := o.buffer.WriteString(line); err != nil {
		log.Printf("Unexpected error: %v", err)
	}
}

// LogBufferContents waits for data to be sent to the buffer and logs each line until the buffer can't be read from anymore
// if a copyBuffer is passed (not nil) the output will also be copied to it
func LogBufferContents(buffer *bytes.Buffer) {
	reader := bufio.NewReader(buffer)

	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			log.Println("[DEBUG] Finished logging buffer")
			return
		}

		log.Printf("[DEBUG] %s", line)
	}
}
