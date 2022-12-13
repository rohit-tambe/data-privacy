package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"

	"golang.org/x/crypto/ssh"
)

func CheckError(e error) {
	if e != nil {
		fmt.Println(e.Error())
	}
}
func RSA_OAEP_Decrypt(cipherText string, privKey rsa.PrivateKey) string {
	ct, _ := base64.StdEncoding.DecodeString(cipherText)
	label := []byte("OAEP Encryptedd")
	rng := rand.Reader
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rng, &privKey, ct, label)
	CheckError(err)
	fmt.Println("Plaintext:", string(plaintext))
	return string(plaintext)
}
func RSA_OAEP_Encrypt(secretMessage string, key rsa.PublicKey) string {
	label := []byte("OAEP Encrypted")
	rng := rand.Reader
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, &key, []byte(secretMessage), label)
	CheckError(err)
	return base64.StdEncoding.EncodeToString(ciphertext)
}
func main() {
	// The GenerateKey method takes in a reader that returns random bits, and
	// the number of bits
	// privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	privateByte, err := ioutil.ReadFile("rsa.private")
	if err != nil {
		log.Fatalln(err)
	}
	privateK, err := ssh.ParseRawPrivateKey(privateByte)
	if err != nil {
		log.Fatalln(err)
	}

	privateKey := privateK.(*rsa.PrivateKey)

	// publicByte, err := ioutil.ReadFile("rsa.public")
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// public, err := ssh.

	// The public key is a part of the *rsa.PrivateKey struct
	// publicKey := privateKey.PublicKey
	// publicKey := public
	pub, err := ioutil.ReadFile("rsa.public")
	CheckError(err)
	pubPem, _ := pem.Decode(pub)

	parsedKey, err := x509.ParsePKIXPublicKey(pubPem.Bytes)
	CheckError(err)

	// var pubKey *rsa.PublicKey
	pubKey := parsedKey.(*rsa.PublicKey)
	cypherText := RSA_OAEP_Encrypt("this is supr secret key", *pubKey)
	fmt.Println(cypherText)
	plainText := RSA_OAEP_Decrypt(cypherText, *privateKey)
	fmt.Println(plainText)

	msg := []byte("verifiable message")

	// Before signing, we need to hash our message
	// The hash is what we actually sign
	msgHash := sha512.New()
	_, err = msgHash.Write(msg)
	if err != nil {
		log.Fatalln(err)
	}
	msgHashSum := msgHash.Sum(nil)

	// In order to generate the signature, we provide a random number generator,
	// our private key, the hashing algorithm that we used, and the hash sum
	// of our message
	signature, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA512, msgHashSum, nil)
	if err != nil {
		log.Fatalln(err)
	}

	// To verify the signature, we provide the public key, the hashing algorithm
	// the hash sum of our message and the signature we generated previously
	// there is an optional "options" parameter which can omit for now
	err = rsa.VerifyPSS(pubKey, crypto.SHA512, msgHashSum, signature, nil)
	if err != nil {
		fmt.Println("could not verify signature: ", err)
		return
	}
	// If we don't get any error from the `VerifyPSS` method, that means our
	// signature is valid
	fmt.Println("signature verified")
}
