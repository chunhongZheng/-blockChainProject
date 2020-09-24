package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
)

const nodeVersion = 0x00
const commonLength = 12

type Version struct {
	Version    int32
	BestHeight int32
	AddrFrom   string
}

var knowNodes = []string{"localhost:3000"}
var nodeAddress string

func (v *Version) String() {
	fmt.Printf("version:%d\n", v.Version)
	fmt.Printf("Height:%d\n", v.BestHeight)
	fmt.Printf("AddrFrom:%s\n", v.AddrFrom)
}

func StartServer(nodeID, mineAddress string, bc *BlockChain) {
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	ln, err := net.Listen("tcp", nodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()
	//bc := NewBlockChain("1GgxEy8DC3Q1w8vrCZ4FjLE8n26HGb6748")
	if nodeAddress != knowNodes[0] {
		sendVersion(knowNodes[0], bc)
	}
	for {
		fmt.Println("a1")
		conn, err2 := ln.Accept()
		if err2 != nil {
			log.Panic(err2)
		}
		fmt.Println("a2")
		go handleConnection(conn, bc)
	}
}

func handleConnection(conn net.Conn, bc *BlockChain) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}
	//获取命令
	command := bytesToCommand(request[:commonLength])
	fmt.Println(command)
	switch command {
	case "version":
		//	fmt.Println("收到version")
		handleVersion(request, bc)
	}
}

func handleVersion(request []byte, bc *BlockChain) {
	var buff bytes.Buffer
	var payLoad Version
	buff.Write(request[commonLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payLoad)
	if err != nil {
		log.Panic(err)
	}
	payLoad.String()
	myBestHeight := bc.GetBestHeight()
	foreignerBestHeight := payLoad.BestHeight
	if myBestHeight < foreignerBestHeight {
		//当前区块高度小于外部区块，需从外部更新数据
	} else {
		//当前区块高度大于外部，向外部外送数据
		//fmt.Println("当前区块高度大于外部，向外部外送数据")
		sendVersion(payLoad.AddrFrom, bc)
	}
	if !nodeIsKnow(payLoad.AddrFrom) {
		knowNodes = append(knowNodes, payLoad.AddrFrom)
	}
}
func nodeIsKnow(addr string) bool {
	for _, node := range knowNodes {
		if node == addr {
			return true
		}
	}
	return false
}

func sendVersion(addr string, bc *BlockChain) {
	bestHeight := bc.GetBestHeight()
	payload := gobEncode(Version{nodeVersion, bestHeight, nodeAddress})
	request := append(commandToBytes("version"), payload...)
	sendData(addr, request)
}

func commandToBytes(command string) []byte {
	var bytes [commonLength]byte
	for i, c := range command {
		bytes[i] = byte(c)
	}
	return bytes[:]

}
func bytesToCommand(bytes []byte) string {
	var command []byte
	for _, b := range bytes {
		if b != 0x00 {
			command = append(command, b)
		}
	}
	return fmt.Sprintf("%s", command)
}

func sendData(addr string, data []byte) {
	con, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("%s is not available", addr)
		var updateNodes []string
		for _, node := range knowNodes {
			if node != addr {
				updateNodes = append(updateNodes, node)
			}
		}
		knowNodes = updateNodes
	}
	defer con.Close()
	_, err = io.Copy(con, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}

}

func gobEncode(data interface{}) []byte {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}
