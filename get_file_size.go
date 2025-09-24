package main

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"sync"
)

// 获取文件/文件夹大小，返回 "-1" 表示存在错误
func getFileSize(file_path string) (size string) {
	file_info, err := os.Stat(file_path)

	if err != nil {
		size = "-1"
		return
	}

	if !file_info.IsDir() {
		size = fmt.Sprintf("%d", file_info.Size())
		return
	}

	total := big.NewInt(0)
	filepath.Walk(file_path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			total.SetInt64(-1)
			return err
		}

		if !info.IsDir() {
			total.Add(total, big.NewInt(info.Size()))
		}

		return err
	})
	size = total.String()

	return
}

// 获取文件/文件夹大小，单位为字节
func getFilesSize(file_paths []string) []string {
	file_sizes := make([]string, len(file_paths)) // 计算出来的文件大小
	var wg sync.WaitGroup                         // 等待所有goroutine完成

	for i, file_path := range file_paths {
		wg.Go(func() {
			file_sizes[i] = getFileSize(file_path)
		})
	}

	wg.Wait()

	return file_sizes
}
