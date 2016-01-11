package main

import (
	"fmt"
	"net"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/pglotfel/udpserver/protocols"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error : ", err)
	}
}

func main() {

	serverAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8080")
	CheckError(err)

	localAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8079")
	CheckError(err)

	conn, err := net.DialUDP("udp", localAddr, serverAddr)
	CheckError(err)

	defer conn.Close()
	i := 0

	fmt.Println("Client connected to: ", serverAddr)

	for {
		t := test.Packet_Type(1)
		h := true
		msg := &test.Packet{Type: &t, Heartbeat: &test.Heartbeat{Heartbeat: &h}}
		i++
		buff, _ := proto.Marshal(msg)
		_, err := conn.Write(buff)
		CheckError(err)
		time.Sleep(time.Second * 1)
	}

	fmt.Println("Connection closed...")
}
