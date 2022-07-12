package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

var (
	conns   []net.Conn
	connCh  = make(chan net.Conn)
	closeCh = make(chan net.Conn)
	msgCh   = make(chan string)
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide port number")
		return
	}
	PORT := ":" + arguments[1]
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		for {
			c, err := l.Accept()
			if err != nil {
				fmt.Println(err)
				return
			}

			conns = append(conns, c)

			connCh <- c

		}
	}()

	for {
		select {
		case conn := <-connCh:
			go onMessage(conn)
		case msg := <-msgCh:
			fmt.Print(msg)
		case conn := <-closeCh:
			fmt.Println("Client exit")
			removeConn(conn)
			conn.Close()
		}
	}
}

func removeConn(conn net.Conn) {
	var i int
	for i := range conns {
		if conns[i] == conn {
			break
		}
	}
	conns = append(conns[i:], conns[:i+1]...)
}

func publicMsg(conn net.Conn, msg string) {
	for i := range conns {
		if conns[i] != conn {
			conns[i].Write([]byte(msg))
		}
	}
}

func onMessage(c net.Conn) {
	for {
		reader := bufio.NewReader(c)
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		msgCh <- msg
		publicMsg(c, msg)
	}
	closeCh <- c
}
