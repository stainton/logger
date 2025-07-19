package logger

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/spf13/cobra"
)

const (
	KiB         = 1 << 10
	MiB         = KiB * 1024
	DEFAULT_THD = 10 * MiB
)

var globalLogger *BaseLogger
var once = &sync.Once{}
var stopCh chan struct{}

type BaseLogger struct {
	serviceName string
	logChan     chan<- string
}

type Options struct {
	// Service name for log messages.
	// Default to hostname or hostname or "default" if failed to get it.
	ServiceName string
	// Output directory for log files.
	// Default to "./logs".
	Output string
	// Maximum number of log messages to keep in memory before writing to a file.
	// Default to 1000 messages.
	MaxMessage int
	// Maximum size of a log file in bytes before it is compressed and a new file is opened.
	// Default to 10MiB.
	Threshold int

	// Interval in seconds for compressing log files.
	// Default to 3600 seconds (1 hour).
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
		threshold = DEFAULT_THD
	}
	interval, err := strconv.Atoi(os.Getenv("LOG_COMPRESS_INTERVAL"))
	if err != nil || interval <= 0 {
		fmt.Println("Invalid LOG_COMPRESS_INTERVAL environment variable, using default value 3600 seconds. Please set it correctly. ^^")
		interval = 60 * 60
	}
	serviceName, err := os.Hostname()
	if err != nil || serviceName == "" {
		fmt.Printf("Failed to get hostname, using default value 'default'. Please set the hostname correctly. ^^")
		serviceName = "default"
	}
	return &Options{
		ServiceName:      serviceName,
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
		compresserExited := make(chan struct{})
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
			<-compresserExited
			fmt.Println("compresser exit")
			close(stopCh)
		}()
		if opts.ServiceName == "" {
			serviceName, err := os.Hostname()
			if err != nil || serviceName == "" {
				fmt.Printf("Failed to get hostname, using default value 'default'. Please set the hostname correctly. ^^")
				serviceName = "default"
			}
			opts.ServiceName = serviceName
		}
		globalLogger = &BaseLogger{
			serviceName: opts.ServiceName,
			logChan:     messages,
		}
	})
	return globalLogger, stopCh
}

func NewLogOptions(serviceName string) *Options {
	return &Options{
		ServiceName: serviceName,
	}
}

func (o *Options) FlagSet(cmd *cobra.Command) {
	cmd.Flags().IntVar(&o.CompressInterval, "log-compress-interval", 3600, "日志文件压缩间隔")
	cmd.Flags().IntVar(&o.Threshold, "log-compress-size", DEFAULT_THD, "日志文件压缩阈值")
	cmd.Flags().IntVar(&o.MaxMessage, "log-max_message", 1000, "允许存在内存中的最多日志条数")
	cmd.Flags().StringVar(&o.Output, "log-output", "./log", "日志输出目录")
}
