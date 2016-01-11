package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/pglotfel/udpserver/protocols"
)

//Add go routine to detect live clients.  Need to add channels for
//Different message types...

type Message struct {
	payload test.Packet
	address *net.UDPAddr
}

type UDPServer struct {
	//UDP server
	address         *net.UDPAddr
	connection      *net.UDPConn
	buffer          []byte
	clients         map[string]bool
	messageChannels map[string]chan Message
	waitGroup       *sync.WaitGroup
	closed          bool
}

func (u *UDPServer) Open(bufferSize int, port string, messageTypes []string) error {

	//Initialize internal channel and buffer
	u.messageChannels = make(map[string]chan Message)
	u.buffer = make([]byte, bufferSize)
	u.clients = make(map[string]bool)
	u.closed = false

	for _, v := range messageTypes {
		fmt.Println("CHANNEL TYPE: ", v)
		u.messageChannels[v] = make(chan Message)
	}

	//Create waitgroup for go routines.  We'll have 2 of them
	u.waitGroup = &sync.WaitGroup{}
	u.waitGroup.Add(2)

	//Initialize UDP address

	address, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		return err
	}
	u.address = address

	//Start listening to UDP socket
	u.connection, err = net.ListenUDP("udp", u.address)
	if err != nil {
		return err
	}

	//Process raw data
	go func() {
		defer func() {
			u.connection.Close()
			fmt.Println("Server successfully exited!")
			for _, v := range u.messageChannels {
				close(v)
			}
			//Signal that go routine is done
			u.waitGroup.Done()
		}()

		for {

			if u.closed {
				break
			}

			count, addr, err := u.connection.ReadFromUDP(u.buffer)
			if err != nil {
				u.closed = true
			}

			packet := &test.Packet{}

			err = proto.Unmarshal(u.buffer[0:count], packet)

			if err != nil {
				//Do something clever
			}

			//Put message in channel
			m := Message{payload: *packet, address: addr}
			fmt.Println(test.Packet_Type_name[int32(*packet.Type)], "ASFD")
			u.messageChannels[test.Packet_Type_name[int32(*packet.Type)]] <- m
		}
	}()

	//Do something with the data
	go func() {
		defer func() {
			fmt.Println("Channel closed!")
			//Signal that go routine is done
			u.waitGroup.Done()
		}()

		//Currently handling only HEARTBEAT
		for msg := range u.messageChannels["HEARTBEAT"] {
			fmt.Println("Received: ", msg.payload.GetHeartbeat(), " from: ", msg.address.String())
			if !u.clients[msg.address.String()] {
				u.clients[msg.address.String()] = true
			}
		}
	}()

	fmt.Println("Server successfully opened!")

	return nil
}

func (u *UDPServer) Close() {
	u.closed = true
	fmt.Println(u.clients)
	u.waitGroup.Wait()
}

func main() {

	server := &UDPServer{}
	server.Open(1024, ":8080", []string{"HEARTBEAT"})

	bufio.NewReader(os.Stdin).ReadBytes('\n')

	server.Close()

	//time.Sleep(time.Second * 1)
}
