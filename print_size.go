package main

import (
	"fmt"
	"strings"
)

// 输出文件大小
func getSizeInFmt(file_sizes []FileSize) string {

	// 补全所有文件路径长度到最长文件路径长度
	file_paths := make([]string, len(file_sizes))
	for i, file_size := range file_sizes {
		file_paths[i] = file_size.Path
	}
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

	// 转换文件大小单位
	file_sizes_str := make([]string, len(file_sizes)) // 转换后的文件大小
	units := make([]string, len(file_sizes))          // 转换后的文件大小单位
	for i, file_size := range file_sizes {
		file_sizes_str[i], units[i] = file_size.GetSizeWithUnit(2, 512)
		if units[i] == "B" {
			units[i] = " B"
		}
		if units[i] == "" {
		    units[i] = "  "
		}
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
	res := make([]string, len(file_sizes))
	for i := range file_paths {
		res[i] = fmt.Sprintf("%s : %s %s\n", file_paths[i], file_sizes_str[i], units[i])
	}
	return strings.Join(res, "")
}
