package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

const subsidy = 10 //是挖出新块的奖励金

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

type TXInput struct {
	Txid []byte
	Vout int
	//ScriptSig string
	Signature []byte // 签名数据
	PubKey    []byte // 公钥
}

type TXOutput struct {
	Value int
	//ScriptPubKey string
	PubKeyHash []byte // 公钥的hash，对地址进行Base58编码
}

func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vout) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
	// 文章中的代码是错误的，tx.Vin表示输入，coinbase没有输入。输出只有一个，所以应该是len(tx.Vout)=1
	//return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

// func (tx *Transaction) SetID() {
// 	var encoded bytes.Buffer
// 	var hash [32]byte
// 	encode := gob.NewEncoder(&encoded)
// 	error := encode.Encode(tx)
// 	if error != nil {
// 		log.Panic(error)
// 	}
// 	hash = sha256.Sum256(encoded.Bytes())
// 	tx.ID = hash[:]
// }

func (tx Transaction) Serialize() []byte {
	var encoded bytes.Buffer
	encoder := gob.NewEncoder(&encoded)
	encoder.Encode(tx)
	return encoded.Bytes()
}

func (tx *Transaction) Hash() []byte {
	var hash [32]byte
	//encoded := tx.Serialize()
	//hash = sha256.Sum256(encoded)

	txCopy := *tx
	txCopy.ID = []byte{}
	hash = sha256.Sum256(txCopy.Serialize())
	return hash[:]
}

func (input *TXInput) UsesKey(pubKeyHash []byte) bool {
	// 对公钥计算hash
	lockingHash := HashPubKey(input.PubKey)

	// 比较计算结果与给定的Public Key Hash是否相同
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}

// 给定一个地址，解码出其中的公钥哈希，这个是输出的字段，表示交易的输出对某个地址进行锁定
// 注意：给定公钥，可以计算出地址。但是给定地址，不能计算出公钥，只能解码出公钥的哈希
func (output *TXOutput) Lock(address []byte) {
	// 根据（Bitcoin）地址，计算出公钥的hash
	pubKeyHash := Base58Decode(address)
	// 去掉最开始的version，和最末尾的checksum，中间部分就是public key hash
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	// 设置TXOutput的公钥哈希字段
	output.PubKeyHash = pubKeyHash
}

// 如果给定的public key hash与输出的公钥哈希一致，说明这个public key hash可以解锁交易的输出
// 参数中的pubKeyHash一般是输入的PubKey经过哈希计算的结果
func (output *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(output.PubKeyHash, pubKeyHash) == 0
}

func NewTXOutput(value int, address string) *TXOutput {
	txo := TXOutput{value, nil}
	txo.Lock([]byte(address))
	return &txo
}

// func (out *TXOutput) CanBeUnlockedWith(address string) bool {
// 	//return out.ScriptPubKey == address
// 	pubKeyHash := Base58Decode([]byte(address))
// 	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
// 	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
// }

// func (input *TXInput) CanUnlockOutputWith(address string) bool {
// 	//return input.ScriptSig == unlockingdata
// 	pubKey := input.PubKey
// 	encodeAddress := Base58Encode(pubKey)
// 	return bytes.Compare(encodeAddress, []byte(address)) == 0
// }

// NewCoinbaseTX 构建coinbase交易，没有输入，只有一个输出
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	//txinput := TXInput{[]byte{}, -1, data}
	//txoutput := TXOutput{subsidy, to}
	txinput := TXInput{[]byte{}, -1, nil, []byte(data)}
	txoutput := NewTXOutput(subsidy, to)

	tx := Transaction{[]byte{}, []TXInput{txinput}, []TXOutput{*txoutput}}
	//tx.SetID()
	tx.ID = tx.Hash()
	return &tx
}
