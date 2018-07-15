package model

// Difficulty is how hard it is to generate a valid hash
const Difficulty = 5

// PasswordHash is hash of the client password that allow server side hash calculation
const PasswordHash = "9f735e0df9a1ddc702bf0a1a7b83033f9f7153a00c29de82cedadc9957289b05"

// Block represents each 'item' in the blockchain
type Block struct {
	Index           int
	Timestamp       string
	BPM             int
	Hash            string
	PrevHash        string
	Difficulty      int
	Nonce           string
	NumCalculations int
}

// Blockchain is a series of validated Blocks
var Blockchain []Block

// Message takes incoming JSON payload for writing heart rate
type Message struct {
	BPM      int
	Password string
}
