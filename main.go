package main

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"os"
)

func main() {

	// 输出帮助信息
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n", helpInfo)
		flag.PrintDefaults()
	}

	var json bool
	var file_paths []string
	var csv bool

	flag.BoolVarP(&json, "json", "j", false, "以json格式输出")
	flag.BoolVarP(&csv, "csv", "c", false, "以csv格式输出")
	flag.StringSliceVarP(&file_paths, "paths", "p", []string{}, "文件/目录路径(英文逗号分隔)")
	flag.Parse()

	file_sizes := getFilesSize(file_paths) // 文件大小

	res := ""
	if json {
		res = getSizeInJSON(file_sizes) // 打印json格式文件大小
	} else if csv {
		res = getSizeInCSV(file_sizes) // 打印csv格式文件大小
	} else {
		res = getSizeInFmt(file_sizes) // 打印文件大小
	}
	fmt.Print(res)

}
