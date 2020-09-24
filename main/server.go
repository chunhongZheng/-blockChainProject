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
type getBlocks struct {
	AddrFrom string
}

type inv struct {
	AddrFrom string
	Type     string
	Items    [][]byte
}
type getData struct {
	AddrFrom string
	Type     string
	ID       []byte
}

type blockSend struct {
	AddrFrom string
	Block    []byte
}

var knowNodes = []string{"localhost:3000"} //主节点
var nodeAddress string
var blockInTransit [][]byte

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
	case "getBlocks":
		handleGetBlock(request, bc)
	case "inv":
		handleInv(request, bc)
	case "getData":
		handleGetData(request, bc)
	case "block":
		handleBlock(request, bc)
	}
}

func handleBlock(request []byte, bc *BlockChain) {
	var buff bytes.Buffer
	var payLoad blockSend
	buff.Write(request[commonLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payLoad)
	if err != nil {
		log.Panic(err)
	}
	blockData := payLoad.Block
	block := DeserializeBlock(blockData)
	bc.AddBlock(block)
	fmt.Printf("Recieve a new Block")
	if len(blockInTransit) > 0 {
		blockHash := blockInTransit[0]
		sendGetData(payLoad.AddrFrom, "block", blockHash)
		blockInTransit = blockInTransit[1:]
	} else {
		//更新utxo
		set := UTXOSet{bc}
		set.Reindex()
	}
}

func handleGetData(request []byte, bc *BlockChain) {
	var buff bytes.Buffer
	var payLoad getData
	buff.Write(request[commonLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payLoad)
	if err != nil {
		log.Panic(err)
	}
	if payLoad.Type == "block" {
		block, err := bc.GetBlock([]byte(payLoad.ID))
		if err != nil {
			log.Panic(err)
		}
		sendBlock(payLoad.AddrFrom, &block)
	}
}

func sendBlock(addr string, block *Block) {
	data := blockSend{nodeAddress, block.Serialize()}
	payLoad := gobEncode(data)
	request := append(commandToBytes("block"), payLoad...)
	sendData(addr, request)
}

func handleInv(request []byte, bc *BlockChain) {
	var buff bytes.Buffer
	var payLoad inv
	buff.Write(request[commonLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payLoad)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Receieve inventory %d,%s,%s\n", len(payLoad.Items), payLoad.Type, payLoad.AddrFrom)

	if payLoad.Type == "block" {
		blockInTransit := payLoad.Items //所有区块hash集合
		blockHash := payLoad.Items[0]   //最近一个区块hash
		sendGetData(payLoad.AddrFrom, "block", blockHash)
		newTransit := [][]byte{} //剔除最近区块hash后的新的所有区块hash集合

		for _, b := range blockInTransit {
			if bytes.Compare(b, blockHash) != 0 {
				newTransit = append(newTransit, b)
			}
		}
		blockInTransit = newTransit
	}

}

func sendGetData(addr string, kind string, id []byte) {
	payLoad := gobEncode(getData{nodeAddress, kind, id})
	request := append(commandToBytes("getData"), payLoad...)
	sendData(addr, request)
}

/**处理请求区块数据*/
func handleGetBlock(request []byte, bc *BlockChain) {
	var buff bytes.Buffer
	var payLoad getBlocks
	buff.Write(request[commonLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payLoad)
	if err != nil {
		log.Panic(err)
	}
	block := bc.GetLockHash()
	sendInv(payLoad.AddrFrom, "block", block)
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
		fmt.Println("当前区块高度小于外部区块，需从外部更新数据")
		sendGetBlock(payLoad.AddrFrom)
	} else {
		//当前区块高度大于外部，向外部广播，通知其他子节点更新数据
		fmt.Println("当前区块高度大于外部，向外部外送数据")
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
func sendInv(addr string, kind string, item [][]byte) {
	inventory := inv{nodeAddress, kind, item}
	payLoad := gobEncode(inventory)
	request := append(commandToBytes("inv"), payLoad...)
	sendData(addr, request)
}
func sendGetBlock(address string) {
	//payLoad:=gobEncode(getBlocks{nodeAddress})
	payLoad := gobEncode(getBlocks{nodeAddress}) //向主节点获取block信息
	request := append(commandToBytes("getBlocks"), payLoad...)
	sendData(address, request)
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
