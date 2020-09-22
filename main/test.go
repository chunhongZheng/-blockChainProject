package main

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/ripemd160"
	"math/big"
)

//const gensisData="blockChainInitData"
//测试区块序列化
func TestNewSerialize() {
	//初始化区块
	block := &Block{
		2,
		[]byte{},
		[]byte{},
		[]byte{},
		1418755780,
		404454260,
		0,
		[]*Transation{},
	}
	deBlock := DeserializeBlock(block.Serialize())
	deBlock.String()
}

func TestCreateMerkleTreeRoot() {
	//初始化区块
	//初始化区块
	block := &Block{
		2,
		[]byte{},
		[]byte{},
		[]byte{},
		1418755780,
		404454260,
		0,
		[]*Transation{},
	}

	txin := TXInput{[]byte{}, -1, nil, []byte(gensisData)}
	txout := NewTxOutput(subsidy, "first")
	tx := Transation{nil, []TXInput{txin}, []TXOutput{*txout}}

	txin2 := TXInput{[]byte{}, -1, nil, []byte(gensisData)}
	txout2 := NewTxOutput(subsidy, "second")
	tx2 := Transation{nil, []TXInput{txin2}, []TXOutput{*txout2}}

	var Transations []*Transation
	Transations = append(Transations, &tx, &tx2)
	block.createMerkelTreeRoot(Transations)
	fmt.Printf("%x\n", block.Merkleroot)
}

func TestPow() {
	block := &Block{
		2,
		[]byte{},
		[]byte{},
		[]byte{},
		1418755780,
		404454260,
		0,
		[]*Transation{},
	}
	pow := NewProofOfWork(block)
	nonce, _ := pow.Run()
	block.Nonce = nonce
	fmt.Println("POW:", pow.validate())

}

//测试下载第三方包
func TestThirdPackage() {
	//github.com/golang/crypto v0.0.0-20180820150726-614d502a4dac // indirect
	//github.com/astaxie/beego-1.12.2 // indirect
	ripemd160.New()
}

func TestBoltDB() {
	blockchain := NewBlockChain("1GgxEy8DC3Q1w8vrCZ4FjLE8n26HGb6748")
	blockchain.MineBlock([]*Transation{})
	//	blockchain.AddBlock()
	blockchain.printBlockchain()
}

func TestNewGensisBlock() {
	//创世区块测试
	transation := NewCoinBaseTx("caspar", "blockChain init data") //生成矿工交易
	NewGensisBlock([]*Transation{transation})
}

func TestProofOfWorkRun() {
	secondhash, _ := hex.DecodeString("0000faf709ed3b20e39ddc1619095822c40774b197c0606c640a068bc46e971c")
	var currenthash big.Int
	currenthash.SetBytes(secondhash[:])
	fmt.Printf("%x\n", secondhash)
	//比较

}
func TestCLI() {
	bc := NewBlockChain("1GgxEy8DC3Q1w8vrCZ4FjLE8n26HGb6748")
	defer bc.db.Close()
	cli := CLI{bc}
	cli.Run()
}

//测试钱包
func TestWallet() {
	wallet := NewWallet()
	fmt.Printf("私钥:%x\n", wallet.PrivateKey.D.Bytes()) //打印私钥
	fmt.Printf("公钥:%x\n", wallet.PublicKey)            //打印公钥
	//	address,_:=hex.DecodeString(wallet.PrivateKey.D.Bytes())
	bitAddress := wallet.GetAddress() //接口操作
	//fmt.Printf("生成的比特币地址16进制为为:%x\n",bitAddress)
	fmt.Printf("生成的比特币地址为:%s\n", string(bitAddress))
	bool := IsVaildBitcoinAddress(string(bitAddress))
	fmt.Println(bool)
}
func main() {
	//	TestNewSerialize()
	//	TestCreateMerkleTreeRoot()
	//TestPow()
	//TestBoltDB()
	//	TestProofOfWorkRun()
	//	TestNewGensisBlock()
	TestCLI()

}
