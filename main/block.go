package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"
)

var (
	maxnonce int32 = math.MaxInt32
)

//定义区块结构体

type Block struct {
	Version       int32
	PrevBlockHash []byte
	MerkleRoot    []byte
	Hash          []byte
	Time          int32
	Bits          int32
	Nonce         int32
	Transactions  []*Transation
	Height        int32
}

//根据前一个hash增加区块
func NewBlock(transations []*Transation, prevBlockHash []byte, height int32) *Block {
	block := &Block{
		2,
		prevBlockHash,
		[]byte{},
		[]byte{},
		int32(time.Now().Unix()),
		404454260,
		0,
		transations,
		height,
	}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Nonce = nonce
	block.Hash = hash
	return block
}

//创世区块
func NewGensisBlock(transactions []*Transation) *Block {
	block := &Block{
		2,
		[]byte{},
		[]byte{},
		[]byte{},
		int32(time.Now().Unix()),
		404454260,
		0,
		transactions,
		0,
	}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Nonce = nonce
	block.Hash = hash
	//	block.String()
	return block
}

//序列化
func (b *Block) Serialize() []byte {
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return encoded.Bytes()
}

//反序列化
func DeserializeBlock(d []byte) *Block {
	var block Block
	decode := gob.NewDecoder(bytes.NewReader(d))
	err := decode.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}

//打印区块
func (b *Block) String() {
	fmt.Printf("version:%s\n", strconv.FormatInt(int64(b.Version), 10))
	fmt.Printf("Prev.BlockHash:%x\n", b.PrevBlockHash)
	fmt.Printf("Prev.merkleroot:%x\n", b.MerkleRoot)
	fmt.Printf("Prev.Hash:%x\n", b.Hash)
	fmt.Printf("Time:%s\n", strconv.FormatInt(int64(b.Time), 10))
	fmt.Printf("Bits:%s\n", strconv.FormatInt(int64(b.Bits), 10))
	fmt.Printf("nonce:%s\n", strconv.FormatInt(int64(b.Nonce), 10))
}

//根据交易创建merkleROOT
func (b *Block) createMerkelTreeRoot(transations []*Transation) {
	var tranHash [][]byte
	for _, tx := range transations {
		tranHash = append(tranHash, tx.Hash())
	}
	mTree := NewMerkleTree(tranHash)
	b.MerkleRoot = mTree.RootNode.Data
}
