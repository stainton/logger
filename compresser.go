package logger

import (
	"archive/zip"
	"container/list"
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

// Compresser continuously compresses files received through a channel into zip archives.
// It runs as a long-lived goroutine, processing files as they become available and
// handling cancellation via a context.
//
// Parameters:
//   - ctx: A context.Context for cancellation control.
//   - dir: The directory path where the zip files will be created and where source files are located.
//   - compressible: A receive-only channel of strings, each string representing a filename to be compressed.
//
// The function does not return any value. It runs indefinitely until the context is cancelled.
func Compresser(ctx context.Context, dir string, compressible <-chan string, interval int, stop chan<- error) {
	fmt.Printf("[%s]compresser started... ^^\n", time.Now().String())
	nowTime := time.Now()
	fileToBeCompressed := list.New()
	for {
		select {
		case file, ok := <-compressible:
			fmt.Println("DEBUG 有文件要压缩")
			if (fileToBeCompressed.Len() > 10) || (time.Now().Hour() != nowTime.Hour()) || !ok {
				fmt.Println("DEBUG 准备压缩1")
				nowTime = time.Now()
				compress(dir, &nowTime, fileToBeCompressed)
				fmt.Println("DEBUG 压缩完毕1")
			}
			if !ok {
				fmt.Printf("[%s]compresser stopped... ^^\n", time.Now().String())
				stop <- nil
				return
			}
			// TODO: 有可能超过10个，移到队末？
			fileToBeCompressed.PushBack(file)
		case <-ctx.Done():
			fmt.Println("DEBUG 准备要退出")
			for file := range compressible {
				fileToBeCompressed.PushBack(file)
			}
			fmt.Println("DEBUG 已经清空带压缩队列")
			nowTime = time.Now()
			fmt.Println("DEBUG 准备压缩")
			compress(dir, &nowTime, fileToBeCompressed)
			fmt.Println("DEBUG 压缩完成")
			fmt.Printf("[%s]compresser stopped... ^^\n", time.Now().String())
			stop <- nil
			return
		default:
			time.Sleep(time.Second * time.Duration(interval))
		}
	}
}

func compress(dir string, nowTime *time.Time, fileList *list.List) {
	if fileList.Len() == 0 {
		return
	}
	fileToDelete := list.New()
	prefix := strings.Replace(nowTime.Format("20060102150405.999"), ".", "", -1)
	zipFileName := path.Join(dir, fmt.Sprintf("%s.zip", prefix))
	fd, err := os.OpenFile(zipFileName, os.O_CREATE|os.O_WRONLY|os.O_EXCL, os.ModePerm)
	if err != nil {
		// 创建压缩文件失败，不删除队列中的文件
		fmt.Printf("[%s][ERROR] failed to create zip file: %v\n", nowTime.String(), err)
		return
	}
	defer fd.Close()
	w := zip.NewWriter(fd)
	defer w.Close()
	for fileToDelete.Len() != fileList.Len() {
		fileToCompress := fileList.Front()
		fileName, ok := fileToCompress.Value.(string)
		if !ok {
			// 解析文件名称错误，不删除队列文件
			// TODO: 会有压缩文件残留(是否删除)
			fmt.Printf("[%s][ERROR] incorrect filename in list: %v\n", nowTime.String(), fileName)
			return
		}
		f, err := w.Create(fileName)
		if err != nil {
			// 创建文件出错，不删除队列文件
			// TODO: 会有压缩文件残留(是否删除)
			fmt.Printf("[%s][ERROR] failed to compress file %s: %v\n", nowTime.String(), f, err)
			return
		}
		content, err := os.ReadFile(fileName)
		if err != nil {
			// 读源文件出错，不删除队列文件
			// TODO: 会有压缩文件残留(是否删除)
			fmt.Printf("[%s][ERROR] failed to compress file %s: %v\n", nowTime.String(), fileName, err)
			return
		}
		total := 0
		for n, err := f.Write(content[total:]); total < len(content); total += n {
			if err != nil {
				// 写压缩文件出错
				// TODO: 会有压缩文件残留(是否删除)
				fmt.Printf("[%s][ERROR] failed to compress file %s: %v\n", nowTime.String(), fileName, err)
				return
			}
		}
		fileToDelete.PushBack(fileName)
		fileList.MoveToBack(fileToCompress)
	}
	fileList.Init()
	for fileToDelete != nil && fileToDelete.Len() > 0 {
		f := fileToDelete.Remove(fileToDelete.Front())
		name := f.(string)
		err := os.Remove(name)
		if err != nil {
			fmt.Printf("[%s][ERROR] failed to remove log file %s\n", nowTime.String(), name)
			return
		}
	}
}
