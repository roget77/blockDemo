package blockDemo

import (
	"github.com/boltdb/bolt"
)

const dbFile string = "boldb.file"
// 创世区块基础信息
const genesisCoinbaseData string = "genesisCoinbaseData"

// type Blockchain struct {
// 	blocks []*Block
// }

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.db}

	return bci
}

func (bc *Blockchain) CloseDb() {
	bc.db.Close()
}

func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)
		return nil
	})

	if err != nil {
		panic(err)
	}

	i.currentHash = block.PrevBlockHash

	return block
}

// func (bc *Blockchain) AddBlock(data string) {
// 	prevBlock := bc.blocks[len(bc.blocks)-1]
// 	newBlock := NewBlock(data, prevBlock.Hash)
// 	bc.blocks = append(bc.blocks, newBlock)
// }

func (bc *Blockchain) AddBlock(newBlock *Block) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		panic("db view fail")
	}

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

// 未花费交易输出 可计算余额
func (bc *Blockchain) FindUTXO(address string) []TXOutput {
	var UTXOs []TXOutput
	unspentTransactions := bc.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

// 读取区块链，不存在则根据地址创建创世块
func NewBlockchain(address string) *Blockchain {
	// return &Blockchain{[]*Block{NewGenesisBlock()}}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			// genesis := NewGenesisBlock()

			cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
			genesis := NewGenesisBlock(cbtx)

			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				panic(err)
			}
			err = b.Put(genesis.Hash, genesis.Serialize())
			err = b.Put([]byte("l"), genesis.Hash)
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	bc := Blockchain{tip, db}

	return &bc

}

func BlockchainRun(address string) *Blockchain {
	bc := NewBlockchain(address)
	return bc
}
