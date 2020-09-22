package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct {
	bc *Blockchain
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		fmt.Println("参数小于1")
		os.Exit(1)
	}
	//	fmt.Println(os.Args)
}

//命令行   ./main addblock
func (cli *CLI) addBlock() {
	fmt.Println("增加区块")
	cli.bc.MineBlock([]*Transation{})
}

//命令行   ./main printchain
func (cli *CLI) printChain() {
	fmt.Println("打印区块")
	cli.bc.printBlockchain()
}

//命令行   ./main getbalance --address=caspar   查询地址为caspar的余额
//查询余额
func (cli *CLI) getBalance(address string) {
	balance := 0
	//由命令行address获取公钥hash：publicKeyHash
	publicKeyHash := GetPublicKeyHashFromAddress(address)
	fmt.Printf("address %s,公钥哈希:%x\n", address, publicKeyHash)
	set := UTXOSet{cli.bc}
	UTXOs := set.FindUTXOByPublicKeyHash(publicKeyHash)
	//UTXOs:=cli.bc.FindUTXO(publicKeyHash)
	for _, out := range UTXOs {
		fmt.Printf("%x\n", out.PublicKeyHash)
		fmt.Printf("%d\n", out.Value)
		balance = balance + out.Value
	}
	fmt.Printf("\n balance of '%s' :%d \n", address, balance)
}

//转账
//命令行  ./main transfer -from caspar -to leo -amount 20
func (cli *CLI) transfer(from, to string, amount int) {
	tx := NewUTXOTransation(from, to, amount, cli.bc) //构建交易
	//进行挖矿后加入区块链
	cli.bc.MineBlock([]*Transation{tx})
	fmt.Printf("\n transfer from :'%s' to '%s' %d success! \n", from, to, amount)
}

//创建钱包
func (cli *CLI) createWallet() {
	wallets, _ := NewWallets()
	address := wallets.CreateWallet()
	wallets.SaveToFile()
	fmt.Printf("your address:%s\n", address)
}
func (cli *CLI) listAddress() {
	wallets, err := NewWallets()
	if err != nil {
		log.Panic(err)
	}
	addresses := wallets.getAddress()
	for _, address := range addresses {
		fmt.Println(address)
	}
}

//打印提示
func (cli *CLI) printUsage() {
	fmt.Println("usage:")
	fmt.Println("addblock:增加区块")
	fmt.Println("printchain:打印区块")
	fmt.Println("getbalance:获取余额")
	fmt.Println("createWallet:创建地址")
	fmt.Println("listAddress:获取地址列表")
}

//"入口" 是 Run 函数
func (cli *CLI) Run() {
	cli.validateArgs()
	/***添加区块**/
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	/***打印区块*/
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	//获取余额
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	getBalanceAddress := getBalanceCmd.String("address", "", "the address to get balance of")
	/**转账测试*/
	transferCmd := flag.NewFlagSet("transfer", flag.ExitOnError)
	transferFrom := transferCmd.String("from", "", "付款地址/人")
	transferTo := transferCmd.String("to", "", "收款地址/人")
	transferAmount := transferCmd.Int("amount", 0, "转账金额")
	/***钱包测试**/
	createWalletCmd := flag.NewFlagSet("createWallet", flag.ExitOnError)
	listAddressCmd := flag.NewFlagSet("listAddress", flag.ExitOnError)

	switch os.Args[1] {

	case "createWallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "listAddress":
		err := listAddressCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "transfer":
		err := transferCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	case "addblock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}
	if addBlockCmd.Parsed() {
		cli.addBlock()
	}
	if printChainCmd.Parsed() {
		cli.printChain()
	}
	if getBalanceCmd.Parsed() {
		//	fmt.Printf("getBalanceAddress:%s\n",*getBalanceAddress)
		if *getBalanceAddress == "" {
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress)
	}
	if transferCmd.Parsed() {
		//	fmt.Println("cli transfer transferCmd")
		if *transferFrom == "" || *transferTo == "" || *transferAmount < 0 {
			os.Exit(1)
		}
		//	fmt.Println("cli transfer start")
		cli.transfer(*transferFrom, *transferTo, *transferAmount)
	}

	if createWalletCmd.Parsed() {
		cli.createWallet()
	}
	if listAddressCmd.Parsed() {
		cli.listAddress()
	}

}
