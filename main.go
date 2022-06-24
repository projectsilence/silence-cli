package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

type IP struct {
	Query string
}

func main() {
	/*
		_, err := os.Stat("logged.dat")
		if os.IsNotExist(err) {
			fmt.Println("[ Silence ] - Not logged in.. see login.go or create_account.go for more information")
			os.Exit(1)
		}

		url := "https://api.ipify.org?format=text"
		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		ip, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		resp.Body.Close()
	*/

	var p_type string
	var ip string
	var port string

	flag.StringVar(&p_type, "m", "l", "Please specify whether you want to listen or connect..")
	flag.StringVar(&ip, "ip", "192.168.0.1", "Please specify the listening IP..")
	flag.StringVar(&port, "p", "8199", "Please specify the listening Port..")
	flag.Parse()

	argLength := len(os.Args[1:])
	if argLength != 6 {
		fmt.Println("Usage: main.go -m l -ip (IP) -p (port)")
		fmt.Println("OR     main.go -m c -ip (IP) -p (port)")
		os.Exit(1)
	}
	switch os.Args[2] {
	case "l":
		listen, err := net.Listen("tcp", ip+":"+port)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		defer listen.Close()
		fmt.Println("[ Silence ] - Starting listener on " + ip + ":" + port)

		for {
			conn, err := listen.Accept()
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			defer conn.Close()

			fmt.Println("[ Silence ] - Connection Established..  CTRL^C TO EXIT")
			go handleRequest(conn)
			go sendMessage(conn)
		}
	case "c":
		conn, err := net.Dial("tcp", ip+":"+port)
		if err != nil {
			log.Fatalln(err)
		}
		defer conn.Close()
		fmt.Println("[ Silence ] - Connection Established..  CTRL^C TO EXIT")
		go handleRequest(conn)
		go sendMessage(conn)

	default:
		fmt.Println("Usage: main.go -m l -ip (IP) -p (port)")
		fmt.Println("OR     main.go -m c -ip (IP) -p (port)")
		os.Exit(1)
	}

}

func sendMessage(conn net.Conn) {
	var message string
	fmt.Println("1")
	fmt.Scanln(&message)
	if message != "test123" {
		conn.Write([]byte(message))
	} else {
		conn.Write([]byte("test123"))
		os.Exit(1)
	}
}

func handleRequest(conn net.Conn) {
	// store incoming data
	buffer := make([]byte, 1024)
	fmt.Println("2")
	message_in, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}
	if string(message_in) == "test123" {
		os.Exit(1)
	}
	fmt.Println(string(message_in))
}
