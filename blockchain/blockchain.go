package blockchain

import (
	"encoding/hex"
	"fmt"
	"os"
	"runtime"

	"github.com/dgraph-io/badger"
)

const (
	dbPath      = "./tmp/blocks"
	dbFile      = "./tmp/blocks/MANIFEST"
	genesisData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks."
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func DBexists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

func InitBlockChain(address string) *BlockChain {

	if DBexists() {
		fmt.Println("Blockchain already exists")
		runtime.Goexit()
	}

	var lastHash []byte

	db, err := badger.Open(badger.DefaultOptions(dbPath))
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {

		coinbase := CoinbaseTx(address, genesisData)
		genesis := Genesis(coinbase)

		err = txn.Set(genesis.Hash, genesis.Serialize())
		Handle(err)

		err = txn.Set([]byte("lh"), genesis.Hash)

		lastHash = genesis.Hash

		return err
	})
	Handle(err)

	blockchain := BlockChain{lastHash, db}

	return &blockchain
}

func ContinueBlockChain() *BlockChain {
	if DBexists() == false {
		fmt.Println("No existing blockchain found, create one!")
		runtime.Goexit()
	}

	var lastHash []byte

	db, err := badger.Open(badger.DefaultOptions(dbPath))
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		lastHash, err = item.ValueCopy(nil)

		return err
	})
	Handle(err)

	chain := BlockChain{lastHash, db}

	return &chain
}

func (chain *BlockChain) AddBlock(txns []*Transaction) {

	newBlock := CreateBlock(txns, chain.LastHash)

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

type TXOInfoPair struct {
	txID   string
	outIdx int
}

type UTXO struct {
	Output TxOutput
	txID   string
	outIdx int
}

func (chain *BlockChain) FindUTXOs(address string) []UTXO {
	var utxos []UTXO

	spentTXOs := make(map[TXOInfoPair]bool)

	iter := chain.Iterator()
	for {
		block := iter.Next()

		for _, tx := range block.Transactions {

			if !tx.IsCoinbase() {
				for _, in := range tx.Inputs {
					if in.CanUnlock(address) {
						inTxID := hex.EncodeToString(in.ID)
						spentTXOs[TXOInfoPair{inTxID, in.Out}] = true
					}
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	iter = chain.Iterator()
	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

			for outIdx, out := range tx.Outputs {

				if out.CanBeUnlocked(address) && !spentTXOs[TXOInfoPair{txID, outIdx}] {

					utxo := UTXO{out, txID, outIdx}

					utxos = append(utxos, utxo)
				}
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return utxos
}

func (chain *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	utxos := chain.FindUTXOs(address)
	accumulated := 0

	unspentOuts := make(map[string][]int)

Work:
	for _, utxo := range utxos {
		txID := utxo.txID
		out := utxo.Output

		if out.CanBeUnlocked(address) && accumulated < amount {
			accumulated += out.Value
			unspentOuts[txID] = append(unspentOuts[txID], utxo.outIdx)

			if accumulated >= amount {
				break Work
			}
		}
	}

	return accumulated, unspentOuts
}
