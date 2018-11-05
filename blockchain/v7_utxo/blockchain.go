package main

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

type BlockChain struct {
	//blocks []*Block
	db  *bolt.DB
	tip []byte
}

type BlockChainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (bc *BlockChain) Iterator() *BlockChainIterator {
	bci := &BlockChainIterator{bc.tip, bc.db}
	return bci
}

func (i *BlockChainIterator) Next() *Block {
	var block *Block
	error := i.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		encodedBlock := bucket.Get(i.currentHash)
		block = Deserialize(encodedBlock)
		return nil
	})

	if error != nil {
		log.Panic(error)
	}

	i.currentHash = block.PrevBlockHash

	return block
}

func (blockChain *BlockChain) AddBlock(transactions []*Transaction) *Block {
	// prevBlock := blockChain.blocks[len(blockChain.blocks)-1]
	// newBlock := NewBlock(data, prevBlock.Hash)
	// blockChain.blocks = append(blockChain.blocks, newBlock)

	var lastHash []byte

	// 在创建区块之前，要验证每一笔交易
	for _, tx := range transactions {
		// 只要有一笔交易验证失败，就不能创建区块
		if blockChain.VerifyTransaction(tx) != true {
			log.Panic("ERROR: Invalid transaction!")
		}
	}

	err := blockChain.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))
		return nil
	})

	//newBlock := NewBlock(data, lastHash)
	newBlock := NewBlock(transactions, lastHash)

	err = blockChain.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		b.Put(newBlock.Hash, newBlock.Serialize())
		b.Put([]byte("l"), newBlock.Hash)
		blockChain.tip = newBlock.Hash
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	return newBlock
}

// func NewGenesisBlock() *Block {
// 	return NewBlock("Genesis Block", []byte{})
// }

func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

// func NewBlockChain() *BlockChain {
// 	return &BlockChain{[]*Block{NewGenesisBlock()}}
// }

// 使用该方法时，必须先调用过CreateBlockChain（而且只能调用一次）
func NewBlockChain() *BlockChain {
	if dbExists() == false {
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}

	var tip []byte

	db, err := bolt.Open(dbFile, 0600, nil)
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		tip = bucket.Get([]byte("l"))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	//fmt.Println("new blockchain tip:", tip)
	return &BlockChain{db, tip}
}

func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

// 创建一个新的区块链数据库，address用来接收挖出创世块的奖励
func CreateBlockchain(address string) *BlockChain {
	// 只能调用一次CreateBlockChain方法
	if dbExists() {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
		genesis := NewGenesisBlock(cbtx)
		bucket, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic(err)
		}
		genesisHash := genesis.Hash
		//fmt.Println("genesis hash:", genesisHash)
		err = bucket.Put(genesisHash, genesis.Serialize())
		if err != nil {
			log.Panic(err)
		}
		err = bucket.Put([]byte("l"), genesisHash)
		if err != nil {
			log.Panic(err)
		}
		tip = genesisHash

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	blockchain := BlockChain{db, tip}
	return &blockchain
}

func (bc *BlockChain) FindTransaction(ID []byte) (Transaction, error) {
	bci := bc.Iterator()
	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
	return Transaction{}, errors.New("")
}

func (bc *BlockChain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)

	// Find current tx's all vin's reference Transaction
	for _, vin := range tx.Vin {
		prevTX, _ := bc.FindTransaction(vin.Txid)
		prevTXs[hex.EncodeToString(vin.Txid)] = prevTX
	}

	tx.Sign(privKey, prevTXs)
}

func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {
	prevTXs := make(map[string]Transaction)

	// Find current tx's all vin's reference Transaction
	for _, vin := range tx.Vin {
		prevTX, _ := bc.FindTransaction(vin.Txid)
		prevTXs[hex.EncodeToString(vin.Txid)] = prevTX
	}

	return tx.Verify(prevTXs)
}
