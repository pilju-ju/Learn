// main.go
package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"os"
)

var iv = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F}

func main() {
	keyValue := os.Args[1]
	key := make([]byte, 16)
	md5Key := md5.Sum([]byte(keyValue))

	for i, v := range md5Key {
		key[i] = v
	}

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		ucl := scanner.Text()
		fmt.Println("----------------------------------")
		fmt.Println("Plain Data : ")
		fmt.Printf("[")
		for i, val := range ucl {
			if i == len(ucl)-1 {
				fmt.Printf("0x%02X", val)
			} else {

				fmt.Printf("0x%02X, ", val)
			}
		}
		fmt.Printf("]\n")
		fmt.Println(ucl)

		encrypted := make([]byte, len(ucl))
		err := EncryptAESCFB(encrypted, []byte(ucl), key, iv)
		if err != nil {
			panic(err)
		}
		fmt.Println("\nEncrypted Data : ")
		fmt.Printf("[")
		for i, val := range encrypted {
			if i == len(encrypted)-1 {
				fmt.Printf("0x%02X", val)
			} else {

				fmt.Printf("0x%02X, ", val)
			}
		}
		fmt.Printf("]\n")

		// fmt.Printf("[ 0x% X, ]\n", encrypted)
		fmt.Println(string(encrypted))

		sEnc := base64.StdEncoding.EncodeToString(encrypted)
		fmt.Println("\nbase64 Data :")

		fmt.Printf("[")
		for i, val := range sEnc {
			if i == len(sEnc)-1 {
				fmt.Printf("0x%02X", val)
			} else {

				fmt.Printf("0x%02X, ", val)
			}
		}
		fmt.Printf("]\n")
		fmt.Println(sEnc)

	}

}

func EncryptAESCFB(dst, src, key, iv []byte) error {
	aesBlockEncrypter, err := aes.NewCipher([]byte(key))
	if err != nil {
		return err
	}
	aesEncrypter := cipher.NewCFBEncrypter(aesBlockEncrypter, iv)
	aesEncrypter.XORKeyStream(dst, src)
	return nil
}

func DecryptAESCFB(dst, src, key, iv []byte) error {
	aesBlockDecrypter, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil
	}
	aesDecrypter := cipher.NewCFBDecrypter(aesBlockDecrypter, iv)
	aesDecrypter.XORKeyStream(dst, src)
	return nil
}
