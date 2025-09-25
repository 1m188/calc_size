package rtprint

import (
	"fmt"
	"sync"
)

var count uint64     // 已完成多少个文件的大小扫描
var mutex sync.Mutex // 互斥锁

/* 初始化 */
func Init() {
	count = 0
}

/*
输出已完成多少个文件的扫描
  - x: 增加的文件扫描数量
*/
func CountAddPrint(x uint64) {
	mutex.Lock()
	count += x
	fmt.Printf("\r已完成 %d 个文件的扫描", count) // 打印进度
	mutex.Unlock()
}
