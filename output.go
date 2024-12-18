package logger

import (
	"io"
	"net"
	"net/netip"
	"os"
)

type writer struct {
	conn net.Conn
	f    *os.File
}

// 负责往远端和本地写日志
func (w *writer) Write(p []byte) (n int, err error) {
	if w.conn != nil {
		// 写远端日志
		n, err = w.conn.Write(p)
	}
	// 写本地日志
	n, err = w.f.Write(p)
	return
}

// 关闭链接和文件
func (w *writer) Close() {
	if w.conn != nil {
		w.conn.Close()
	}
	w.f.Close()
}

// Output 将远端的输出(connectionString)和本地文件输出(filepath)绑定在一起
func Output(connectionString string, filePath string) io.Writer {
	peer := ""
	var err error
	var conn net.Conn = nil
	if connectionString != "" {
		ipport := netip.MustParseAddrPort(connectionString)
		peer = ipport.String()
		conn, err = net.Dial("tcp", peer)
		if err != nil {
			panic(err)
		}
	}

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		panic(err)
	}
	return &writer{
		conn: conn,
		f:    f,
	}
}
