package blockchain

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

const Difficulty = 16

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, 256-Difficulty)
	return &ProofOfWork{b, target}
}

func (pow *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)

	return data
}

func (pow *ProofOfWork) Run() (int, []byte) {

	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	for ; nonce < math.MaxInt64; nonce++ {

		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\r %x", hash)

		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.Target) == -1 {
			break
		}
	}

	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.InitData(pow.Block.Nonce)
	hash := sha256.Sum256(data)

	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.Target) == -1
}

func ToHex(num int64) []byte {
	hex := fmt.Sprintf("%x", num)
	return []byte(hex)
}
