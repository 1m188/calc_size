package get_files_size

import (
	"example.com/calc_size/data"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"sync"
)

var count uint64     // 已完成多少个文件的大小扫描
var mutex sync.Mutex // 互斥锁

var IsCnt bool // 是否需要统计已完成多少个文件的大小扫描

/* 输出已完成多少个文件的扫描 */
func print_count() {
	mutex.Lock()
	count++
	fmt.Printf("\r已完成 %d 个文件的扫描", count) // 打印进度
	mutex.Unlock()
}

/*
获取文件/文件夹大小，返回 Size=-1 表示存在错误
  - file_path: 文件/文件夹路径
  - return: FileSize
*/
func getFileSize(file_path string) data.FileSize {
	size := data.FileSize{Path: file_path, Size: *big.NewInt(-1)}
	file_info, err := os.Stat(file_path)

	if err != nil {
		return size
	}

	if !file_info.IsDir() {
		size.Size = *big.NewInt(file_info.Size())
		if IsCnt {
			print_count()
		}
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
			if IsCnt {
				print_count()
			}
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
func GetFilesSize(file_paths []string) []data.FileSize {
	file_sizes := make([]data.FileSize, len(file_paths)) // 计算出来的文件大小
	var wg sync.WaitGroup                                // 等待所有goroutine完成

	for i, file_path := range file_paths {
		wg.Go(func() {
			file_sizes[i] = getFileSize(file_path)
		})
	}

	wg.Wait()

	if IsCnt {
		fmt.Println()
	}

	return file_sizes
}
