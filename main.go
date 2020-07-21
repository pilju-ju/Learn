// main
package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	fileName := os.Args[1]
	outputName := os.Args[2]
	arrayName := os.Args[3]

	file, err := os.Open(fileName)
	if err != nil {
		log.Println("file open err", err)
		return
	}
	defer file.Close()

	output, err := os.OpenFile(outputName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("output file err", err)
		return
	}
	defer output.Close()

	rb := bufio.NewReader(file)

	count := 0
	comma := 0
	output.WriteString("#include <stdint.h>\n\n")
	output.WriteString("uint8_t " + arrayName + "[]={\n ")

	for {
		data, err := rb.ReadByte()
		if err != nil {
			break
		}
		if comma != 0 {
			output.WriteString(", ")
		}

		if count == 16 {
			count = 0
			output.WriteString("\n")
		}

		count++
		if data < 16 {
			output.WriteString("0x0" + strings.ToUpper(strconv.FormatUint(uint64(data), 16)))
		} else {
			output.WriteString("0x" + strings.ToUpper(strconv.FormatUint(uint64(data), 16)))
		}
		comma = 1

	}
	output.WriteString("\n}; ")
}
