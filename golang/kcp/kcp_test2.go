package main

import (
	"github.com/xtaci/kcp-go"
	"io"
	"log"
	"time"
)

func main() {
	if listener, err := kcp.ListenWithOptions("127.0.0.1:8082", nil, 3, 2); err == nil {
		// spin-up the client
		go client()
		for {
			s, err := listener.AcceptKCP()
			if err != nil {
				log.Fatal(err)
			}
			go handleEcho2(s)
		}
	} else {
		log.Fatal(err)
	}
}

// handleEcho send back everything it received
func handleEcho2(conn *kcp.UDPSession) {
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}

		n, err = conn.Write(buf[:n])
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func client() {
	// wait for server to become ready
	time.Sleep(time.Second)

	// dial to the echo server
	if sess, err := kcp.DialWithOptions("127.0.0.1:8082", nil, 3, 2); err == nil {
		for {
			data := time.Now().String()
			buf := make([]byte, len(data))
			log.Println("sent:", data)
			if _, err := sess.Write([]byte(data)); err == nil {
				// read back the data
				if _, err := io.ReadFull(sess, buf); err == nil {
					log.Println("recv:", string(buf))
				} else {
					log.Fatal(err)
				}
			} else {
				log.Fatal(err)
			}
			time.Sleep(time.Second * 10)
		}
	} else {
		log.Fatal(err)
	}
}

