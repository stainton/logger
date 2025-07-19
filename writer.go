package logger

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

func createNewLogFile(dir string) (*os.File, error) {
	prefix := strings.Replace(time.Now().Format("20060102150405.999"), ".", "", -1)
	filename := path.Join(dir, fmt.Sprintf("%s.log", prefix))
	return os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, os.ModePerm)
}

// Writer is a long-running function that writes messages to log files in a specified directory.
// It listens to incoming messages on the 'messages' channel and writes them to a log file.
// If the size of the log file exceeds the specified threshold, it closes the current file,
// compresses it (if enabled), and opens a new log file.
// The function also listens for a cancellation signal on the provided context.
// If the context is canceled, it closes the current log file and sends a nil error to the 'stop' channel.
//
// Parameters:
// - ctx: A context that can be used to cancel the function.
// - dir: The directory where log files will be stored.
// - messages: A channel of strings representing log messages.
// - compressible: A channel to which the name of the log file to be compressed is sent.
// - stop: A channel to which an error (if any) is sent when the function completes.
// - threshold: The maximum size (in bytes) of a log file before it is compressed and a new file is opened.
func Writer(ctx context.Context, dir string, messages <-chan string, compressible chan<- string, stop chan<- error, threshold int) {
	var message string
	size := 0
	fmt.Printf("[%s]writer started... ^^\n", time.Now().String())
	fd, err := createNewLogFile(dir)
	if err != nil {
		stop <- err
		return
	}
	for {
		select {
		case message = <-messages:
			n, err := fd.WriteString(message)
			if err != nil {
				fmt.Printf("write message to file failed: %+v, log is: `%s`\n", err, message)
				continue
			}
			if n != len(message) {
				fmt.Printf("write message to file success: `%s`\n", message)
			}
			size += n
			if size >= threshold {
				if err := fd.Sync(); err != nil {
					fmt.Println("writer flush file failed:", err)
				}
				name := fd.Name()
				fd.Close()
				compressible <- name
				size = 0
				fd, err = createNewLogFile(dir)
				if err != nil {
					stop <- err
					return
				}
			}
		case <-ctx.Done():
			fmt.Println("DEBUG ctx被cancel了")
			err := fd.Sync()
			fd.Close()
			stop <- err
			fmt.Println("DEBUG 写入stopch")
			return
		}
	}
}
