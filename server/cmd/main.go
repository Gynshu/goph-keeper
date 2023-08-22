package main

import (
	"fmt"
	"log"

	"github.com/gynshu-one/goph-keeper/server/pkg/utils"
)

func main() {
	// Generate a master key for the user
	key, err := utils.GenerateMasterKeyForUser()
	if err != nil {
		log.Fatal(err)
	}

	// Encrypt some data using the master key
	plaintext := []byte(`ok`)
	ciphertext, err := utils.EncryptData(plaintext, key)
	if err != nil {
		log.Fatal(err)
	}

	// Decrypt the data using the master key
	decrypted, err := utils.DecryptData(ciphertext, key)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Key: %x\n", key)
	fmt.Printf("Plaintext: %s\n", plaintext)
	fmt.Printf("Ciphertext: %s\n", ciphertext)
	fmt.Printf("Decrypted: %s\n", decrypted)
}
