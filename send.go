// Package kakscmd provides functions for writing commands directly to
// a Kakoune socket. Example program using the functions is provided in
// cmd/kak-raw-send package.
package kakscmd

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type commandMessage struct {
	// header

	messageType uint8
	messageSize uint32

	// command

	commandSize uint32
	command     []byte
}

const (
	msgTypeCommand = 0x02

	headerAndSizeFields = 4 + 4 + 1
	finalNewLineSize    = 1
)

// SocketPath returns path to the socket for given named session.
func SocketPath(session string) string {
	runtime := os.Getenv("XDG_RUNTIME_DIR")
	if os.Getenv("XDG_RUNTIME_DIR") == "" {
		return filepath.Join(os.TempDir(), "kakoune", session)
	}
	return filepath.Join(runtime, "kakoune", session)
}

func prepareMsg(cmd []byte) commandMessage {
	l := len(cmd) + finalNewLineSize
	return commandMessage{
		messageType: msgTypeCommand,
		messageSize: uint32(headerAndSizeFields + l),
		commandSize: uint32(l),
		command:     cmd,
	}
}

// Write command to w in format accepted by Kakoune socket.
func Write(w io.Writer, cmd string) (n int, err error) {
	cmdB := []byte(cmd)
	if strings.HasSuffix(cmd, "\n") {
		cmdB = cmdB[:len(cmdB)-1]
	}
	msg := prepareMsg(cmdB)
	if int(msg.commandSize) != len(cmdB)+1 {
		panic("library error: wrong command size")
	}

	// Using LittleEndian by default since AMD64 and ARM64 are both litte endian.

	err = binary.Write(w, binary.LittleEndian, msg.messageType)
	if err != nil {
		return
	}
	n++
	err = binary.Write(w, binary.LittleEndian, msg.messageSize)
	if err != nil {
		return
	}
	n += 4
	err = binary.Write(w, binary.LittleEndian, msg.commandSize)
	if err != nil {
		return
	}
	n += 4
	err = binary.Write(w, binary.LittleEndian, msg.command)
	if err != nil {
		return
	}
	n += int(msg.commandSize)
	err = binary.Write(w, binary.LittleEndian, byte('\n'))
	if err != nil {
		return
	}
	if n != int(msg.messageSize) {
		panic(fmt.Sprintf("library error: n: %d msg size %d", n, msg.messageSize))
	}
	return n, err
}
