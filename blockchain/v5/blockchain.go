package main

import (
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
	i.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		encodedBlock := bucket.Get(i.currentHash)
		block = Deserialize(encodedBlock)
		return nil
	})

	i.currentHash = block.PrevBlockHash

	return block
}

func (blockChain *BlockChain) AddBlock(transactions []*Transaction) {
	// prevBlock := blockChain.blocks[len(blockChain.blocks)-1]
	// newBlock := NewBlock(data, prevBlock.Hash)
	// blockChain.blocks = append(blockChain.blocks, newBlock)

	var lastHash []byte

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

func NewBlockChain(address string) *BlockChain {
	var tip []byte

	db, err := bolt.Open(dbFile, 0600, nil)
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		if bucket == nil {
			fmt.Println("No existing blockchain found. Creating a new one...")
			//genesis := NewGenesisBlock()
			coinbaseTx := NewCoinbaseTX(address, genesisCoinbaseData)
			genesis := NewGenesisBlock(coinbaseTx)
			bucket, err = tx.CreateBucket([]byte(blocksBucket))
			err = bucket.Put(genesis.Hash, genesis.Serialize())
			err = bucket.Put([]byte("l"), genesis.Hash)
			tip = genesis.Hash
		} else {
			tip = bucket.Get([]byte("l"))
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
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
		err = bucket.Put(genesis.Hash, genesis.Serialize())
		err = bucket.Put([]byte("1"), genesis.Hash)
		tip = genesis.Hash
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	blockchain := BlockChain{db, tip}
	return &blockchain
}
