package main

import (
	"example.com/calc_size/utils"
	"fmt"
	"os"
)

func main() {

	args := os.Args

	// 输出帮助信息
	if len(args) == 2 && (args[1] == "--help" || args[1] == "--h" || args[1] == "-help" || args[1] == "-h") {
		fmt.Println(utils.HelpInfo)
		return
	}

	file_paths := args[1:] // 文件路径

	file_sizes := utils.GetFileSize(file_paths) // 文件大小

	utils.PrintSize(file_paths, file_sizes) // 打印文件大小

}
