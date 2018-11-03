package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

const subsidy = 10 //是挖出新块的奖励金

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

type TXInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
}

type TXOutput struct {
	Value        int
	ScriptPubKey string
}

func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vout) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
	// 文章中的代码是错误的，tx.Vin表示输入，coinbase没有输入。输出只有一个，所以应该是len(tx.Vout)=1
	//return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encode := gob.NewEncoder(&encoded)
	error := encode.Encode(tx)
	if error != nil {
		log.Panic(error)
	}

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func (out *TXOutput) CanBeUnlockedWith(unlockingdata string) bool {
	return out.ScriptPubKey == unlockingdata
}

func (input *TXInput) CanUnlockOutputWith(unlockingdata string) bool {
	return input.ScriptSig == unlockingdata
}

// NewCoinbaseTX 构建coinbase交易，没有输入，只有一个输出
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	txinput := TXInput{[]byte{}, -1, data}
	txoutput := TXOutput{subsidy, to}

	tx := Transaction{[]byte{}, []TXInput{txinput}, []TXOutput{txoutput}}
	tx.SetID()
	return &tx
}
