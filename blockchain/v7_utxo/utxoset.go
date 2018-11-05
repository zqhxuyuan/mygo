package main

import (
	"encoding/hex"
	"log"

	"github.com/boltdb/bolt"
)

const utxoBucket = "chainstate"

type UTXOSet struct {
	Blockchain *BlockChain
}

// UTXO保存在DB中，需要用到未花费的输出，直接查询数据库
// 注意：数据库中保存的key是transaction的ID，与用户什么无关，
// 所以要查询某个地址，还是需要遍历数据库。但是好歹比遍历整个区块链要快
func (u UTXOSet) Reindex() {
	db := u.Blockchain.db
	bucketName := []byte(utxoBucket)

	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(bucketName)
		_, err = tx.CreateBucket(bucketName)
		return err
	})

	UTXO := u.Blockchain.FindUTXOs()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)

		for txID, outs := range UTXO {
			key, err := hex.DecodeString(txID)

			// outs是一个TXOutput[]数组，可以在FindUTXOs中直接返回TXOutputs对象，
			// 当然了，也可以像这里，只有在涉及到序列化和反序列化时，才需要转为TXOutputs
			outputs := TXOutputs{outs}
			err = b.Put(key, outputs.Serialize())

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

func (u UTXOSet) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	accumulated := 0
	unspentOutputs := make(map[string][]int)

	// 获取一个帐号的所有未花费Transactions，但是最后并不需要返回所有未花费输出的余额总和
	//unspentTransactions := bc.FindUnspentTransactions(pubKeyHash)

	db := u.Blockchain.db

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			txId := hex.EncodeToString(k)
			outs := DeserializeOutputs(v)

			for outIdx, out := range outs.TXOutputs {
				if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
					accumulated += out.Value
					unspentOutputs[txId] = append(unspentOutputs[txId], outIdx)
				}
			}
		}
		return nil
	})

	return accumulated, unspentOutputs
}

func (u UTXOSet) FindUTXO(pubKeyHash []byte) []TXOutput {
	var UTXOs []TXOutput
	db := u.Blockchain.db

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			outs := DeserializeOutputs(v)

			for _, out := range outs.TXOutputs {
				if out.IsLockedWithKey(pubKeyHash) {
					UTXOs = append(UTXOs, out)
				}
			}
		}

		return nil
	})

	return UTXOs
}

func (u UTXOSet) Update(block *Block) {
	db := u.Blockchain.db

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(utxoBucket))

		for _, tx := range block.Transactions {
			if tx.IsCoinbase() == false {
				// 新的Transaction，输入引用了未花费的输出
				// 根据txId从UTXOSet中获取出所有未花费的输出集合
				// 如果输入引用的输出索引与未花费输出集合中的某个索引相等
				// 那么对应这个索引表示被花费了！
				for _, vin := range tx.Vin {
					updatedOuts := TXOutputs{}
					outsBytes := b.Get(vin.Txid)
					outs := DeserializeOutputs(outsBytes)

					for outIdx, out := range outs.TXOutputs {
						if outIdx != vin.Vout {
							updatedOuts.TXOutputs = append(updatedOuts.TXOutputs, out)
						}
					}

					if len(updatedOuts.TXOutputs) == 0 {
						b.Delete(vin.Txid)
					} else {
						b.Put(vin.Txid, updatedOuts.Serialize())
					}

				}
			}

			newOutputs := TXOutputs{}
			for _, out := range tx.Vout {
				newOutputs.TXOutputs = append(newOutputs.TXOutputs, out)
			}

			b.Put(tx.ID, newOutputs.Serialize())
		}
		return nil
	})
}
