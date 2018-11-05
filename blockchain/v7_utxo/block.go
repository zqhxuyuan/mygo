package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	Timestamp int64
	//Data          []byte
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int // add nonce field
}

func (block *Block) HashTransactions() []byte {
	var txhashes [][]byte
	var txhash [32]byte

	for _, tx := range block.Transactions {
		txhashes = append(txhashes, tx.ID)
	}

	txhash = sha256.Sum256(bytes.Join(txhashes, []byte{}))
	return txhash[:]
}

// func (block *Block) SetHash() {
// 	timestamp := []byte(strconv.FormatInt(block.Timestamp, 10))
// 	headers := bytes.Join([][]byte{timestamp, block.Data, block.PrevBlockHash}, []byte{})
// 	hash := sha256.Sum256(headers)
// 	block.Hash = hash[:]
// }

func (b *Block) Serialize() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

func Deserialize(buf []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(buf))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}

func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{
		time.Now().Unix(),
		//[]byte(data),
		transactions,
		prevBlockHash,
		[]byte{},
		0,
	}
	//block.SetHash()
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()
	block.Nonce = nonce
	block.Hash = hash[:]
	return block
}
