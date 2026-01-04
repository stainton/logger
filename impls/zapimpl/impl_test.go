package zapimpl

import (
	"sync"
	"testing"

	"github.com/stainton/logger"
)

func TestLogger(t *testing.T) {
	t.Run("single-log", func(t *testing.T) {
		opt := logger.NewOptions(
			logger.WithEncoding(logger.EncodingJSON),
			logger.WithMinLevel(logger.Debug),
			logger.WithOutputPaths([]string{"/home/hanjunzheng/projects/github.com/InHuanLe/logger/_log/t.json"}),
		)
		l, err := NewDevelopLogger(opt)
		if err != nil {
			t.Errorf("new logger failed: %+v", err)
		}
		l.Debugf("print test log")
	})
	t.Run("multiple-goroutine", func(t *testing.T) {
		opt := logger.NewOptions(
			logger.WithEncoding(logger.EncodingJSON),
			logger.WithMinLevel(logger.Debug),
			logger.WithOutputPaths([]string{"/home/hanjunzheng/projects/github.com/InHuanLe/logger/_log/t-multi.json"}),
		)
		l, err := NewDevelopLogger(opt)
		if err != nil {
			t.Errorf("new logger failed: %+v", err)
		}
		wg := sync.WaitGroup{}
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				for j := 0; j < 100; j++ {
					l.Debugf("goroutine(%+v), log-index(%+v)", i, j)
				}
			}(i)
		}
		wg.Wait()
	})
}
