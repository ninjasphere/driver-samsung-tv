package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net"

	"github.com/davecgh/go-spew/spew"
)

func sendCommand(ip, cmd string) error {

  log.Infof("Sending command %s to TV %s", cmd, ip)

	//ip := "10.10.0.219"
	mac := "00:00:00:02"
	name := "Ninja Sphere         " // Need extra spaces or bits of "samsung remote" get added on the end??

	appString := "iphone..iapp.samsung"
	tvAppString := "iphone.UN60D6000.iapp.samsung"
	port := 55000

	header := header(ip, mac, name, appString)
	command := command(cmd, tvAppString)

	spew.Dump("header", header)
	spew.Dump("command", command)

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		log.Warningf("Failed to connect to tv: %s", err)
		return err
	}

	conn.Write(header)
	conn.Write(command)
	conn.Close()
	return nil
}

func header(ip, mac, name, app string) []byte {
	msg := bytes.Buffer{}

	msg.WriteByte(0x64)
	msg.WriteByte(0x0)
	addB64(ip, &msg)
	addB64(mac, &msg)
	addB64(name, &msg)

	wrapped := bytes.Buffer{}

	wrapped.WriteByte(0x0)
	wrapped.WriteByte(uint8(len(app)))
	wrapped.WriteByte(0x0)
	wrapped.Write([]byte(app))
	wrapped.WriteByte(uint8(msg.Len()))
	wrapped.WriteByte(0x0)
	wrapped.Write(msg.Bytes())

	return wrapped.Bytes()
}

func command(command, app string) []byte {
	msg := bytes.Buffer{}

	msg.Write([]byte{0, 0, 0})
	addB64(command, &msg)

	wrapped := bytes.Buffer{}

	wrapped.WriteByte(0x0)
	wrapped.WriteByte(uint8(len(app)))
	wrapped.WriteByte(0x0)
	wrapped.Write([]byte(app))
	wrapped.WriteByte(uint8(msg.Len()))
	wrapped.WriteByte(0x0)
	wrapped.Write(msg.Bytes())

	return wrapped.Bytes()
}

func addB64(str string, msg *bytes.Buffer) {
	enc := []byte(base64.StdEncoding.EncodeToString([]byte(str)))

	msg.WriteByte(uint8(len(enc)))
	msg.WriteByte(0)
	msg.Write(enc)
}
