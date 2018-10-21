package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"runtime"
	"strings"
	"time"
)

//GameConnection is the base connection to the server and client
type GameConnection struct {
	InLoop       bool
	Running      bool
	Killed       bool
	SocketDebug  bool
	DidLogin     bool
	Phase        int
	LocalHandle  net.Conn
	RemoteHandle net.Conn
	LocalSocket  *bufio.ReadWriter
	RemoteSocket *bufio.ReadWriter
}

var (
	remoteHost   string
	remotePort   string
	localAddress = "127.0.0.1"
	useSocks     = false
)

//ConnectionLoop is the main loop for reading from the sockets
func (g *GameConnection) ConnectionLoop() {
	for {
		if g.Killed == false && g.Running == false {
			g.Running = true

			if g.InLoop == false {
				g.InLoop = true
				go g.Receive(g.LocalSocket, true)
				go g.Receive(g.RemoteSocket, false)
			}
		} else {
			time.Sleep(30 * time.Second)
		}
	}
}

//Receive reads from the socket
func (g *GameConnection) Receive(c *bufio.ReadWriter, isLocal bool) {
	for g.Running == true {
		recbuf := make([]byte, 1028)
		bytesRead, err := c.Read(recbuf)
		if err != nil {
			log.Println("Error reading:", err)
			if strings.Contains(err.Error(), "EOF") == true {
				g.InLoop = false
				g.Running = false
			}
			return
		}
		fmt.Println(recbuf[:bytesRead])

		//send out
		p := new(Packet)
		p.Data = make([]byte, len(recbuf))
		p.Data = recbuf[:bytesRead]
		if isLocal == true {
			g.Send(g.RemoteSocket, p)
		} else {
			g.Send(g.LocalSocket, p)
		}
	}
}

//Send should be self explantory
func (g *GameConnection) Send(c *bufio.ReadWriter, p *Packet) {
	if p == nil {
		return
	}
	if g.SocketDebug == true {
		fmt.Println("Send:", p.Data)
	}
	if c.Reader == nil || c.Writer == nil { //dont write anything if we've disconnected
		return
	}
	_, err := c.Writer.Write(p.Data)
	if err != nil {
		fmt.Println("Write error:", err)
		return
	}
	err = c.Flush()
	if err != nil {
		fmt.Println("Flush error:", err)
		return
	}
}

func InitListener(g *GameConnection) {
	lstn, err := net.Listen("tcp", localAddress+":"+remotePort)
	if err != nil {
		log.Println("Error starting listener:", err)
	}
	for {
		conn, err := lstn.Accept()
		if err != nil {
			log.Println("Error accepting client:", err)
		}
		g.WrapSocket(conn, true)
		InitGameConnection(g)
		go g.ConnectionLoop()
	}
}

//InitGameConnection connects to the server
func InitGameConnection(g *GameConnection) {
	if useSocks == true {

	} else {
		hostIP, err := net.ResolveTCPAddr("tcp", remoteHost+":"+remotePort)
		if err != nil {
			log.Println("Error resolving server:", err)
			return
		}
		conn, err := net.DialTCP("tcp", nil, hostIP)
		if err != nil {
			log.Println("Error dialing server:", err)
			return
		}
		err = conn.SetNoDelay(true)
		if err != nil {
			log.Println("Error setting no delay:", err)
			//shouldnt have to return since we still have a valid socket
		}
		g.WrapSocket(conn, false)
		log.Println("Dialed server")
	}
}

//Kill destroys the connection and cleans up
func (g *GameConnection) Kill() {
	g.Killed = true
	g.InLoop = false
	g.Running = false
	g.LocalHandle.Close()
	g.RemoteHandle.Close()
	g.LocalSocket = nil
	g.RemoteSocket = nil
	runtime.GC()
	time.Sleep(1000 * time.Millisecond)
}

//WrapSocket self explanatory
func (g *GameConnection) WrapSocket(c net.Conn, isLocal bool) {
	if isLocal == true {
		g.LocalHandle = c
		g.LocalSocket = bufio.NewReadWriter(bufio.NewReader(c), bufio.NewWriter(c))
	} else {
		g.RemoteHandle = c
		g.RemoteSocket = bufio.NewReadWriter(bufio.NewReader(c), bufio.NewWriter(c))
	}

}
