package handlers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"os"
)

func outputPrivateKey(w io.Writer, key *rsa.PrivateKey) error {
	_, err := w.Write(x509.MarshalPKCS1PrivateKey(key))
	return err
}

func outputPublicKey(w io.Writer, key *rsa.PublicKey) error {
	//Version of tomcrypt packet
	const TOMCRYPT_VERSION = 0x98
	//Section tags
	const PACKET_SECT_RSA = 0
	//Subsection Tags for the first three sections
	const PACKET_SUB_KEY = 0
	//Public key type
	const PK_PRIVATE = 0           //PK private keys
	const PK_PUBLIC = 1            //PK public keys
	const PK_PRIVATE_OPTIMIZED = 2 //PK private key [rsa optimized]

	var keyType uint8 = PK_PUBLIC

	if err := binary.Write(w, binary.LittleEndian, uint16(TOMCRYPT_VERSION)); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, uint8(PACKET_SECT_RSA)); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, uint8(PACKET_SUB_KEY)); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, uint8(keyType)); err != nil {
		return err
	}

	nBytes := key.N.Bytes()
	if err := binary.Write(w, binary.LittleEndian, uint32(len(nBytes))); err != nil {
		return err
	}
	if _, err := w.Write(nBytes); err != nil {
		return err
	}

	bigE := big.NewInt(int64(key.E))
	eBytes := bigE.Bytes()

	if err := binary.Write(w, binary.LittleEndian, uint32(len(eBytes))); err != nil {
		return err
	}
	if _, err := w.Write(eBytes); err != nil {
		return err
	}

	return nil
}

func generateRsaKey() *rsa.PrivateKey {
	const RSA_KEY_SIZE = 1024

	fmt.Printf("Generating new %d-bit RSA key...", RSA_KEY_SIZE)
	privKey, err := rsa.GenerateKey(rand.Reader, RSA_KEY_SIZE)

	if err != nil {
		panic("Error generating RSA key :" + err.Error())
	}

	fmt.Println("DONE")
	return privKey
}

func LoadOrGenerateRsaKey(privateKeyName string, publicKeyName string) *rsa.PrivateKey {
	var loadedFromFile bool
	var rsaKey *rsa.PrivateKey
	{
		privBytes, err := ioutil.ReadFile(privateKeyName)
		if err != nil {
			fmt.Printf("err loading private key file: %s ", err.Error())
		} else {
			rsaKey, err = x509.ParsePKCS1PrivateKey(privBytes)
			if err != nil {
				fmt.Printf("err parsing privkey: %s ")
			} else {
				loadedFromFile = true
			}
		}
	}
	if rsaKey == nil {
		rsaKey = generateRsaKey()
		loadedFromFile = false
	}
	if !loadedFromFile {
		privFile, err := os.Create(privateKeyName)
		if err != nil {
			panic("error opening RSA privkey: " + err.Error())
		}
		err = outputPrivateKey(privFile, rsaKey)
		privFile.Close()
		if err != nil {
			panic("error writing RSA privkey: " + err.Error())
		}
	}
	var pubKeyExists bool = true
	if _, err := os.Stat(publicKeyName); err != nil {
		if os.IsNotExist(err) {
			pubKeyExists = false
		}
	}
	if !loadedFromFile || !pubKeyExists {
		pubFile, err := os.Create(publicKeyName)
		if err != nil {
			panic("error opening RSA pubkey: " + err.Error())
		}
		err = outputPublicKey(pubFile, &rsaKey.PublicKey)
		pubFile.Close()
		if err != nil {
			panic("error writing RSA pubkey: " + err.Error())
		}
	}

	return rsaKey
}
