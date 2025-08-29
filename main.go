package main

import (
	"fmt"
	"os"
)

func main() {

	file_name := "go.mod"

	file_info, err := os.Stat(file_name)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("file size (%s): %d bytes\n", file_name, file_info.Size())

}
