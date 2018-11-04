package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
)

const targetBits = 24
const dbFile = "blockchain.db"
const blocksBucket = "blocks"
const genesisCoinbaseData = "Genesis..."

var (
	maxNonce = math.MaxInt64
)

func main() {
	//blockchain := NewBlockChain()
	//defer blockchain.db.Close()

	//cli := CLI{blockchain}
	cli := CLI{}
	cli.Run()
}

type CLI struct {
	//blockchain *BlockChain
}

// 原先的AddBlock由send方法代替。另外，增加了getbalance方法
func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  getbalance -address ADDRESS - Get balance of ADDRESS")
	fmt.Println("  createblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("  printchain - Print all the blocks of the blockchain")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT - Send AMOUNT of coins from FROM address to TO")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) Run() {
	cli.validateArgs()

	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	//addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	//addBlockData := addBlockCmd.String("data", "", "BlockData")

	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	createBlockchainAddress := createBlockchainCmd.String("address", "", "the address to send genesis block reward")

	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	sendFrom := sendCmd.String("from", "", "transfer sender address")
	sendTo := sendCmd.String("to", "", "transfer receiver address")
	sendAmount := sendCmd.Int("amount", 0, "transfer how much money")

	getBalanceCmd := flag.NewFlagSet("getBalance", flag.ExitOnError)
	getBalanceAddress := getBalanceCmd.String("address", "", "address")

	switch os.Args[1] {
	// case "addblock":
	// 	error := addBlockCmd.Parse(os.Args[2:])
	// 	if error != nil {
	// 		log.Panic(error)
	// 	}
	case "createblockchain":
		error := createBlockchainCmd.Parse(os.Args[2:])
		if error != nil {
			log.Panic(error)
		}
	case "printchain":
		error := printChainCmd.Parse(os.Args[2:])
		if error != nil {
			log.Panic(error)
		}
	case "send":
		error := sendCmd.Parse(os.Args[2:])
		if error != nil {
			log.Panic(error)
		}
	case "getBalance":
		error := getBalanceCmd.Parse(os.Args[2:])
		if error != nil {
			log.Panic(error)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	// if addBlockCmd.Parsed() {
	// 	if *addBlockData == "" {
	// 		addBlockCmd.Usage()
	// 		os.Exit(1)
	// 	}
	// 	cli.addBlock(*addBlockData)
	// }

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}
		cli.createBlockchain(*createBlockchainAddress)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if sendCmd.Parsed() {
		if *sendFrom != "" && *sendTo != "" && *sendAmount > 0 {
			cli.send(*sendFrom, *sendTo, *sendAmount)
		} else {
			sendCmd.Usage()
			os.Exit(1)
		}
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress)
	}
}

func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}
