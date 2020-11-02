package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
)

var key = "1234567890abcdef"

func main() {

	data := "가나다라마바사아자차카타파하 ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	fmt.Println("data len", len(data))
	fmt.Println("origin", []byte(data), "\n")
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		fmt.Println(err)
		return
	}

	ciphertext := newCBCEncrypt(block, []byte(data))
	fmt.Println("cipher", ciphertext, "\n")

	plaintext := newCBCDecrypt(block, ciphertext)
	fmt.Println("plain", string(plaintext), "\n")
	fmt.Println("convert", plaintext)

}

func newCBCEncrypt(b cipher.Block, plaintext []byte) []byte {
	plaintext = pkcs7Pad(plaintext, b.BlockSize())

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		fmt.Println(err)
		return nil
	}

	cbc := cipher.NewCBCEncrypter(b, iv)
	cbc.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext
}

func newCBCDecrypt(b cipher.Block, ciphertext []byte) (plaintext []byte) {
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	cbc := cipher.NewCBCDecrypter(b, iv)
	cbc.CryptBlocks(ciphertext, ciphertext)
	plaintext = pkcs7Unpad(ciphertext, aes.BlockSize)

	return plaintext
}

func pkcs7Pad(plaintext []byte, blocksize int) []byte {
	padBuf := &bytes.Buffer{}
	n := blocksize - (len(plaintext) % blocksize)

	binary.Write(padBuf, binary.LittleEndian, plaintext)
	binary.Write(padBuf, binary.LittleEndian, bytes.Repeat([]byte{byte(n)}, n))
	fmt.Println("Padding Data : ", padBuf.Bytes(), "\n")
	return padBuf.Bytes()
	// padBuf := make([]byte, len(b)+n)
	// copy(padBuf, b)
	// copy(padBuf[len(b):], bytes.Repeat([]byte{byte(n)}, n))
	// return padBuf

}

func pkcs7Unpad(ciphertext []byte, blocksize int) []byte {
	a := ciphertext[len(ciphertext)-1]
	num := int(a)
	if num == 0 || num > len(ciphertext) {
		return nil
	}
	for i := 0; i < num; i++ {
		if ciphertext[len(ciphertext)-num+i] != a {
			return nil
		}
	}
	fmt.Println("Unpad Data : ", ciphertext[:len(ciphertext)-num], "\n")
	return ciphertext[:len(ciphertext)-num]
}
