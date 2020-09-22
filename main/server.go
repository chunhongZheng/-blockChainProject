package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

type Version struct {
	Version    int32
	BestHeight int32
	AddrFrom   string
}

var knowNodes = []string{"localhost:3000"}
var nodeAddress = ""

func StartServer(nodeID, mineAddress string) {
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	ln, err := net.Listen("tcp", nodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()
	bc := NewBlockChain("1GgxEy8DC3Q1w8vrCZ4FjLE8n26HGb6748")
	if nodeAddress != knowNodes[0] {
		sendVersion(knowNodes[0], bc)
	}
	/*	for(
			conn , err2 :=ln.Accept()
		    if err2!=nil{
			log.Panic(err2)
		}
			go handleConnection(conn,bc)
			)*/
}

func handleConnection(conn interface{}, chain *BlockChain) {

}

func sendVersion(addr string, bc *BlockChain) {
	bestHeight := bc.GetBestHeight()
	payload := gobEncode(Version{nodeversion, bestHeight, nodeAddress})
	request := append([]byte("version"), payload...)
	sendData(addr, request)
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
