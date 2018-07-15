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

// IsBlockValid returns if the block is valid by checking index, and comparing
// the hash of the previous block
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

// CalculateHash returns the hash of the the block using SHA256
func CalculateHash(block model.Block) string {

	// Represents the object as a string
	record := strconv.Itoa(block.Index) +
		block.Timestamp +
		strconv.Itoa(block.BPM) +
		block.PrevHash +
		block.Nonce

	// Hashes string that represented the block
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// GenerateBlock creates a new block using previous block's hash
func GenerateBlock(oldBlock model.Block, BPM int) model.Block {

	t := time.Now()

	// create new block
	var newBlock model.Block
	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Difficulty = model.Difficulty

	nonce := 0
	for {
		// create a hexadecimal which is to be used as a nonce from i
		hex := fmt.Sprintf("%x", nonce)
		newBlock.Nonce = hex

		// check if nonce of new block is valid
		if isHashValid(CalculateHash(newBlock), newBlock.Difficulty) {
			fmt.Println(CalculateHash(newBlock), " work done!")
			newBlock.Hash = CalculateHash(newBlock)

			// As nonce began at zero and is incremened once for each calculation,
			// the nonce value will equal the number of calculations done
			newBlock.NumCalculations = nonce
			break
		}
		nonce++ // incrementing allows a new nonce can be generated each time
	}

	return newBlock
}

// HashStr hashes a string using sha256
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
