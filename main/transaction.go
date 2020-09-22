package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
)

const subsidy = 100 //矿工收益

type Transation struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}
type TXInput struct {
	TXId      []byte //上一个交易的id
	VOutIndex int
	Signature []byte
	PubKey    []byte //公钥
}

type TXOutput struct {
	Value         int
	PublicKeyHash []byte //公钥hash
}

type TXOutputs struct {
	Outputs []TXOutput
}

func (outputs TXOutputs) Serialize() []byte {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(outputs)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}
func DeserializeOutputs(data []byte) TXOutputs {
	var Outputs TXOutputs
	dec := gob.NewDecoder(bytes.NewReader(data))
	//err:=dec.Decode(Outputs)
	err := dec.Decode(&Outputs)
	if err != nil {
		log.Panic(err)
	}
	return Outputs
}
func (out *TXOutput) Lock(address []byte) {
	decodeAddress := Base58Decode(address)
	publicHash := decodeAddress[1 : len(decodeAddress)-4]
	out.PublicKeyHash = publicHash
}

//交易序列化
func (tx Transation) Serialize() []byte {
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	return encoded.Bytes()
}

//交易序列化进行hash运算，得出惟一hash值，进行交易标识
func (tx *Transation) Hash() []byte {
	txcopy := tx
	txcopy.ID = []byte{}
	hash := sha256.Sum256(txcopy.Serialize()) //slice类型切片
	return hash[:]
}

//根据金额与地址构造一个交易输出
func NewTxOutput(value int, address string) *TXOutput {
	txo := &TXOutput{value, nil}
	//txo.PublicKeyHash=[]byte(address)
	txo.Lock([]byte(address))
	return txo
}

//生成矿工交易
func NewCoinBaseTx(to, data string) *Transation {
	txIn := TXInput{[]byte{}, -1, nil, []byte(data)}
	txOutput := NewTxOutput(subsidy, to)
	tx := Transation{nil, []TXInput{txIn}, []TXOutput{*txOutput}}
	tx.ID = tx.Hash()
	return &tx
}
func (out *TXOutput) CanBeUnlockedWith(pubickeyHash []byte) bool {
	return bytes.Compare(out.PublicKeyHash, pubickeyHash) == 0
	//return string(out.PublicKeyHash)==unlockdata
}
func (in *TXInput) CanUnlockOutputWith(publickeyHash []byte) bool {
	lockinghash := GeneratePublicKeyHash(in.PubKey)
	return bytes.Compare(lockinghash, publickeyHash) == 0
	//return string(in.Signature)==unlockdata
}
func (tx Transation) IsCoinBase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].TXId) == 0 && tx.Vin[0].VOutIndex == -1
}

func NewUTXOTransation(from, to string, amount int, bc *BlockChain) *Transation {
	var inputs []TXInput
	var outputs []TXOutput
	wallets, err := NewWallets()
	if err != nil {
		log.Panic(err)
	}
	wallet := wallets.GetWallet(from)

	acc, validoutputs := bc.FindSpendableOutputs(GeneratePublicKeyHash(wallet.PublicKey), amount)
	//acc,validoutputs:=bc.FindSpendableOutputs(from ,amount)
	if acc < amount {
		log.Panic("Error: Not enough fund")
	}
	for txid, outs := range validoutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}
		//所有输出都需要变成输入，变成已花费的输出
		for _, out := range outs {
			//交易
			input := TXInput{txID, out, nil, []byte(wallet.PublicKey)} //构建输入
			inputs = append(inputs, input)
		}
	}
	outputs = append(outputs, *NewTxOutput(amount, to)) //构建输出
	if acc > amount {
		//构建余额输出
		outputs = append(outputs, *NewTxOutput(acc-amount, from))
	}
	tx := Transation{nil, inputs, outputs}
	tx.ID = tx.Hash()
	bc.SignTransation(&tx, wallet.PrivateKey) //对交易进行签名
	return &tx
}

func (tx *Transation) Sign(privateKey ecdsa.PrivateKey, prevTXs map[string]Transation) {
	if tx.IsCoinBase() {
		//矿工交易不需要处理，因为并没有上一个输出交易
		return
	}
	//检查过程
	for _, vin := range tx.Vin {
		if prevTXs[hex.EncodeToString(vin.TXId)].ID == nil {
			log.Panic("当前交易出错，当前交易的输入并没有找到其对应的前面交易的输出")
		}
	}
	txcopy := tx.TrimmedCopy()

	for inId, vin := range txcopy.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.TXId)] //获取当前交易的当前输入的前一笔交易
		txcopy.Vin[inId].Signature = nil
		txcopy.Vin[inId].PubKey = prevTx.Vout[vin.VOutIndex].PublicKeyHash //这笔交易的这笔输入的引用的有一笔交易的输出的公钥哈希
		txcopy.ID = txcopy.Hash()
		fmt.Printf("签名hash txcopy.ID%x\n", txcopy.ID)
		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, txcopy.ID) //对复制后的交易的hash值进行签名
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)
		tx.Vin[inId].Signature = signature
	}
}

func (tx *Transation) TrimmedCopy() Transation {
	var inputs []TXInput
	var outputs []TXOutput

	for _, vin := range tx.Vin {
		inputs = append(inputs, TXInput{vin.TXId, vin.VOutIndex, nil, nil})
	}
	for _, vout := range tx.Vout {
		outputs = append(outputs, TXOutput{vout.Value, vout.PublicKeyHash})
	}
	txCopy := Transation{tx.ID, inputs, outputs}
	return txCopy
}

func (tx *Transation) Verify(prevTXs map[string]Transation) bool {
	if tx.IsCoinBase() {
		return true
	}
	txcopy := tx.TrimmedCopy()
	//椭圆曲线
	curve := elliptic.P256()

	for inID, vin := range tx.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.TXId)]
		txcopy.Vin[inID].Signature = nil
		txcopy.Vin[inID].PubKey = prevTx.Vout[vin.VOutIndex].PublicKeyHash
		txcopy.ID = txcopy.Hash()

		r := big.Int{}
		s := big.Int{}
		signLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(signLen / 2)])
		s.SetBytes(vin.Signature[(signLen / 2):])
		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keyLen / 2)])
		y.SetBytes(vin.PubKey[(keyLen / 2):])
		rawPubkey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubkey, txcopy.ID, &r, &s) == false {
			return false
		}
		txcopy.Vin[inID].PubKey = nil
	}
	return true
}
