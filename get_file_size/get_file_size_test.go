package get_files_size

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"example.com/calc_size/data"
	rtp "example.com/calc_size/get_file_size/internal/rtprint"
)

// 创建固定大小文件
func makeFile(file_path string, size int64) error {
	path := filepath.Dir(file_path)
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}
	file, err := os.Create(file_path)
	if err != nil {
		return err
	}
	defer file.Close()
	for b := int64(0); b < size; {
		var block int64 = 1024 * 1024 * 1024
		if size-b < block {
			block = size - b
		}
		_, err := file.Write(make([]byte, block))
		if err != nil {
			return err
		}
		b += block
	}
	return nil
}

/* 测试 getFileSize 基本函数功能 */
func TestGetFileSize(t *testing.T) {
	sizes := []int64{0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024*1024 + 1234}
	var sum big.Int
	for i := range sizes {
		sum.Add(&sum, big.NewInt(sizes[i]))
	}

	// 测试文件大小
	for i := range len(sizes) {
		file_name := fmt.Sprintf(".test%d", i)
		file, err := os.Create(file_name)
		if err != nil {
			t.Fatal(err)
		}
		for b := int64(0); b < sizes[i]; {
			var block int64 = 1024 * 1024 * 1024
			if sizes[i]-b < block {
				block = sizes[i] - b
			}
			_, err := file.Write(make([]byte, block))
			if err != nil {
				t.Fatal(err)
			}
			b += block
		}
		file.Close()

		fileSize := getFileSize(file_name)
		if fileSize.Size.Cmp(big.NewInt(sizes[i])) != 0 {
			t.Errorf("Expected %d, got %s", sizes[i], fileSize.Size.String())
		}

		os.Remove(file_name)
	}

	// 测试文件夹大小
	err := os.Mkdir(".test", 0755)
	if err != nil {
		t.Fatal(err)
	}

	for i := range len(sizes) {
		file_name := fmt.Sprintf("%s/.test%d", ".test", i)
		file, err := os.Create(file_name)
		if err != nil {
			t.Fatal(err)
		}
		for b := int64(0); b < sizes[i]; {
			var block int64 = 1024 * 1024 * 1024
			if sizes[i]-b < block {
				block = sizes[i] - b
			}
			_, err := file.Write(make([]byte, block))
			if err != nil {
				t.Fatal(err)
			}
			b += block
		}
		file.Close()
	}
	fileSize := getFileSize(".test")
	if fileSize.Size.Cmp(&sum) != 0 {
		t.Errorf("Expected %s, got %s", sum.String(), fileSize.Size.String())
	}
	os.RemoveAll(".test")

	// 测试不存在的文件或文件夹
	fileSize = getFileSize("./path/not/exist/file")
	if fileSize.Size.Cmp(big.NewInt(-1)) != 0 {
		t.Errorf("Expected -1, got %s", fileSize.Size.String())
	}
}

// 测试对比单线程vs并发，在不同文件规模下的性能表现
func TestRoutinePerformance(t *testing.T) {

	// 单线程——获取文件/文件夹大小（尽可能和并发函数保持一致的流程）
	GetFilesSizeSingle := func(file_paths []string) []data.FileSize {
		file_sizes := make([]data.FileSize, len(file_paths)) // 计算出来的文件大小

		if IsCnt {
			rtp.Init()
		}

		for i, file_path := range file_paths {
			file_sizes[i] = getFileSize(file_path)
		}

		if IsCnt {
			fmt.Println()
		}

		return file_sizes
	}

	test := func(sizes []int64) {
		var sum big.Int
		for i := range sizes {
			sum.Add(&sum, big.NewInt(sizes[i]))
		}

		file_names := make([]string, len(sizes))
		for i := range sizes {
			file_names[i] = fmt.Sprintf(".test%d", i)
		}

		// 创建文件
		for i := range len(sizes) {
			err := makeFile(file_names[i], sizes[i])
			if err != nil {
				t.Fatal(err)
			}
			err = makeFile(".test/"+file_names[i], sizes[i])
			if err != nil {
				t.Fatal(err)
			}
		}

		file_names = append(file_names, ".test/")                 // 文件夹
		file_names = append(file_names, "./path/not/exist/.file") // 不存在的文件

		// 多线程
		start := time.Now()
		filesSize := GetFilesSize(file_names)
		fmt.Printf("time cost with routine: %d μs\n", time.Since(start).Microseconds())
		for i := range sizes {
			fileSize := filesSize[i]
			if fileSize.Size.Cmp(big.NewInt(sizes[i])) != 0 {
				t.Errorf("Expected %d, got %s", sizes[i], fileSize.Size.String())
			}
		}
		if fileSize := filesSize[len(filesSize)-2]; fileSize.Size.Cmp(&sum) != 0 {
			t.Errorf("Expected %s, got %s", sum.String(), fileSize.Size.String())
		}
		if fileSize := filesSize[len(filesSize)-1]; fileSize.Size.Cmp(big.NewInt(-1)) != 0 {
			t.Errorf("Expected %d, got %s", -1, fileSize.Size.String())
		}

		// 单线程
		start = time.Now()
		filesSize = GetFilesSizeSingle(file_names)
		fmt.Printf("time cost without routine: %d μs\n", time.Since(start).Microseconds())
		for i := range sizes {
			fileSize := filesSize[i]
			if fileSize.Size.Cmp(big.NewInt(sizes[i])) != 0 {
				t.Errorf("Expected %d, got %s", sizes[i], fileSize.Size.String())
			}
		}
		if fileSize := filesSize[len(filesSize)-2]; fileSize.Size.Cmp(&sum) != 0 {
			t.Errorf("Expected %s, got %s", sum.String(), fileSize.Size.String())
		}
		if fileSize := filesSize[len(filesSize)-1]; fileSize.Size.Cmp(big.NewInt(-1)) != 0 {
			t.Errorf("Expected %d, got %s", -1, fileSize.Size.String())
		}

		// 删除文件
		for i := range len(sizes) {
			file_name := fmt.Sprintf(".test%d", i)
			err := os.Remove(file_name)
			if err != nil {
				t.Error(err)
			}
		}
		err := os.RemoveAll(".test")
		if err != nil {
			t.Error(err)
		}
	}

	fmt.Println("规模1")
	test([]int64{
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024*1024 + 1234,
	})

	fmt.Println("规模2")
	test([]int64{
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
		0, 10 + 1, 1024 + 12, 1024*1024*2 + 123, 1024*1024 + 1234,
	})

}
