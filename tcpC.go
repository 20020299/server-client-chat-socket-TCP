package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

type ConfigData struct {
	Username string
	Password string
}

var Config []ConfigData

func LoadConfig() {
	data, err := os.ReadFile("database.json")
	if err == nil {
		decoder := json.NewDecoder(strings.NewReader(string(data)))
		err = decoder.Decode(&Config)
	}
	return
}

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide host:port.")
		return
	}

	CONNECT := arguments[1]
	c, err := net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Username: ")
	userName, _ := reader.ReadString('\n')
	userName = userName[:len(userName)-1]
	fmt.Print("Password: ")
	passWord, _ := reader.ReadString('\n')
	passWord = passWord[:len(passWord)-1]
	LoadConfig()

	for i := range Config {
		if userName == Config[i].Username && passWord == Config[i].Password {
			fmt.Println("************MESSAGE************")
			go onMsg(c)
			for {
				msgReader := bufio.NewReader(os.Stdin)

				msg, err := msgReader.ReadString('\n')
				if err != nil {
					break
				}

				msg = fmt.Sprintf("%v: %v\n", userName, msg)
				c.Write([]byte(msg))
			}
		}
	}

	c.Close()

}

func onMsg(conn net.Conn) {
	for {
		reader := bufio.NewReader(conn)
		msg, _ := reader.ReadString('\n')
		fmt.Print(msg)
	}
	conn.Close()
}
