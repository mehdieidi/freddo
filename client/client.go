package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("udp", "127.0.0.1:1234")
	if err != nil {
		panic(err)
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			slog.Error("error closing connection", "err", err)
		}
	}(conn)

	err = writeTo(conn, "hey my dear peer!")
	if err != nil {
		panic(err)
	}

	go prompt(conn)

	for {
		buf := make([]byte, 2048)
		_, err = bufio.NewReader(conn).Read(buf)
		if err != nil {
			slog.Error("error reading from connection", "err", err)
			continue
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

		err := writeTo(conn, text)
		if err != nil {
			slog.Error("error writing to connection", "err", err)
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
