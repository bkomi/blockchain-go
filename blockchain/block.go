package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash
	block.Nonce = nonce

	return block
}

func Genesis() *Block {
	return CreateBlock("The Times 03/Jan/2009 Chancellor on brink of second bailout for banks.", []byte{})
}

func (b *Block) Serialize() []byte {
	var buffer bytes.Buffer        // Stand-in
	enc := gob.NewEncoder(&buffer) // Will write to buffer
	// Encode the value.
	err := enc.Encode(b)
	Handle(err)

	return buffer.Bytes()
}

func DeSerialize(data []byte) *Block {
	var block Block

	dec := gob.NewDecoder(bytes.NewBuffer(data))

	err := dec.Decode(&block)
	Handle(err)

	return &block
}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
