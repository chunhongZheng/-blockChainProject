package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const walletFile = "wallet.dat"

type Wallets struct {
	WalletsStore map[string]*Wallet
}

func NewWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.WalletsStore = make(map[string]*Wallet)
	err := wallets.LoadFromFile()
	return &wallets, err
}

func (ws *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := fmt.Sprintf("%s", wallet.GetAddress())
	ws.WalletsStore[address] = wallet
	return address
}
func (ws *Wallets) GetWallet(address string) Wallet {
	return *ws.WalletsStore[address]
}
func (ws *Wallets) getAddress() []string {
	var addresses []string
	for address, _ := range ws.WalletsStore {
		addresses = append(addresses, address)
	}
	return addresses
}

//反序列化
func (ws *Wallets) LoadFromFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}
	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Panic(err)
	}
	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Panic(err)
	}
	ws.WalletsStore = wallets.WalletsStore
	return nil
}

//序列化
func (ws *Wallets) SaveToFile() {
	var content bytes.Buffer
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}
	err = ioutil.WriteFile(walletFile, content.Bytes(), 0777)
	if err != nil {
		log.Panic(err)
	}
}
