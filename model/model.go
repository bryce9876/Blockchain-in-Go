package model

/*********Constants********/
const Difficulty = 5
const PasswordHash = "9f735e0df9a1ddc702bf0a1a7b83033f9f7153a00c29de82cedadc9957289b05" // sha256 hash

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