// 工作量证明模块
package blockDemo

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"strconv"
)

/**
在比特币中，当一个块被挖出来以后，“target bits” 代表了区块头里存储的难度，也就是开头有多少个 0。
这里的 24 指的是算出来的哈希前 24 位必须是 0，如果用 16 进制表示，就是前 6 位必须是 0，
这一点从最后的输出可以看出来。目前我们并不会实现一个动态调整目标的算法，
所以将难度定义为一个全局的常量即可。
 */
const targetBits = 24

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}

	return pow
}

func IntToHex(num int64) []byte {
	return []byte(strconv.FormatInt(num, 16))
}

/**
我们需要有数据来进行哈希，准备数据
 */
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

/**
在这个循环中，我们做的事情有：
1.准备数据
2.用 SHA-256 对数据进行哈希
3.将哈希转换成一个大整数
4.将这个大整数与目标进行比较
 */
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)
	for nonce < math.MaxInt64 {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			fmt.Printf("\r%x", hash)
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}

/**
对工作量证明进行验证
 */
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}






