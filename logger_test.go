package logger

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestInterface(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	var wg sync.WaitGroup
	l := NewLoggerWithWG(ctx, &wg, "logger_test", "", "my-log")
	l.Errorf("test message from 2024-12-18 01:28")
	time.Sleep(5 * time.Second)
	cancel()
	wg.Wait()
}
