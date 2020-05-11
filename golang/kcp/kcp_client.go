package main

import (
	"fmt"
	"github.com/xtaci/kcp-go"
	"strings"
	"time"
)

const serverPortEcho = "127.0.0.1:8081"

func dialEcho() (*kcp.UDPSession, error) {
	conn, err := kcp.Dial(serverPortEcho)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return conn.(*kcp.UDPSession), err
}

func main() {
	cli, err := dialEcho()
	if err != nil {
		panic(err)
	}
	cli.SetStreamMode(false)
	cli.SetWindowSize(4096, 4096)
	cli.SetWriteDelay(true)
	cli.SetACKNoDelay(false)
	// NoDelay options
	// fastest: ikcp_nodelay(kcp, 1, 20, 2, 1)
	// nodelay: 0:disable(default), 1:enable
	// interval: internal update timer interval in millisec, default is 100ms
	// resend: 0:disable fast resend(default), 1:enable fast resend
	// nc: 0:normal congestion control(default), 1:disable congestion control
	cli.SetNoDelay(1, 100, 2, 0)
	const N = 100

	sb := strings.Builder{}
	for i := 0; i < 1000; i++ {
		sb.WriteString("hello")
	}
	for i := 0; i < N; i++ {
		time.Sleep(1 * time.Second)
		msg := fmt.Sprintf("%v", i) + sb.String()
		cli.Write([]byte(msg))
		buf := make([]byte, 10000)
		if n, err := cli.Read(buf); err == nil {
			if string(buf[:n]) != msg {
				fmt.Println("不一致", len(string(buf[:n])), len([]byte(msg)))
			}
		} else {
			panic(err)
		}

	}
	cli.Close()
}

