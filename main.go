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
	flag.Parse()

	file_paths := os.Args[1:] // 文件路径

	file_sizes := getFilesSize(file_paths) // 文件大小

	printSize(file_paths, file_sizes) // 打印文件大小

}
