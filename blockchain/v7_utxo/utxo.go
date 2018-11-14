package main

import (
	"encoding/hex"
	"log"
)

// 获取某个地址的可用余额
// 如果input解锁脚本可以解锁对前一个transaction的output，那么前一个transaction的output被花费了
// 如果output锁定脚本没有被花费，它就是可用的
// v6: change address string to pubKeyHash []byte
func (bc *BlockChain) FindUnspentTransactions(pubKeyHash []byte) []Transaction {
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

				//if output.CanBeUnlockedWith(address) {
				if output.IsLockedWithKey(pubKeyHash) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			// 后处理inputs
			// 思考下为什么不应该是coinbase? 难道因为coinbase特点是：只有输出，没有输入！
			if tx.IsCoinbase() == false {
				for _, input := range tx.Vin {
					//if input.CanUnlockOutputWith(address) {
					if input.UsesKey(pubKeyHash) {
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

// v6: change address string to pubKeyHash []byte
// 获取某个地址（解码为PubKeyHash）的所有未花费输出，注意：这里的返回值是TXOutput数组，
// 没有Transaction相关的信息。因此无法从中获取Transaction的上下文。
// 由于这个方法只在获取余额使用，而TXOutput中包括了Value属性。注：获取余额不需要Transaction相关信息，即不管Transaction的位置
func (bc *BlockChain) FindUTXO(pubKeyHash []byte) []TXOutput {
	var UTXOs []TXOutput

	unspentTransactions := bc.FindUnspentTransactions(pubKeyHash)
	for _, tx := range unspentTransactions {
		for _, output := range tx.Vout {
			if output.IsLockedWithKey(pubKeyHash) {
				UTXOs = append(UTXOs, output)
			}
		}
	}

	return UTXOs
}

// v6: FindUnspentTransactions返回的是包含有未花费输出的Transactions，整个Transaction对象，里面有些输出已经花费，但至少有一个输出未花费
// v7: 下面这个方法返回的是含有未花费输出Transaction的ID，以及具体哪几个输出是未花费的
func (bc *BlockChain) FindUnspentTransactionsX() map[string][]TXOutput {
	//var unspentTXs []Transaction
	//未花费
	UTXO := make(map[string][]TXOutput)
	//已花费，key为txid，value为已经花费的输出在trasaction中的index
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

			// 先处理outputs
			// outIdx: TXOutput在outputs数组中的indx
			// output: outIdx对应的TXOutput对象
		Outputs:
			for outIdx, output := range tx.Vout {
				//当前输出在已经花费的输出集合中，跳过，处理下一条
				//如果不在已花费的输出中，说明遍历完了整个事务，都没有找到这个输出，说明这个输出是未花费的，可以加入UTXO集合中
				if spentTXOs[txID] != nil {
					// spentOut(inputSpentOutIdx): 通过(后面先遍历的Tx)inputs，记录的index
					for _, inputSpentOutIdx := range spentTXOs[txID] {
						if inputSpentOutIdx == outIdx {
							// 如果在(txid对应的)spent output index数组里，说明已经被花费了，则它一定不会被加入到unspentTXs中
							continue Outputs
							//注意Outputs的位置！不能放在block.Transactions那里！因为一笔Trasaction会有多个输出,必须处理每个输出
							//只不过有些输出因为被花费过了，所以不需要加入到UTXO集合中
						}
					}
				}

				//if output.CanBeUnlockedWith(address) {
				//if output.IsLockedWithKey(pubKeyHash) {
				//unspentTXs = append(unspentTXs, *tx)
				outputs := UTXO[txID]
				outputs = append(outputs, output)
				UTXO[txID] = outputs
				//}
			}

			// 后处理inputs
			// 思考下为什么不应该是coinbase? 难道因为coinbase特点是：只有输出，没有输入！
			if tx.IsCoinbase() == false {
				for _, input := range tx.Vin {
					//if input.CanUnlockOutputWith(address) {
					//if input.UsesKey(pubKeyHash) {
					prevTxId := hex.EncodeToString(input.Txid)
					prevTxIndex := input.Vout
					// 花费了某个地址对应的transaction，这里的transaction是当前input的前一个transaction，而不是当前transaction
					spentTXOs[prevTxId] = append(spentTXOs[prevTxId], prevTxIndex)
					//}
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return UTXO
}

// v7: 返回值改为map[string][]TXOutput，其中key是TxID，即Transaction的ID。value是Transaction对应的未花费输出集合
func (bc *BlockChain) FindUTXOs() map[string][]TXOutput {
	var UTXOs map[string][]TXOutput

	UTXOs = bc.FindUnspentTransactionsX()

	return UTXOs
}

// return value:
// int: accumulated amount
// map: txId -> output index array 只需要返回超过amount的未花费TXOutput集合即可，不需要返回账户的所有未花费TXOutput
// v6: change address string to pubKeyHash []byte
func (bc *BlockChain) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	accumulated := 0
	unspentOutputs := make(map[string][]int)

	// 获取一个帐号的所有未花费Transactions，但是最后并不需要返回所有未花费输出的余额总和
	unspentTransactions := bc.FindUnspentTransactions(pubKeyHash)

Work:
	for _, tx := range unspentTransactions {
		txId := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Vout {
			// TXOutput属于address，且收集到的余额还不够，继续找下一个TXOutput
			//if out.CanBeUnlockedWith(address) && accumulated < amount {
			if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
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

// v6: bc *BlockChain
// v7: UTXOSet
func NewUTXOTransaction(from, to string, amount int, UTXOSet *UTXOSet) *Transaction {
	var txinputs []TXInput
	var txoutputs []TXOutput

	// 从钱包集合中找出from这个地址对应的钱包对象，并根据钱包里的公钥，计算出公钥哈希
	// 获取区块链的已花费输出，传递的是pubKeyHash，而不是原先的from
	// from表示的是bit coin address。实际上要获取pubKeyHash，可以直接对address进行解码也是可以的
	wallets, error := NewWallets()
	if error != nil {
		log.Panic(error)
	}
	wallet := wallets.GetWallet(from)
	pubKeyHash := HashPubKey(wallet.PublicKey)

	// 转账发起方必须要有足够的余额，从区块链中找出未花费的输出TXOutput
	//accumulated, validOutput := bc.FindSpendableOutputs(from, amount)
	//accumulated, validOutput := bc.FindSpendableOutputs(pubKeyHash, amount)
	accumulated, validOutput := UTXOSet.FindSpendableOutputs(pubKeyHash, amount)
	if accumulated < amount {
		log.Panic("ERROR: 余额不够！")
	}

	// 未花费的输出，作为Transaction的inputs
	// txid是未花费输出的transaction id, outputs是index数组
	for txid, outputidxs := range validOutput {
		txID, _ := hex.DecodeString(txid)

		for _, outIdx := range outputidxs {
			//txinput := TXInput{txID, outIdx, from}
			txinput := TXInput{txID, outIdx, nil, wallet.PublicKey}
			txinputs = append(txinputs, txinput)
		}
	}

	// 转账接收方的输出
	//toTXOutput := TXOutput{amount, to}
	//txoutputs = append(txoutputs, toTXOutput)
	toTXOutput := NewTXOutput(amount, to)
	txoutputs = append(txoutputs, *toTXOutput)

	// 如果转账发起方的余额大于amount，零钱找回
	if accumulated > amount {
		//fromTXOutput := TXOutput{accumulated - amount, from}
		//txoutputs = append(txoutputs, fromTXOutput)
		fromTXOutput := NewTXOutput(accumulated-amount, from)
		txoutputs = append(txoutputs, *fromTXOutput)
	}

	tx := Transaction{nil, txinputs, txoutputs}
	//tx.SetID()
	tx.ID = tx.Hash()

	//签名
	//bc.SignTransaction(&tx, wallet.PrivateKey)
	UTXOSet.Blockchain.SignTransaction(&tx, wallet.PrivateKey)

	return &tx
}
