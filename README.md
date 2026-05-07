# logger

基于 zap 的轻量日志接口封装，提供统一的日志方法与可选配置。

## 安装

```bash
go get github.com/stainton/logger
```

## 包结构

- `github.com/stainton/logger`: 日志接口与通用配置
- `github.com/stainton/logger/impls/zapimpl`: zap 实现

## 快速开始

```go
package main

import (
	"log"

	"github.com/stainton/logger"
	"github.com/stainton/logger/impls/zapimpl"
)

func main() {
	opt := logger.NewOptions(
		logger.WithEncoding(logger.EncodingJSON),
		logger.WithMinLevel(logger.Info),
		logger.WithOutputPaths([]string{"stdout"}),
	)

	l, err := zapimpl.NewProductionLogger(opt)
	if err != nil {
		log.Fatalf("create logger failed: %v", err)
	}
	defer l.Sync()

	l.Infof("service started, name=%s version=%s", "demo", "v1.0.0")
	l.Warn("this is a warning")
	l.Errorf("request failed, code=%d", 500)
}
```

## 创建日志器

- `zapimpl.NewDevelopLogger(opts)`: 开发模式，输出更易读
- `zapimpl.NewProductionLogger(opts)`: 生产模式，适合结构化日志

两者都返回 `logger.Logger` 接口。

## 可配置项

通过 `logger.NewOptions(...)` 传入可选项：

- `logger.WithOutputPaths([]string{...})`
  - 默认值：`[]string{"stdout"}`
  - 可写入标准输出或文件路径（例如 `./_log/app.json`）
- `logger.WithEncoding(logger.EncodingJSON | logger.EncodingConsole)`
  - 默认值：`logger.EncodingConsole`
- `logger.WithMinLevel(logger.Debug ~ logger.Fatal)`
  - 默认值：`logger.Info`

可用日志级别：

- `logger.Debug`
- `logger.Info`
- `logger.Warn`
- `logger.Error`
- `logger.DPanic`
- `logger.Fatal`

## 写入文件示例

```go
opt := logger.NewOptions(
	logger.WithEncoding(logger.EncodingJSON),
	logger.WithMinLevel(logger.Debug),
	logger.WithOutputPaths([]string{"./_log/t.json"}),
)

l, err := zapimpl.NewDevelopLogger(opt)
if err != nil {
	panic(err)
}
defer l.Sync()

l.Debugf("user=%s login", "alice")
```

## 接口能力

`logger.Logger` 提供以下方法族：

- `Debug / Debugf / Debugln`
- `Info / Infof / Infoln`
- `Warn / Warnf / Warnln`
- `Error / Errorf / Errorln`
- `DPanic / DPanicf / DPanicln`
- `Fatal / Fatalf / Fatalln`
- `Sync() error`

## 注意事项

- 建议在程序退出前调用 `Sync()`，避免缓冲日志丢失。
- `WithOutputPaths` 中的文件路径需确保进程有写权限。