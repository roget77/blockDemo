package blockDemo

import (
	"bytes"
	"encoding/gob"
	"github.com/boltdb/bolt"
	"time"
)

const blocksBucket string = "blocksBucket"

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

// func (b *Block) SetHash() {
//     timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
//     headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
//     hash := sha256.Sum256(headers)

//     b.Hash = hash[:]
// }

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block

}

// func (bc *Blockchain) AddBlock(data string) {
// 	prevBlock := bc.blocks[len(bc.blocks)-1]
// 	newBlock := NewBlock(data, prevBlock.Hash)
// 	bc.blocks = append(bc.blocks, newBlock)
// }

func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		panic("db view fail")
	}

	newBlock := NewBlock(data, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if err := b.Put(newBlock.Hash, newBlock.Serialize());err != nil {
			panic("db put fail")
		}
		err = b.Put([]byte("l"), newBlock.Hash)
		bc.tip = newBlock.Hash

		return nil
	})
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	if err := encoder.Encode(b);err != nil {
		panic("encoder fail")
	}


	return result.Bytes()
}

func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	if err := decoder.Decode(&block);err != nil {
		panic("decode fail")
	}

	return &block
}