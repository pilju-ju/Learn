// main.go
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

var bufSize = 512

func main() {
	file1 := os.Args[1]
	file2 := os.Args[2]

	f1, idx1 := readFile(file1)
	f2, _ := readFile(file2)
	okCnt := 0
	match := 0
	for i := 0; i < idx1+1; i++ {
		match = 0
		for j := 0; j < len(f1[idx1]); j++ {
			if f1[idx1][j] != f2[idx1][j] {
				fmt.Println("not match")
				match = 1
			}
		}
		if match == 0 {
			okCnt++
		}
	}
	fmt.Println("okcnt : ", okCnt)
}

func readFile(inputFile string) ([][]byte, int) {
	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	testBufSize := bufSize
	fileStat, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return nil, -1
	}
	fb := bufio.NewReader(file)
	var fileData [][]byte
	fileSize := fileStat.Size()
	idx := int(fileSize) / testBufSize
	remain := int(fileSize) % testBufSize
	fmt.Println("filesize : ", fileSize, "idx :", idx, "remain", remain)
	if remain != 0 {
		idx = idx + 1
	}
	fileData = make([][]byte, idx)

	dataIdx := 0
	// filelen := 0

	for i := 0; i < idx; i++ {
		if remain != 0 {
			if i == idx-1 {
				fileData[i] = make([]byte, remain)
				// filelen += len(decryptDataBuf[i])
				io.ReadFull(fb, fileData[dataIdx])
				// fmt.Println("F Len ", filelen, "dlen :", len(fileData[dataIdx]))
				break
			}
		}
		fileData[i] = make([]byte, testBufSize)
		// filelen += len(decryptDataBuf[i])

		io.ReadFull(fb, fileData[dataIdx])
		// fmt.Println("F Len ", filelen, "dlen :", len(fileData[dataIdx]))
		dataIdx++
	}

	return fileData, dataIdx

}
