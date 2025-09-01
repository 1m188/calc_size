package utils

import (
	"os"
	"path/filepath"
	"sync"
)

// 获取文件/文件夹大小，单位为字节
func GetFileSize(file_paths []string) []int64 {
	file_sizes := make([]int64, len(file_paths)) // 计算出来的文件大小
	var wg sync.WaitGroup                        // 等待所有goroutine完成

	for i, file_path := range file_paths {
		wg.Go(func() {
			file_info, err := os.Stat(file_path)

			if err != nil {
				file_sizes[i] = -1
				return
			}

			if !file_info.IsDir() {
				file_sizes[i] = file_info.Size()
				return
			}

			filepath.Walk(file_path, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					file_sizes[i] = -1
					return err
				}

				if !info.IsDir() {
					file_sizes[i] += info.Size()
				}

				return err
			})
		})
	}

	wg.Wait()

	return file_sizes
}
