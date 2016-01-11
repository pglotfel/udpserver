package main

import (
  "fmt"
  "net"
  "os"
  "github.com/golang/protobuf/proto"
  "github.com/pglotfel/hello/protocols"
)

type message struct {
  payload test.Packet
  count int
  address *net.UDPAddr
}

func CheckError(err error) {
  if err != nil {
    fmt.Println("Error : ", err)
    os.Exit(0)
  }
}

func ProcessMessages(c chan message) {

  defer fmt.Println("Done!")

  for msg := range c {
    fmt.Println("Received: ", msg.payload.GetMessage(), " from: ", msg.address)
  }
}

func main() {

  channel := make(chan message)

  serverAddr, err := net.ResolveUDPAddr("udp", ":8080")
  CheckError(err)

  serverConn, err := net.ListenUDP("udp", serverAddr)
  //Handle error
  CheckError(err)

  //Happen after everything else...
  defer func() {
          fmt.Println("Server exited!")
          serverConn.Close()
        }()

  //Buffer to hold incoming data...
  buff := make([]byte, 1024)
  test := &test.Packet{}

  fmt.Printf("Server initialized\n")

  go ProcessMessages(channel);

  for {
    n,addr,err := serverConn.ReadFromUDP(buff)
    //Process
    CheckError(err)
    err = proto.Unmarshal(buff[0:n], test)
    CheckError(err)
    m := message{payload: *test, count: n, address: addr}

    if m.count > 3 {
      fmt.Println(m.count)
      close(channel)
      break
    } else {
      channel <- m
    }
  }
}
