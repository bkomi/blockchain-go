package main

import (
	"fmt"

	"github.com/bkomi/blockchain-go/blockchain"
)

func main() {

	chain := blockchain.InitBlockChain()

	chain.AddBlock("First Block")
	chain.AddBlock("Second Block")
	chain.AddBlock("Third Block")

	for _, block := range chain.Blocks {
		fmt.Printf("%x \n", block.Hash)
		fmt.Printf("%v \n", string(block.Data))
		fmt.Println()
	}

}
