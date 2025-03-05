package logger

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
)

var globalLogger *BaseLogger
var once = &sync.Once{}
var stopCh chan struct{}

type BaseLogger struct {
	logChan chan<- string
}

type Options struct {
	// Output directory for log files.
	Output string
	// Maximum number of log messages to keep in memory before writing to a file.
	MaxMessage int
	// Maximum size of a log file in bytes before it is compressed and a new file is opened.
	Threshold int

	// Interval in seconds for compressing log files.
	CompressInterval int
}

func initialOptions() *Options {
	dir := os.Getenv("LOG_OUTPUT")
	if dir == "" {
		dir = "./logs"
	}
	max_message, err := strconv.Atoi(os.Getenv("LOG_MAX_MESS"))
	if err != nil || max_message <= 0 {
		fmt.Println("Invalid LOG_MAX_MESS environment variable, using default value 1000. Please set it correctly. ^^")
		max_message = 1000
	}
	threshold, err := strconv.Atoi(os.Getenv("LOG_THRESHOLD_BYTE"))
	if err != nil || threshold <= 0 {
		fmt.Println("Invalid LOG_THRESHOLD_BYTE environment variable, using default value 10MiB. Please set it correctly. ^^")
		threshold = Threshold
	}
	interval, err := strconv.Atoi(os.Getenv("LOG_COMPRESS_INTERVAL"))
	if err != nil || interval <= 0 {
		fmt.Println("Invalid LOG_COMPRESS_INTERVAL environment variable, using default value 3600 seconds. Please set it correctly. ^^")
		interval = 60 * 60
	}
	return &Options{
		Output:           dir,
		MaxMessage:       max_message,
		Threshold:        threshold,
		CompressInterval: interval,
	}
}

func NewLogger(ctx context.Context, opts *Options) (Logger, <-chan struct{}) {
	once.Do(func() {
		if opts == nil {
			opts = initialOptions()
		}
		dir := opts.Output
		messages := make(chan string, opts.MaxMessage)
		compressible := make(chan string, 10)
		errOccur := make(chan error)
		compresserExited := make(chan error)
		stopCh = make(chan struct{})

		// 创建一个goroutine用于压缩文件
		go Compresser(ctx, dir, compressible, opts.CompressInterval, compresserExited)
		// 创建一个goroutine用于写文件
		go Writer(ctx, dir, messages, compressible, errOccur, opts.Threshold)
		go func() {
			// 等待writer和compresser退出
			wrErr := <-errOccur
			close(compressible)
			fmt.Println("writer exit with error:", wrErr)
			cpErr := <-compresserExited
			fmt.Println("compresser exit with error:", cpErr)
			close(stopCh)
		}()

		globalLogger = &BaseLogger{
			logChan: messages,
		}
	})
	return globalLogger, stopCh
}
