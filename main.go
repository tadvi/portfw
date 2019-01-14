package main

import (
	"io"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	if len(os.Args) <= 2 {
		log.Fatal("usage: portfw local:port remote:port")
	}
	localAddrString := os.Args[1]
	remoteAddrString := os.Args[2]

	localAddr, err := net.ResolveTCPAddr("tcp", localAddrString)
	if localAddr == nil {
		log.Fatalf("net.ResolveTCPAddr failed: %s", err)
	}
	local, err := net.ListenTCP("tcp", localAddr)
	if local == nil {
		log.Fatalf("portfw: %s", err)
	}
	log.Printf("portfw listen on %s", localAddr)

	for {
		conn, err := local.Accept()
		if conn == nil {
			log.Printf("accept failed: %s", err)
			continue
		}
		go forward(conn, remoteAddrString)
	}
}

/// forward requests to other host.
func forward(local net.Conn, remoteAddr string) {

	remote, err := net.DialTimeout("tcp", remoteAddr, time.Duration(5*time.Second))
	if remote == nil {
		log.Printf("remote dial failed: %v\n", err)
		local.Close()
		return
	}
	go func() {
		defer local.Close()
		io.Copy(local, remote)
	}()
	go func() {
		defer remote.Close()
		io.Copy(remote, local)
	}()
	log.Printf("forward %s to %s", local.RemoteAddr(), remoteAddr)
}
