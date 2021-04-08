package blockchain

import (
	"github.com/dgraph-io/badger"
)

const (
	dbpath = "./tmp/blocks"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func InitBlockChain() *BlockChain {
	var lastHash []byte

	db, err := badger.Open(badger.DefaultOptions(dbpath))
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {

		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound {
			genesis := Genesis()

			err = txn.Set(genesis.Hash, genesis.Serialize())
			Handle(err)

			err = txn.Set([]byte("lh"), genesis.Hash)

			lastHash = genesis.Hash

			return err
		} else {
			item, err := txn.Get([]byte("lh"))
			Handle(err)
			lastHash, err = item.ValueCopy(nil)
			return err
		}
	})
	Handle(err)

	blockchain := BlockChain{lastHash, db}

	return &blockchain
}

func (chain *BlockChain) AddBlock(data string) {

	newBlock := CreateBlock(data, chain.LastHash)

	err := chain.Database.Update(func(txn *badger.Txn) error {

		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)

		err = txn.Set([]byte("lh"), newBlock.Hash)

		return err
	})
	Handle(err)

	chain.LastHash = newBlock.Hash
}

func (chain *BlockChain) Iterator() *BlockChainIterator {

	return &BlockChainIterator{chain.LastHash, chain.Database}
}

func (itr *BlockChainIterator) Next() *Block {
	var block *Block

	err := itr.Database.View(func(txn *badger.Txn) error {

		item, err := txn.Get(itr.CurrentHash)
		Handle(err)

		encodedBlock, err := item.ValueCopy(nil)
		block = DeSerialize(encodedBlock)

		return err
	})
	Handle(err)

	itr.CurrentHash = block.PrevHash

	return block
}
