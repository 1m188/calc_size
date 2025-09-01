package main

import (
	"fmt"
	"os"
	"path/filepath"
)

var help_info string = `输入文件名或文件夹名称来计算文件或文件夹大小: 
	calc_size [Your file name or folder name]

或者是连续计算多个文件或者文件夹的大小: 
	calc_size [file1] [file2] [file3] ...

输出-1则表示该文件不存在或文件夹存在路径错误或者其他错误`

func main() {

	args := os.Args

	// 输出帮助信息
	if len(args) == 2 && (args[1] == "--help" || args[1] == "--h" || args[1] == "-help" || args[1] == "-h") {
		fmt.Println(help_info)
		return
	}

	file_paths := args[1:]                       // 文件路径
	file_sizes := make([]int64, len(file_paths)) // 计算出来的文件大小

	for i, file_path := range file_paths {
		file_info, err := os.Stat(file_path)

		if err != nil {
			file_sizes[i] = -1
			continue
		}

		if !file_info.IsDir() {
			file_sizes[i] = file_info.Size()
			continue
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
	}

	for i, file_path := range file_paths {
		size := float64(file_sizes[i])
		unit := "B"

		if size > 1024 {
			size /= 1024
			unit = "KB"

			if size > 1024 {
				size /= 1024
				unit = "MB"

				if size > 1024 {
					size /= 1024
					unit = "GB"
				}
			}
		}

		fmt.Printf("%s: %.2f %s\n", file_path, size, unit)
	}

}
