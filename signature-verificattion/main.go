package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"log"

	sslconfig "github.com/rohit-tambe/data-privacy/ssl-config"
)

func main() {
	sslConf := sslconfig.NewGetSSLKey("../rsa.private", "../rsa.public")
	msg := []byte("verifiable message")

	// Before signing, we need to hash our message
	// The hash is what we actually sign
	msgHash := sha512.New()
	_, err := msgHash.Write(msg)
	if err != nil {
		log.Fatalln(err)
	}
	msgHashSum := msgHash.Sum(nil)

	// In order to generate the signature, we provide a random number generator,
	// our private key, the hashing algorithm that we used, and the hash sum
	// of our message
	signature, err := rsa.SignPSS(rand.Reader, sslConf.PrivateKey, crypto.SHA512, msgHashSum, nil)
	if err != nil {
		log.Fatalln(err)
	}
	base64Signature := base64.StdEncoding.EncodeToString(signature)
	fmt.Println(base64Signature, "\n\n")

	// base64SignatureByte, err := base64.StdEncoding.DecodeString(fmt.Sprintf("%s,rohit", base64Signature))
	base64SignatureByte, err := base64.StdEncoding.DecodeString(base64Signature)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(string(base64SignatureByte), "\n\n")
	// To verify the signature, we provide the public key, the hashing algorithm
	// the hash sum of our message and the signature we generated previously
	// there is an optional "options" parameter which can omit for now
	err = rsa.VerifyPSS(sslConf.PublicKey, crypto.SHA512, msgHashSum, base64SignatureByte, nil)
	if err != nil {
		fmt.Println("could not verify signature: ", err)
		return
	}
	// If we don't get any error from the `VerifyPSS` method, that means our
	// signature is valid
	fmt.Println("signature verified")
}
