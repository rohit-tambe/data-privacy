package sslconfig

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"

	"golang.org/x/crypto/ssh"
)

type SSLConfig struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func NewGetSSLKey(privateKeyPath, publicKeyPath string) SSLConfig {
	return SSLConfig{
		PrivateKey: readPrivateKey(privateKeyPath),
		PublicKey:  readPublicKey(publicKeyPath),
	}

}
func CheckError(e error) {
	if e != nil {
		fmt.Println(e.Error())
	}
}
func readPrivateKey(privateKeyPath string) *rsa.PrivateKey {
	privateByte, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatalln(err)
	}
	privateK, err := ssh.ParseRawPrivateKey(privateByte)
	if err != nil {
		log.Fatalln(err)
	}

	return privateK.(*rsa.PrivateKey)
}

func readPublicKey(keyPath string) *rsa.PublicKey {
	pub, err := ioutil.ReadFile(keyPath)
	CheckError(err)
	pubPem, _ := pem.Decode(pub)

	parsedKey, err := x509.ParsePKIXPublicKey(pubPem.Bytes)
	CheckError(err)

	// var pubKey *rsa.PublicKey
	return parsedKey.(*rsa.PublicKey)
}
