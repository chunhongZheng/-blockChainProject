package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct {
	bc *BlockChain
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
	fmt.Println("addBlock:增加区块")
	fmt.Println("printChain:打印区块")
	fmt.Println("getBalance:获取余额")
	fmt.Println("createWallet:创建地址")
	fmt.Println("listAddress:获取地址列表")
	fmt.Println("getBlockHeight:获取区块高度")
	fmt.Println("startNode:启动节点")
}

//获取区块高度
func (cli *CLI) getBlockHeight() {
	fmt.Printf("区块高度为:%d\n", cli.bc.GetBestHeight())
}

//启动节点服务
func (cli *CLI) startNode(nodeID string, minerAddress string) {
	fmt.Printf("Starting node %s\n", nodeID)
	if len(minerAddress) > 0 {
		if IsVaildBitcoinAddress(minerAddress) {
			fmt.Printf("%s miner is on", minerAddress)
			//	 StartServer(nodeID,minerAddress,cli.bc)
		} else {
			log.Panic("invalid miner Address")
		}
	} else {
		// 	log.Panic("miner Address must not be empty")
	}
	StartServer(nodeID, minerAddress, cli.bc)
}

//"入口" 是 Run 函数
func (cli *CLI) Run() {
	cli.validateArgs()

	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		fmt.Printf("NODE_ID is not set")
		os.Exit(1)
	}

	/***添加区块**/
	addBlockCmd := flag.NewFlagSet("addBlock", flag.ExitOnError)
	/***打印区块*/
	printChainCmd := flag.NewFlagSet("printChain", flag.ExitOnError)
	//获取余额
	getBalanceCmd := flag.NewFlagSet("getBalance", flag.ExitOnError)
	getBalanceAddress := getBalanceCmd.String("address", "", "the address to get balance of")
	/**转账测试*/
	transferCmd := flag.NewFlagSet("transfer", flag.ExitOnError)
	transferFrom := transferCmd.String("from", "", "付款地址/人")
	transferTo := transferCmd.String("to", "", "收款地址/人")
	transferAmount := transferCmd.Int("amount", 0, "转账金额")
	/***钱包测试**/
	createWalletCmd := flag.NewFlagSet("createWallet", flag.ExitOnError)
	listAddressCmd := flag.NewFlagSet("listAddress", flag.ExitOnError)
	//区块高度
	getBlockHeightCmd := flag.NewFlagSet("getBlockHeight", flag.ExitOnError)
	//启动服务
	startNodeCmd := flag.NewFlagSet("startNode", flag.ExitOnError)
	startNodeMinner := startNodeCmd.String("minerAddress", "", "矿工地址")
	switch os.Args[1] {
	case "startNode":
		err := startNodeCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "getBlockHeight":
		err := getBlockHeightCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
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

	case "getBalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}

	case "addBlock":
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printChain":
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
	if getBlockHeightCmd.Parsed() {
		cli.getBlockHeight()
	}
	if startNodeCmd.Parsed() {
		nodeID := os.Getenv("NODE_ID")
		if nodeID == "" {
			startNodeCmd.Usage()
			os.Exit(1)
		}
		cli.startNode(nodeID, *startNodeMinner)
	}
}
