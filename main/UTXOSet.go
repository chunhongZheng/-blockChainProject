package main

import (
	"encoding/hex"
	"github.com/boltdb/bolt"
	"log"
)

type UTXOSet struct {
	blockchain *Blockchain
}

const utxoBucket = "chainSet"

func (u UTXOSet) Reindex() {
	db := u.blockchain.db
	bucketName := []byte(utxoBucket)

	err := db.Update(func(tx *bolt.Tx) error {
		err1 := tx.DeleteBucket(bucketName)
		if err1 != nil && err1 != bolt.ErrBucketNotFound {
			log.Panic(err1)
		}
		_, err2 := tx.CreateBucket(bucketName)
		if err2 != nil {
			log.Panic(err2)
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	UTXO := u.blockchain.FindAllUTXO()
	err4 := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		for txID, outs := range UTXO {
			key, err5 := hex.DecodeString(txID)
			if err5 != nil {
				log.Panic(err5)
			}
			err6 := b.Put(key, outs.Serialize())
			if err6 != nil {
				log.Panic(err6)
			}
		}
		return nil
	})
	if err4 != nil {
		log.Panic(err4)
	}

}

func (u UTXOSet) FindUTXOByPublicKeyHash(publickeyHash []byte) []TXOutput {
	var UTXOS []TXOutput
	db := u.blockchain.db
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor() //游标

		for k, v := c.First(); k != nil; k, v = c.Next() {
			outs := DeserializeOutputs(v)
			for _, out := range outs.Outputs {
				if out.CanBeUnlockedWith(publickeyHash) {
					UTXOS = append(UTXOS, out)
				}
			}
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	return UTXOS
}
func (u UTXOSet) update(block *Block) {
	db := u.blockchain.db
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		for _, tx := range block.Transations {
			if tx.IsCoinBase() == false {
				for _, vin := range tx.Vin {
					//针对指定输入找到其对应的引用交易，并更新该引用交易的未花费输出
					updateOuts := TXOutputs{}
					outsbytes := b.Get(vin.TXId) //找到当前交易输入所引用的输出
					outs := DeserializeOutputs(outsbytes)

					for outIdx, out := range outs.Outputs {
						if outIdx != vin.VOutIndex {
							//当前输出依然未被花费，依然需要保留
							updateOuts.Outputs = append(updateOuts.Outputs, out)
						}
					}
					if len(updateOuts.Outputs) == 0 {
						//当前输入对应的引用交易都已经是花费了，直接将其对应的交易未花费输出全部移除掉
						err := b.Delete(vin.TXId)
						if err != nil {
							log.Panic(err)
						}
					} else {
						//将引用对应的交易所对应的未花费输出更新
						err := b.Put(vin.TXId, updateOuts.Serialize())
						if err != nil {
							log.Panic(err)
						}
					}
				}
			}
			//新区块的交易输出都为未花费的输出，直接存储
			newOutputs := TXOutputs{}
			for _, out := range tx.Vout {
				newOutputs.Outputs = append(newOutputs.Outputs, out)
			}
			err := b.Put(tx.ID, newOutputs.Serialize())
			if err != nil {
				log.Panic(err)
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}
