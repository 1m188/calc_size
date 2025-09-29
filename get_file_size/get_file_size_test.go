package get_files_size

import (
	"fmt"
	"math/big"
	"os"
	"testing"
)

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

		fileSizes := GetFilesSize([]string{file_name})
		if fileSizes[0].Size.Cmp(big.NewInt(sizes[i])) != 0 {
			t.Errorf("Expected %d, got %s", sizes[i], fileSizes[0].Size.String())
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
	fileSizes := GetFilesSize([]string{".test"})
	if fileSizes[0].Size.Cmp(&sum) != 0 {
		t.Errorf("Expected %s, got %s", sum.String(), fileSizes[0].Size.String())
	}
	os.RemoveAll(".test")
}
