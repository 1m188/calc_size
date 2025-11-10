package main

import (
	"fmt"
	"math"
	"os"

	gfs "example.com/calc_size/get_file_size"
	"example.com/calc_size/print"
	flag "github.com/spf13/pflag"
)

// 帮助信息
const HELPINFO string = `输入文件名或文件夹名称来计算文件或文件夹大小: 
	calc_size -p [Your file name or folder name]

或者是连续计算多个文件或者文件夹的大小: 
	calc_size -p [file1],[file2],[file3] ...

输出-1则表示该文件不存在或文件夹存在路径错误或者其他错误`

func main() {

	// 输出帮助信息
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n", HELPINFO)
		flag.PrintDefaults()
	}

	var json bool
	var file_paths []string
	var csv bool
	var cnt bool
	var depth int // 递归最大深度

	flag.BoolVarP(&json, "json", "j", false, "以json格式输出")
	flag.BoolVarP(&csv, "csv", "c", false, "以csv格式输出")
	flag.StringSliceVarP(&file_paths, "paths", "p", []string{}, "文件/目录路径(英文逗号分隔)")
	flag.BoolVarP(&cnt, "count", "n", false, "是否实时输出统计文件数量")
	flag.IntVarP(&depth, "depth", "d", math.MaxInt, "递归最大深度")
	flag.Parse()

	gfs.IsCnt = cnt
	gfs.MaxDepth = depth

	var res string
	if json {
		file_sizes := gfs.GetFilesSize(file_paths) // 文件大小
		res = print.GetSizeInJSON(file_sizes)      // 打印json格式文件大小
	} else if csv {
		file_sizes := gfs.GetFilesSize(file_paths) // 文件大小
		res = print.GetSizeInCSV(file_sizes)       // 打印csv格式文件大小
	} else {
		file_sizes := gfs.GetFilesSize(file_paths) // 文件大小
		res = print.GetSizeInFmt(file_sizes)       // 打印文件大小
	}
	fmt.Print(res)

}
