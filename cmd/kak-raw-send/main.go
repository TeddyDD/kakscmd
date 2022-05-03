package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	"go.teddydd.me/kakscmd"
)

var (
	flagSocket = flag.String("session", "", "session")
	flagMsg    = flag.String("cmd", "", "command")
	flagDebug  = flag.Bool("debug", false, "dump debug info")
)

func main() {
	flag.Parse()

	if *flagSocket == "" || *flagMsg == "" {
		log.Fatal("both socket and cmd flags must be provided")
	}

	socketPath := kakscmd.SocketPath(*flagSocket)
	con, err := net.Dial("unix", socketPath)
	if err != nil {
		log.Fatalf("failed to dial socket: %s", err.Error())
	}
	defer con.Close()

	buf := &bytes.Buffer{}
	_, err = kakscmd.Write(buf, *flagMsg)
	if err != nil {
		log.Printf("failed to prepare msg: %s", err.Error())
		return
	}

	bytes := buf.Bytes()
	if *flagDebug {
		fmt.Println(hex.Dump(bytes))
		enc := hex.EncodeToString(bytes)
		enc = strings.ToUpper(enc)
		fmt.Println(enc)
	}

	n, err := con.Write(bytes)
	if err != nil {
		log.Printf("failed to write msg to socket: %s", err.Error())
		return
	}
	if *flagDebug {
		log.Printf("written %d bytes", n)
	}
}
