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

func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	for _, in := range tx.Vin {
		inputs = append(inputs, TXInput{in.Txid, in.Vout, nil, nil})
	}

	for _, out := range tx.Vout {
		outputs = append(outputs, TXOutput{out.Value, out.PubKeyHash})
	}

	txCopy := Transaction{tx.ID, inputs, outputs}
	return txCopy
}

func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsCoinbase() {
		return
	}

	txCopy := tx.TrimmedCopy()

	// 循环副本的每个输入
	for inID, vin := range txCopy.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]

		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash

		txCopy.ID = txCopy.Hash()

		txCopy.Vin[inID].PubKey = nil

		r, s, error := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
		if error != nil {
			log.Panic(error)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		// 对于原始tx对象里TXInput的PubKey是不变的，也不为nil!
		tx.Vin[inID].Signature = signature
	}
}

func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	// 循环实际Transaction的每个输入，注意：不是副本！
	for inID, vin := range tx.Vin {
		// 下面的逻辑和签名的一样
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Vin[inID].PubKey = nil

		// 解析输入中的签名：签名是一对数字
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		// 解析输入中的PubKey：PubKey是一对坐标
		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keyLen / 2)])
		y.SetBytes(vin.PubKey[(keyLen / 2):])

		// 从坐标可以计算出PubKey
		rawPubKey := ecdsa.PublicKey{curve, &x, &y}

		// 验证签名是否正确：
		// txCopy.ID是由prevTx计算出来的签名，而r,s是在vin.Signature中的签名解码出来的
		// 第一个参数是公钥。如果解码出来的r,s的签名正好是prevTx计算出来的签名，则验证通过
		// 对比签名方法：r, s = Sign(privateKey, txCopy.ID)
		if ecdsa.Verify(&rawPubKey, txCopy.ID, &r, &s) == false {
			return false
		}
	}
	return true
}
