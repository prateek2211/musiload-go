package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"log"
	"math"
)

func Decrypt(enc []byte, key string, iv string) string {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		log.Fatal(err.Error())
	}
	size := base64.StdEncoding.DecodedLen(len(enc))
	ciphertext := make([]byte, (int(math.Ceil((float64(size) / 16))) * 16))
	_, err = base64.StdEncoding.Decode(ciphertext, enc)
	if err != nil {
		log.Fatal(err.Error())
	}
	if len(ciphertext)%aes.BlockSize != 0 {
		fmt.Println(len(ciphertext))
		fmt.Println(aes.BlockSize)
		panic("ciphertext is not a multiple of the block size")
	}
	fmt.Println(len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, []byte(iv))
	mode.CryptBlocks(ciphertext, ciphertext)
	return string(ciphertext[:size-1])
}
