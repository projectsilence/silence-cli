package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"os"
)

func generate_keys() {
	log.Printf("[ Silence ] - Generating client keys..")
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		fmt.Printf("Error generating RSA Keys..\n")
		os.Exit(1)
	}

	publicKey := &privateKey.PublicKey
	var privateKeyBytes []byte = x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	pemFilePrivate, err := os.Create("private.pem")
	if err != nil {
		fmt.Printf("Error creating private.pem: %s \n", err)
		os.Exit(1)
	}

	err = pem.Encode(pemFilePrivate, privateKeyBlock)
	if err != nil {
		fmt.Printf("Error writing private.pem: %s \n", err)
		os.Exit(1)
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		fmt.Printf("Error creating publickey: %s \n", err)
		os.Exit(1)
	}

	publicKeyBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	pemFilePublic, err := os.Create("public.pem")
	if err != nil {
		fmt.Printf("Error creating public.pem: %s \n", err)
		os.Exit(1)
	}
	err = pem.Encode(pemFilePublic, publicKeyBlock)
	if err != nil {
		fmt.Printf("Error writing public.pem: %s \n", err)
		os.Exit(1)
	}

	log.Printf("[ Silence ] - Finished setting up client keys, exiting..")
}

func main() {
	argLength := len(os.Args[1:])
	if argLength != 6 {
		fmt.Println("Usage: go run create_account.go -u (username) -p (password) -m (o to override, c to create normally)")
		os.Exit(1)
	}

	type User struct {
		Username string
		Password string
	}

	var username string
	var passwordin string
	var override string
	var override_int int

	flag.StringVar(&username, "u", "easteregg", "Please specify a username")
	flag.StringVar(&passwordin, "p", "easteregg", "Please specify a password")
	flag.StringVar(&override, "m", "c", "Specify to override")
	flag.Parse()

	sha256hash := sha256.New()
	sha256hash.Write([]byte(passwordin))
	md := sha256hash.Sum(nil)
	password := hex.EncodeToString(md)

	if len(username) > 16 || len(username) < 5 {
		fmt.Println("[ Silence ] - Please make sure your username is between 5 and 16 chars")
		os.Exit(1)
	}

	_, err := os.Stat("accounts.json")
	if !os.IsNotExist(err) {
		if override != "o" {
			fmt.Println("[ Silence ] - Account already exists.. use -m o to override")
			os.Exit(1)
		} else {
			log.Println("[ Silence ] - Overriding current account..")
		}
	}
	if override_int == 1 {
		os.Exit(1)
	}
	f, err := os.Create("accounts.json")
	defer f.Close()
	if err != nil {
		log.Fatal("Failed to open file, exiting..")
		os.Exit(1)
	}

	user := User{username, password}
	res, err := json.Marshal(user)

	if err != nil {
		fmt.Println(err)
	}

	file, err := os.OpenFile("accounts.json", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	defer file.Close()
	if _, err := file.Write(res); err != nil {
		log.Fatal(err)
	}

	log.Println("[ Silence ] - Account registered, generating keys..")
	generate_keys()
	os.Exit(1)

}
