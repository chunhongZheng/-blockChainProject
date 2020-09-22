package main

import (
	"bytes"
	"crypto/sha256"
	"math/big"
)

//工作量证明
type ProofOfWork struct {
	block   *Block
	tartget *big.Int
}

const targetBits = 16

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	pow := &ProofOfWork{b, target}
	return pow
}
func (pow *ProofOfWork) prepareData(nonce int32) []byte {
	data := bytes.Join(
		[][]byte{
			Int32ToLittleEndianHex(pow.block.Version),
			pow.block.PrevBlockHash,
			pow.block.MerkleRoot,
			Int32ToLittleEndianHex(pow.block.Time),
			Int32ToLittleEndianHex(pow.block.Bits),
			Int32ToLittleEndianHex(nonce),
		},
		[]byte{},
	)

	return data
}

//挖矿
func (pow *ProofOfWork) Run() (int32, []byte) {
	var nonce int32
	var secondhash [32]byte
	nonce = 0
	var currenthash big.Int
	for nonce < maxnonce {
		//序列化
		data := pow.prepareData(nonce)
		//double hash
		firsthash := sha256.Sum256(data)
		secondhash = sha256.Sum256(firsthash[:])
		// fmt.Printf("secondhash:%x\n",secondhash)
		currenthash.SetBytes(secondhash[:])
		//比较
		if currenthash.Cmp(pow.tartget) == -1 {
			break
		} else {
			nonce++
		}
	}
	return nonce, secondhash[:]
}

func (pow *ProofOfWork) validate() bool {
	var hashInt big.Int
	data := pow.prepareData(pow.block.Nonce)
	firsthash := sha256.Sum256(data)
	secondhash := sha256.Sum256(firsthash[:])
	hashInt.SetBytes(secondhash[:])
	isValid := hashInt.Cmp(pow.tartget) == -1
	return isValid
}
