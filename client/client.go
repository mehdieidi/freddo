package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("udp", "127.0.0.1:1234")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	if err = writeTo(conn, "hey my deer peer!"); err != nil {
		panic(err)
	}

	go prompt(conn)

	for {
		buf := make([]byte, 2048)
		if _, err = bufio.NewReader(conn).Read(buf); err != nil {
			panic(err)
		}

		go handle(conn, buf)
	}
}

func handle(conn net.Conn, buf []byte) {
	fmt.Printf("[from %v]: %s\n", conn.RemoteAddr(), buf)
	prompt(conn)
}

func prompt(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("-> ")

		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if err := writeTo(conn, text); err != nil {
			panic(err)
		}
	}
}

func writeTo(conn net.Conn, msg string) error {
	_, err := conn.Write([]byte(msg))
	if err != nil {
		return err
	}
	return nil
}
