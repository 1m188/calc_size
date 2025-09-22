package main

import (
	"example.com/calc_size/utils"
	"fmt"
	flag "github.com/spf13/pflag"
	"os"
)

func main() {

	// 输出帮助信息
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n", utils.HelpInfo)
		flag.PrintDefaults()
	}
	flag.Parse()

	file_paths := os.Args[1:] // 文件路径

	file_sizes := utils.GetFileSize(file_paths) // 文件大小

	utils.PrintSize(file_paths, file_sizes) // 打印文件大小

}
