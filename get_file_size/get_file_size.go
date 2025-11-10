package get_files_size

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"example.com/calc_size/data"
	rtp "example.com/calc_size/get_file_size/internal/rtprint"
)

var IsCnt bool   // 是否需要统计已完成多少个文件的大小扫描
var MaxDepth int // 最大目录递归深度

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
			rtp.CountAddPrint(1)
		}
		return size
	}

	// 计算起始目录的深度
	rootDepth := strings.Count(filepath.Clean(file_path), string(os.PathSeparator))

	total := big.NewInt(0)
	filepath.Walk(file_path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			total.SetInt64(-1)
			return err
		}

		// 计算当前路径的深度
		currentDepth := strings.Count(filepath.Clean(path), string(os.PathSeparator))
		// 计算相对深度
		relativeDepth := currentDepth - rootDepth
		// 如果超过最大深度，则跳过该目录
		if info.IsDir() && relativeDepth > int(MaxDepth) {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			total.Add(total, big.NewInt(info.Size()))
			if IsCnt {
				rtp.CountAddPrint(1)
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

	if IsCnt {
		rtp.Init()
	}

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
