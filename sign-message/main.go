package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	sslconfig "github.com/rohit-tambe/data-privacy/ssl-config"
)

func main() {
	message := []byte("message to be signed")
	hashed := sha256.Sum256(message)
	sslConf := sslconfig.NewGetSSLKey("../rsa.private", "../rsa.public")
	signature, err := rsa.SignPKCS1v15(rand.Reader, sslConf.PrivateKey, crypto.SHA256, hashed[:])
	if err != nil {
		panic(err)
	}
	fmt.Println(string(base64.StdEncoding.EncodeToString(signature)))
}
