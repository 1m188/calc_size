package main

import (
	"fmt"
	"strings"
)

// 输出文件大小
func print_size(file_paths []string, file_sizes []int64) {

	// 补全所有文件路径长度到最长文件路径长度
	leftLen := 0
	for _, file_path := range file_paths {
		if len(file_path) > leftLen {
			leftLen = len(file_path)
		}
	}
	for i, file_path := range file_paths {
		if len(file_path) < leftLen {
			file_path += strings.Repeat(" ", leftLen-len(file_path))
			file_paths[i] = file_path
		}
	}

	file_sizes_str := make([]string, len(file_sizes))

	// 转换文件大小单位
	for i := range file_paths {
		size := float64(file_sizes[i])
		unit := " B"

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

		if size < 0 {
			unit = "  "
		}

		file_sizes_str[i] = fmt.Sprintf("%.2f %s", size, unit)
	}

	// 补全所有文件大小长度到最长文件大小长度
	rightLen := 0
	for _, file_size_str := range file_sizes_str {
		if len(file_size_str) > rightLen {
			rightLen = len(file_size_str)
		}
	}
	for i, file_size_str := range file_sizes_str {
		if len(file_size_str) < rightLen {
			file_size_str = strings.Repeat(" ", rightLen-len(file_size_str)) + file_size_str
			file_sizes_str[i] = file_size_str
		}
	}

	// 输出文件路径和文件大小
	for i := range file_paths {
		fmt.Printf("%s  :  %s\n", file_paths[i], file_sizes_str[i])
	}
}
