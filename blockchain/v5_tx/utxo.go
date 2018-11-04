package main

import (
	"encoding/hex"
	"log"
)

// 获取某个地址的可用余额
// 如果input解锁脚本可以解锁对前一个transaction的output，那么前一个transaction的output被花费了
// 如果output锁定脚本没有被花费，它就是可用的
func (bc *BlockChain) FindUnspentTransactions(address string) []Transaction {
	var unspentTXs []Transaction
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.Next()

	Outputs:
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

			// 先处理outputs
			// outIdx: TXOutput在outputs数组中的indx
			// output: outIdx对应的TXOutput对象
			for outIdx, output := range tx.Vout {
				if spentTXOs[txID] != nil {
					// spentOut(inputSpentOutIdx): 通过(后面先遍历的Tx)inputs，记录的index
					for _, inputSpentOutIdx := range spentTXOs[txID] {
						if inputSpentOutIdx == outIdx {
							// 如果在(txid对应的)spent output index数组里，说明已经被花费了，则它一定不会被加入到unspentTXs中
							continue Outputs
						}
					}
				}

				if output.CanBeUnlockedWith(address) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			// 后处理inputs
			// 思考下为什么不应该是coinbase? 难道因为coinbase特点是：只有输出，没有输入！
			if tx.IsCoinbase() == false {
				for _, input := range tx.Vin {
					if input.CanUnlockOutputWith(address) {
						prevTxId := hex.EncodeToString(input.Txid)
						prevTxIndex := input.Vout
						// 花费了某个地址对应的transaction，这里的transaction是当前input的前一个transaction，而不是当前transaction
						spentTXOs[prevTxId] = append(spentTXOs[prevTxId], prevTxIndex)
					}
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return unspentTXs
}

func (bc *BlockChain) FindUTXO(address string) []TXOutput {
	var UTXOs []TXOutput
	unspentTransactions := bc.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions {
		for _, output := range tx.Vout {
			if output.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, output)
			}
		}
	}

	return UTXOs
}

// return value:
// int: accumulated amount
// map: txId -> output index array 只需要返回超过amount的未花费TXOutput集合即可，不需要返回账户的所有未花费TXOutput
func (bc *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	accumulated := 0
	unspentOutputs := make(map[string][]int)

	// 获取一个帐号的所有未花费Transactions，但是最后并不需要返回所有未花费输出的余额总和
	unspentTransactions := bc.FindUnspentTransactions(address)

Work:
	for _, tx := range unspentTransactions {
		txId := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Vout {
			// TXOutput属于address，且收集到的余额还不够，继续找下一个TXOutput
			if out.CanBeUnlockedWith(address) && accumulated < amount {
				accumulated += out.Value
				// 需要记录下来哪些transaction，以及对应的index，需要被后续花费（即转账）
				unspentOutputs[txId] = append(unspentOutputs[txId], outIdx)

				// 收集够了，立即退出
				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}

func NewUTXOTransaction(from, to string, amount int, bc *BlockChain) *Transaction {
	var txinputs []TXInput
	var txoutputs []TXOutput

	// 转账发起方必须要有足够的余额，从区块链中找出未花费的输出TXOutput
	accumulated, validOutput := bc.FindSpendableOutputs(from, amount)
	if accumulated < amount {
		log.Panic("ERROR: 余额不够！")
	}

	// 未花费的输出，作为Transaction的inputs
	// txid是未花费输出的transaction id, outputs是index数组
	for txid, outputidxs := range validOutput {
		txID, _ := hex.DecodeString(txid)

		for _, outIdx := range outputidxs {
			txinput := TXInput{txID, outIdx, from}
			txinputs = append(txinputs, txinput)
		}
	}

	// 转账接收方的输出
	toTXOutput := TXOutput{amount, to}
	txoutputs = append(txoutputs, toTXOutput)

	// 如果转账发起方的余额大于amount，零钱找回
	if accumulated > amount {
		fromTXOutput := TXOutput{accumulated - amount, from}
		txoutputs = append(txoutputs, fromTXOutput)
	}

	tx := &Transaction{nil, txinputs, txoutputs}
	tx.SetID()
	return tx
}
