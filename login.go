package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	type User struct {
		Username string
		Password string
	}

	var l_username string
	var l_password string
	var a_username string

	normal_login := flag.NewFlagSet("-l", flag.ExitOnError)
	normal_login.StringVar(&l_username, "u", "easteregg", "Please specify a username")
	normal_login.StringVar(&l_password, "p", "easteregg", "Please specify a password")

	anon_login := flag.NewFlagSet("-a", flag.ExitOnError)
	anon_login.StringVar(&a_username, "u", "anon", "Please specify a display name")

	argLength := len(os.Args[1:])
	switch os.Args[1] {
	case "-l":
		if argLength != 5 {
			fmt.Println("Usage: login.go -l -u (username) -p (password)")
			fmt.Println("OR     login.go -a -u (username)")
			os.Exit(1)
		}
		normal_login.Parse(os.Args[2:])
		fmt.Println("[ Silence ] - Attempting to login as: " + l_username + "..")
	case "-a":
		if argLength != 3 {
			fmt.Println("Usage: login.go -l -u (username) -p (password)")
			fmt.Println("OR     login.go -a -u (username)")
			os.Exit(1)
		}
		anon_login.Parse(os.Args[2:])
		fmt.Println("[ Silence ] - You are logged in as anon user: " + a_username)

		f, err := os.Create("loggedanon.dat")
		defer f.Close()
		if err != nil {
			log.Fatal("Failed to create file, exiting..")
			os.Exit(1)
		}
		os.Exit(1)

	default:
		fmt.Println("Usage: login.go -l -u (username) -p (password)")
		fmt.Println("OR     login.go -a -u (username)")
		os.Exit(1)
	}

	_, err := os.Stat("accounts.json")
	if os.IsNotExist(err) {
		fmt.Println("[ Silence ] - No account registered.. please use the command 'go run create_account.go -u (username) -p (password) -m c' to create an account")
		os.Exit(1)
	}

	sha256hash := sha256.New()
	sha256hash.Write([]byte(l_password))
	md := sha256hash.Sum(nil)
	password := hex.EncodeToString(md)

	contents, err := ioutil.ReadFile("accounts.json")
	if err != nil {
		log.Fatal("Error..")
	}

	user := User{l_username, password}
	res, err := json.Marshal(user)

	if string(res) != string(contents) {
		fmt.Println("[ Silence ] - Incorrect username:password combination... Try again..")
		os.Exit(1)
	}
	fmt.Println("[ Silence ] - Logged in as " + l_username)

	f, err := os.Create("logged.dat")
	defer f.Close()
	if err != nil {
		log.Fatal("Failed to create file, exiting..")
		os.Exit(1)
	}

}
