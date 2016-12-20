package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"utils"
)

func main() {
	originalText := "encrypt this"
	fmt.Println(originalText)

	key := gen_sym_key()
	fmt.Println(key)

	// encrypt value to base64
	cryptoText := sym_encrypt(key, originalText)
	fmt.Println(cryptoText)

	// encrypt base64 crypto to original value
	text := sym_decrypt(key, cryptoText)
	fmt.Printf(text)
}
// genreate sym key
func gen_sym_key() []byte{

	//Define the key size, always the same
	key := make([]byte, 16)

	//this key needs to be saved and encrypted with dest public key
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}
	return key
}

// encrypt string to base64 crypto using AES
func sym_encrypt(key []byte, text string) string {
	// key := []byte(keyText)
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext)
}

// decrypt from base64 to decrypted string
func sym_decrypt(key []byte, cryptoText string) string {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s\n", ciphertext)
}
