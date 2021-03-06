package blockDemo

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type CLI struct {
	bc *Blockchain
}

func (cli *CLI) addBlock(data string) {
	// cli.bc.AddBlock(data)
	fmt.Println("Success!")
}

func (cli *CLI) printChain() {
	bci := cli.bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		// fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

func (cli *CLI) getBalance(address string) {
	// bc := NewBlockchain(address)
	// defer bc.db.Close()

	balance := 0
	UTXOs := cli.bc.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}

func (cli *CLI) validateArgs() {}

func (cli *CLI) printUsage() {}

func (cli *CLI) Run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "Block data")
	getBalanceData := getBalanceCmd.String("address", "", "address data")

	switch os.Args[1] {
	case "addblock":
		if err := addBlockCmd.Parse(os.Args[2:]);err != nil {
			panic("addBlockCmd parse fail")
		}
	case "printchain":
		if err := printChainCmd.Parse(os.Args[2:]); err != nil {
			panic("printChainCmd parse fail")
		}
	case "getbalance":
		if err := getBalanceCmd.Parse(os.Args[2:]); err != nil {
			panic("getBlanceCmd parse fail")
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}
	// 解析获取余额指令
	if getBalanceCmd.Parsed() {
		if *getBalanceData == "" {
			getBalanceCmd.Usage()
			os.Exit(1)
		}
		cli.getBalance(*getBalanceData)
	}
	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}

func CliRun(bc *Blockchain) {
	cli := CLI{bc}
	cli.Run()
}
