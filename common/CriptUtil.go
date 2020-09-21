package common

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"os"
)

// RSA
const (
	rsaKeySize    = 8192
	rsaClientSize = 2048
)

type Keypair struct {
	Priv *rsa.PrivateKey
	Pub  *rsa.PublicKey
}

func GenerateKeypair(client bool) *Keypair {
	var err error
	var size = rsaKeySize
	if client {
		size = rsaClientSize
	}
	fmt.Println("Size: ", size)
	priv, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return nil
	}
	var pub = &priv.PublicKey
	var kp = &Keypair{
		Priv: priv,
		Pub:  pub,
	}
	return kp
}

func getHash(client bool) hash.Hash {
	if client {
		return sha256.New()
	} else {
		return sha512.New()
	}
}

func Encrypt(message []byte, pub *rsa.PublicKey, client bool) []byte {
	if message == nil {
		fmt.Println("message null")
		return nil
	} else if pub == nil {
		fmt.Println("pub null")
		return nil
	}
	var err error
	fmt.Println(len(message))
	ciphertext, err := rsa.EncryptOAEP(getHash(client), rand.Reader, pub, message, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from encryption: %s\n", err)
		return nil
	}
	// Since encryption is a randomized function, ciphertext will be
	// different each time.
	return ciphertext
}

func Decrypt(cipherText []byte, priv *rsa.PrivateKey, client bool) []byte {
	if cipherText == nil {
		fmt.Println("cipherText null")
		return nil
	} else if priv == nil {
		fmt.Println("priv null")
		return nil
	}
	message, err := rsa.DecryptOAEP(getHash(client), rand.Reader, priv, cipherText, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from decryption: %s\n", err)
		return nil
	}
	return message
}

/*
func sign(message []byte, priv *rsa.PrivateKey) []byte {
	var err error
	signedMessage, err := rsa.SignPKCS1v15(rand.Reader, priv, message)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from encryption: %s\n", err)
		return nil
	}
	return signedMessage
}

func verify(signedMessage []byte, pub *rsa.PublicKey) {
	msgVerified, err := rsa.VerifyPKCS1v15(rand.Reader, pub, signedMessage)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from decryption: %s\n", err)
		return
	}
	fmt.Printf("Verified Message: %s\n", string(msgVerified))
}
*/
