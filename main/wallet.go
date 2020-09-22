package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"log"
)

const VERSION = byte(0x00)
const CHECKSUM_LENGTH = 4

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte
}

func (b *Wallet) newKeyPair() {
	curve := elliptic.P256()
	var err error
	b.PrivateKey, err = ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	b.PublicKey = append(b.PrivateKey.PublicKey.X.Bytes(), b.PrivateKey.PublicKey.Y.Bytes()...)
}
func NewWallet() *Wallet {
	b := &Wallet{nil, nil}
	b.newKeyPair()
	return b
}
func GeneratePublicKeyHash(publicKey []byte) []byte {
	sha256PubKey := sha256.Sum256(publicKey)
	r := ripemd160.New()
	r.Write(sha256PubKey[:])
	ripPubKey := r.Sum(nil)
	return ripPubKey
}

func CheckSumHash(versionPublickeyHash []byte) []byte {
	versionPublickeyHashSha1 := sha256.Sum256(versionPublickeyHash)
	versionPublickeyHashSha2 := sha256.Sum256(versionPublickeyHashSha1[:])
	tailHash := versionPublickeyHashSha2[:CHECKSUM_LENGTH]
	return tailHash
}

//获取地址

func (b *Wallet) GetAddress() []byte {
	//1.ripemd160(sha256(publickey))
	ripPubKey := GeneratePublicKeyHash(b.PublicKey)
	//2.最前面添加一个字节的版本信息获得 versionPublickeyHash
	versionPublickeyHash := append([]byte{VERSION}, ripPubKey[:]...)
	//3.sha256(sha256(versionPublickeyHash))  取最后四个字节的值
	tailHash := CheckSumHash(versionPublickeyHash)
	//4.拼接最终hash versionPublickeyHash + checksumHash
	finalHash := append(versionPublickeyHash, tailHash...)
	//进行base58加密
	address := Base58Encode(finalHash)
	return address
}

//检测比特币地址是否有效
func IsVaildBitcoinAddress(address string) bool {
	adddressByte := []byte(address)
	fullHash := Base58Decode(adddressByte)
	if len(fullHash) != 25 {
		return false
	}
	prefixHash := fullHash[:len(fullHash)-CHECKSUM_LENGTH]
	tailHash := fullHash[len(fullHash)-CHECKSUM_LENGTH:]
	tailHash2 := CheckSumHash(prefixHash)
	if bytes.Compare(tailHash, tailHash2[:]) == 0 {
		return true
	} else {
		return false
	}
}

//通过地址获得公钥
func GetPublicKeyHashFromAddress(address string) []byte {
	addressBytes := []byte(address)
	fullHash := Base58Decode(addressBytes)
	publicKeyHash := fullHash[1 : len(fullHash)-CHECKSUM_LENGTH]
	return publicKeyHash
}
