package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type RegisterData struct {
	Username string
	Password string
}

func NewRegisterData(username, password string) RegisterData {
	return RegisterData{username, password}
}

var AllRegisterData []RegisterData

func LoadRegisterData() {
	data, err := os.ReadFile("database.json")
	if err == nil {
		decoder := json.NewDecoder(strings.NewReader(string(data)))
		err = decoder.Decode(&AllRegisterData)
	}
	return
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	username = username[:len(username)-1]
	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')
	password = password[:len(password)-1]

	LoadRegisterData()

	registerData := NewRegisterData(username, password)

	AllRegisterData := append(AllRegisterData, registerData)

	file, err := os.OpenFile("database.json", os.O_WRONLY|os.O_CREATE, 0666)
	if err == nil {
		defer file.Close()
		encoder := json.NewEncoder(file)
		encoder.Encode(AllRegisterData)
	} else {
		fmt.Println(err)
	}
}
