package main

import (
	"fmt"
	"os"
)

func main() {

	args := os.Args

	// 输出帮助信息
	if len(args) == 2 && (args[1] == "--help" || args[1] == "--h" || args[1] == "-help" || args[1] == "-h") {

		fmt.Println(`输入文件名或文件夹名称来计算文件或文件夹大小: 
	calc_size [Your file name or folder name]

或者是连续计算多个文件或者文件夹的大小: 
	calc_size [file1] [file2] [file3] ...`)

		return
	}

	file_name := "go.mod"

	file_info, err := os.Stat(file_name)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("file size (%s): %d bytes\n", file_name, file_info.Size())

}
