package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// 获取文件/文件夹大小，单位为字节
func getFileSize(file_paths []string) []int64 {
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

func main() {

	args := os.Args

	// 输出帮助信息
	if len(args) == 2 && (args[1] == "--help" || args[1] == "--h" || args[1] == "-help" || args[1] == "-h") {
		fmt.Println(help_info)
		return
	}

	file_paths := args[1:] // 文件路径

	file_sizes := getFileSize(file_paths) // 文件大小

	print_size(file_paths, file_sizes) // 打印文件大小

}
