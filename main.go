package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

type BlockChain struct {
	blocks []*Block
}

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
}

func (b *Block) DeriveHash() {
	key := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	hash := sha256.Sum256(key)
	b.Hash = hash[:]
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash}
	block.DeriveHash()
	return block
}

func (chain *BlockChain) AddBlock(data string) {
	lastBlock := chain.blocks[len(chain.blocks)-1]
	newBlock := CreateBlock(data, lastBlock.Hash)
	chain.blocks = append(chain.blocks, newBlock)
}

func Genesis() *Block {
	return CreateBlock("The Times 03/Jan/2009 Chancellor on brink of second bailout for banks.", []byte{})
}

func InitBlockChain() *BlockChain {
	return &BlockChain{[]*Block{Genesis()}}
}

func main() {

	chain := InitBlockChain()

	chain.AddBlock("First Block")
	chain.AddBlock("Second Block")
	chain.AddBlock("Third Block")

	for _, block := range chain.blocks {
		fmt.Printf("%x \n", block.Hash)
		fmt.Printf("%v \n", string(block.Data))
		fmt.Printf("%x \n", block.PrevHash)
		fmt.Println()
	}

}
