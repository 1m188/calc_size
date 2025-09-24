package main

import (
	"math/big"
	"os"
	"path/filepath"
	"sync"
)

/*
获取文件/文件夹大小，返回 Size=-1 表示存在错误
  - file_path: 文件/文件夹路径
  - return: FileSize
*/
func getFileSize(file_path string) FileSize {
	size := FileSize{Path: file_path, Size: *big.NewInt(-1)}
	file_info, err := os.Stat(file_path)

	if err != nil {
		return size
	}

	if !file_info.IsDir() {
		size.Size = *big.NewInt(file_info.Size())
		return size
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
	size.Size = *total

	return size
}

/*
获取文件/文件夹大小，单位为字节
  - file_paths: 文件/文件夹路径列表
  - return: FileSize 列表
*/
func getFilesSize(file_paths []string) []FileSize {
	file_sizes := make([]FileSize, len(file_paths)) // 计算出来的文件大小
	var wg sync.WaitGroup                           // 等待所有goroutine完成

	for i, file_path := range file_paths {
		wg.Go(func() {
			file_sizes[i] = getFileSize(file_path)
		})
	}

	wg.Wait()

	return file_sizes
}
