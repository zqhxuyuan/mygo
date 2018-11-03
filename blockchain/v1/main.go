package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
}

func (block *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(block.Timestamp, 10))
	headers := bytes.Join([][]byte{timestamp, block.Data, block.PrevBlockHash}, []byte{})
	hash := sha256.Sum256(headers)
	block.Hash = hash[:]
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		time.Now().Unix(),
		[]byte(data),
		prevBlockHash,
		[]byte{},
	}
	block.SetHash()
	return block
}

type BlockChain struct {
	blocks []*Block
}

func (blockChain *BlockChain) AddBlock(data string) {
	prevBlock := blockChain.blocks[len(blockChain.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	blockChain.blocks = append(blockChain.blocks, newBlock)
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func NewBlockChain() *BlockChain {
	return &BlockChain{[]*Block{NewGenesisBlock()}}
}

func main() {
	blockChain := NewBlockChain()
	blockChain.AddBlock("Block1: alice 1 BTC to Bob.")
	blockChain.AddBlock("Block2: bob 1 BTC to Cal.")

	for _, block := range blockChain.blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()
	}
}
