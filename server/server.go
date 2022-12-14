package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

var peersMap = make(map[string]int)
var peersAddr = []*net.UDPAddr{}
var mutex sync.RWMutex

func main() {
	addr := net.UDPAddr{
		Port: 1234,
		IP:   net.ParseIP("127.0.0.1"),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		panic(err)
	}

	go prompt(conn)

	peerIndex := 0

	for {
		buf := make([]byte, 2048)
		n, peerAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		peerKey := peerAddr.IP.String() + strconv.Itoa(peerAddr.Port)

		mutex.Lock()
		if _, ok := peersMap[peerKey]; !ok {
			peersMap[peerKey] = peerIndex
			peersAddr = append(peersAddr, peerAddr)

			peerIndex++
		}
		mutex.Unlock()

		go handle(conn, peerAddr, buf, n)
	}
}

func handle(conn *net.UDPConn, addr *net.UDPAddr, buf []byte, length int) {
	msg := string(buf[:length])

	if msg == "#status" {
		text := fmt.Sprintf("me: [%v] you: [%v:%v]\n", conn.LocalAddr(), addr.IP, addr.Port)
		conn.WriteToUDP([]byte(text), addr)
	} else {
		fmt.Printf("\nfrom [%v:%v]: %s\n", addr.IP, addr.Port, buf)
	}
}

func prompt(conn *net.UDPConn) {
	var peerIDSet bool
	var peerID string

	reader := bufio.NewReader(os.Stdin)

	for {
		prompt := "-> "
		if peerIDSet {
			prompt = "[to peer " + peerID + "]" + " -> "
		}
		fmt.Print(prompt)

		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if peerIDSet {
			peerIDSet = false

			pID, _ := strconv.Atoi(peerID)
			writeTo(conn, pID, text)
		} else if text == "#peers" {
			printPeers()
		} else if len(text) > 9 && text[:9] == "#peer_id " {
			peerID = text[9:]
			peerIDSet = true
		}
	}
}

func printPeers() {
	mutex.RLock()
	defer mutex.RUnlock()

	for _, id := range peersMap {
		fmt.Printf("[%v] -> [%v:%v]  ", id, peersAddr[id].IP, peersAddr[id].Port)
	}
	fmt.Println()
}

func writeTo(conn *net.UDPConn, peerID int, text string) {
	mutex.RLock()
	defer mutex.RUnlock()

	conn.WriteToUDP([]byte(text), peersAddr[peerID])
}
