package main


import (
	"fmt"
	"log"
	"net"
)

var (
	serverIP = "54.234.151.78:2050"
	eusouth = "52.47.150.186:2050"
)

type ProxyConnections struct {
	ClientHandle net.Conn
	ServerHandle net.Conn
	ServerPacket Packet
	ClientPacket Packet
	cChan chan Packet
	sChan chan Packet
}

func StartListener() {
	log.Println("Starting listeners...")
	lConn, err := net.Listen("tcp", "127.0.0.1:2050")
	proxy := ProxyConnections{}
	
	go StartPolicyListener()
	if err != nil {
		fmt.Println("Failure starting listener conn:", err)
	}
	log.Println("Listener Started!")
	cConn, err := lConn.Accept()
	if err != nil {
		fmt.Println("Failure accepting client:", err)
	}
	proxy.ClientHandle = cConn
	sConn, err := net.Dial("tcp", serverIP) //change server ip here
	if err != nil {
		fmt.Println("Error dialing server:", err)
	}
	proxy.ServerHandle = sConn
	fmt.Println("Connected to", sConn.RemoteAddr())
	proxy.cChan = make(chan Packet)
	proxy.sChan = make(chan Packet)
	go func() {
		for {
			proxy.ClientPacket = NewPacket()
			byteRead, _ := proxy.ClientHandle.Read(proxy.ClientPacket.Data)
			_ = byteRead
			proxy.ClientPacket.Length = int(proxy.ClientPacket.ReadUInt32())
			if proxy.ClientPacket.Length <= packetSize { //cut packet down to size
				proxy.ClientPacket.Data = proxy.ClientPacket.Data[:proxy.ClientPacket.Length]
			} else {
				fmt.Println("Got big packet!")
				fmt.Println(proxy.ClientPacket.Data[0:4])
			}
			fmt.Printf("Client packet size is %d\n", proxy.ClientPacket.Length)
			//proxy.ClientPacket.Data = proxy.ClientPacket.Data[:proxy.ClientPacket.Length]
			proxy.cChan <- proxy.ClientPacket //tell it that we are ready
		}
		
	}()
	go func() {
		for {
			proxy.ServerPacket = NewPacket()
			firstRead, _ := proxy.ServerHandle.Read(proxy.ServerPacket.Data[0:5])
			_ = firstRead
			proxy.ServerPacket.Length = int(proxy.ServerPacket.ReadUInt32())
			if proxy.ServerPacket.Length <= packetSize { //cut packet down to size
				proxy.ServerPacket.Data = proxy.ServerPacket.Data[:proxy.ServerPacket.Length]
			} else {
				fmt.Println("Got big packet!")
			}
			byteIndex := 0
			for byteIndex != proxy.ServerPacket.Length-5 { //loop until we have ALL of the packet bytes in the buffer
				secondRead, _ := proxy.ServerHandle.Read(proxy.ServerPacket.Data[byteIndex+5:proxy.ServerPacket.Length])
				byteIndex += secondRead
			}

			fmt.Printf("Server packet size is %d\n", proxy.ServerPacket.Length)
			//proxy.ServerPacket.Data = proxy.ServerPacket.Data[:proxy.ServerPacket.Length]
			proxy.sChan <- proxy.ServerPacket //tell it that we are ready
		}
		
	}()
	for {
		proxy.Balancer()
	}
	
}

func (p *ProxyConnections)Balancer() {
	select { //this is where we will decrypt data
	case data := <-p.cChan:
		fmt.Println("Client:", data.Data[data.Index:data.Length])
		p.ServerHandle.Write(data.Data[data.Index:data.Length])
	case data := <-p.sChan:
		fmt.Println("Server:", data.Data[:data.Length])
		p.ClientHandle.Write(data.Data[:data.Length])
	//default:
	//	fmt.Println("default")
	}
}

func (p *ProxyConnections)Sender(data []byte, code int) {
	switch code {
	case 1: //client
		p.ServerHandle.Write(data)
	case 2: //server
		p.ClientHandle.Write(data)
	}
}

func StartPolicyListener() {
	lConn, err := net.Listen("tcp", "127.0.0.1:843")
	if err != nil {
		fmt.Println("Error listening policy:", err)
	}
	cli, err := lConn.Accept()
	if err != nil {
		fmt.Println("Bad accept policy:", err)
	}
	buffer := make([]byte, 23)
	cli.Read(buffer)
	//correct := "3c706f6c6963792d66696c652d726571756573742f3e00"
	response := []byte{0x3c, 0x3f, 0x78, 0x6d, 0x6c, 0x20, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x3d, 0x22, 0x31, 0x2e, 0x30, 0x22, 0x3f, 0x3e, 0x3c, 0x21, 0x44, 0x4f, 0x43, 0x54, 0x59, 0x50, 0x45, 0x20, 0x63, 0x72, 0x6f, 0x73, 0x73, 0x2d, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2d, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x20, 0x53, 0x59, 0x53, 0x54, 0x45, 0x4d, 0x20, 0x22, 0x2f, 0x78, 0x6d, 0x6c, 0x2f, 0x64, 0x74, 0x64, 0x73, 0x2f, 0x63, 0x72, 0x6f, 0x73, 0x73, 0x2d, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2d, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x2e, 0x64, 0x74, 0x64, 0x22, 0x3e, 0x20, 0x20, 0x3c, 0x63, 0x72, 0x6f, 0x73, 0x73, 0x2d, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2d, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x3e, 0x20, 0x20, 0x3c, 0x73, 0x69, 0x74, 0x65, 0x2d, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x20, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x74, 0x74, 0x65, 0x64, 0x2d, 0x63, 0x72, 0x6f, 0x73, 0x73, 0x2d, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2d, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x69, 0x65, 0x73, 0x3d, 0x22, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x2d, 0x6f, 0x6e, 0x6c, 0x79, 0x22, 0x2f, 0x3e, 0x20, 0x20, 0x3c, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x2d, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x2d, 0x66, 0x72, 0x6f, 0x6d, 0x20, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x3d, 0x22, 0x2a, 0x22, 0x20, 0x74, 0x6f, 0x2d, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x3d, 0x22, 0x2a, 0x22, 0x20, 0x2f, 0x3e, 0x3c, 0x2f, 0x63, 0x72, 0x6f, 0x73, 0x73, 0x2d, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x2d, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x3e, 0x0a, 0x00}
	cli.Write(response)
	cli.Close()
}
