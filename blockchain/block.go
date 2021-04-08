package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

type Block struct {
	Hash         []byte
	Transactions []*Transaction
	PrevHash     []byte
	Nonce        int
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}

	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

func CreateBlock(txns []*Transaction, prevHash []byte) *Block {
	block := &Block{[]byte{}, txns, prevHash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash
	block.Nonce = nonce

	return block
}

func Genesis(coinbase *Transaction) *Block {
	return CreateBlock([]*Transaction{coinbase}, []byte{})
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
