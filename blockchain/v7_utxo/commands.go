package main

import (
	"fmt"
	"strconv"
)

func (cli *CLI) createBlockchain(address string) {
	blockchain := CreateBlockchain(address)
	defer blockchain.db.Close()
	fmt.Println("create block chain done!")

	UTXOSet := UTXOSet{blockchain}
	UTXOSet.Reindex()
}

// 现在不能直接addBlock了
// func (cli *CLI) addBlock(data string) {
// 	cli.blockchain.AddBlock(data)
// 	fmt.Println("add block success!")
// }

// 通过send方式来创建区块
func (cli *CLI) send(from, to string, amount int) {
	blockchain := NewBlockChain()
	defer blockchain.db.Close()

	UTXOSet := UTXOSet{blockchain}

	//tx := NewUTXOTransaction(from, to, amount, blockchain)
	tx := NewUTXOTransaction(from, to, amount, &UTXOSet)

	// 假设发送交易的人，就是挖出块的人，所以coinbase奖励给from
	coinbaseTX := NewCoinbaseTX(from, "")
	txs := []*Transaction{tx, coinbaseTX}

	// change AddBlock method name to MineBlock
	newBlock := blockchain.AddBlock(txs)
	fmt.Println("Mined Block done!")

	UTXOSet.Update(newBlock)
	fmt.Println("Update UTXO success!")
}

func (cli *CLI) printChain() {
	//blockchainiterator := cli.blockchain.Iterator()
	blockchain := NewBlockChain()
	defer blockchain.db.Close()

	blockchainiterator := blockchain.Iterator()
	for {
		block := blockchainiterator.Next()
		fmt.Printf("Prev hash: %x\n", block.PrevBlockHash)
		//fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

// 获取某个地址的总余额
func (cli *CLI) getBalance(address string) {
	blockchain := NewBlockChain()
	defer blockchain.db.Close()

	UTXOSet := UTXOSet{blockchain}

	pubKeyHash := Address2PubKeyHash(address)

	balance := 0
	// BlockChain有多个Find相关的方法
	//utxos := blockchain.FindUTXO(address)
	//utxos := blockchain.FindUTXO(pubKeyHash)
	utxos := UTXOSet.FindUTXO(pubKeyHash)
	for _, output := range utxos {
		balance += output.Value
	}
	fmt.Printf("Balance of '%s': %d\n", address, balance)
}

func (cli *CLI) createWallet() {
	wallets, _ := NewWallets()
	address := wallets.CreateWallet()
	wallets.SaveToFile()
	fmt.Printf("Your new address: %s\n", address)
}

func (cli *CLI) listAddresses() {
	wallets, _ := NewWallets()
	addresses := wallets.GetAddresses()
	for _, address := range addresses {
		fmt.Println(address)
	}
}
