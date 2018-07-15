package blockchainhelpers

import (
	"Blockchain/model"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// make sure block is valid by checking index, and comparing the hash of the previous block
func IsBlockValid(newBlock, oldBlock model.Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if CalculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

// hash the block using SHA256
func CalculateHash(block model.Block) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + strconv.Itoa(block.BPM) + block.PrevHash + block.Nonce
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// create a new block using previous block's hash
func GenerateBlock(oldBlock model.Block, BPM int) model.Block {

	t := time.Now()

	// create new block
	var newBlock model.Block
	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Difficulty = model.Difficulty

	i := 0
	for {
		// create a hexadecimal which is to be used as a nonce from i
		hex := fmt.Sprintf("%x", i)
		newBlock.Nonce = hex

		// check if nonce of new block is valid
		if isHashValid(CalculateHash(newBlock), newBlock.Difficulty) {
			fmt.Println(CalculateHash(newBlock), " work done!")
			newBlock.Hash = CalculateHash(newBlock)
			newBlock.NumCalculations = i
			break
		}
		i++ // increase by a 1 so a new nonce can be generated each time
	}

	newBlock.NumCalculations = i
	return newBlock
}

// hashes a string using sha256 encyption
func HashStr(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// checks if a hash is valid depending on the difficulty
// a hash is considered valid if the number of leading zero's is greater than
// equal to the difficulty
func isHashValid(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}
