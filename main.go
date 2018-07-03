package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

/*********Constants********/
const difficulty = 5
const passwordHash = "9f735e0df9a1ddc702bf0a1a7b83033f9f7153a00c29de82cedadc9957289b05" // sha256 hash

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

var mutex = &sync.Mutex{}

func main() {

	// load enviroment from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		t := time.Now()

		// Create genesis block which is always the very first block
		genesisBlock := Block{}
		genesisBlock = Block{0, t.String(), 0, calculateHash(genesisBlock), "", difficulty, "", -1}
		spew.Dump(genesisBlock)

		// Need to lock to prevent two seperate clients both appending their own node to the same block
		mutex.Lock()
		Blockchain = append(Blockchain, genesisBlock)
		mutex.Unlock()
	}()

	// run the web server
	log.Fatal(run())

}

// set up and run the server
func run() error {

	mux := makeMuxRouter()

	// Remeber ADDR is in the .env file
	httpPort := os.Getenv("ADDR")
	log.Println("HTTP Server Listening on port :", httpPort)
	s := &http.Server{
		Addr:           ":" + httpPort,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

// create handlers
func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	return muxRouter
}

// handles GET request by responding with a JSON representaion of the current blockchain
func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {

	// create a json representation of the current blockchain with indentations
	bytes, err := json.MarshalIndent(Blockchain, "", "  ")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(bytes))
}

// takes JSON payload as an input for heart rate (BPM)
func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var m Message

	// Decode http request into message struct
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}

	// checks if the password is correct
	// if !authenticate(m.Password) {
	// 	respondWithJSON(w, r, http.StatusUnauthorized, r.Body)
	// }

	defer r.Body.Close()

	//ensure atomicity when creating new block
	mutex.Lock()
	newBlock := generateBlock(Blockchain[len(Blockchain)-1], m.BPM)
	mutex.Unlock()

	if isBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
		Blockchain = append(Blockchain, newBlock)
		spew.Dump(Blockchain)
	}

	respondWithJSON(w, r, http.StatusCreated, newBlock)

}

// responds to the POST request with a http status, and if this POST is successful, it responds
// with a payload of the new block created.
func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

// make sure block is valid by checking index, and comparing the hash of the previous block
func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

// check that the hash of the input password matches the stored hash
func authenticate(usrPassword string) bool {
	return hashStr(usrPassword) != passwordHash
}

// hash the block using SHA256
func calculateHash(block Block) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + strconv.Itoa(block.BPM) + block.PrevHash + block.Nonce
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// hashes a string using sha256 encyption
func hashStr(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// create a new block using previous block's hash
func generateBlock(oldBlock Block, BPM int) Block {

	t := time.Now()

	// create new block
	var newBlock Block
	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Difficulty = difficulty

	i := 0
	for {
		// create a hexadecimal which is to be used as a nonce from i
		hex := fmt.Sprintf("%x", i)
		newBlock.Nonce = hex

		// check if nonce of new block is valid
		if isHashValid(calculateHash(newBlock), newBlock.Difficulty) {
			fmt.Println(calculateHash(newBlock), " work done!")
			newBlock.Hash = calculateHash(newBlock)
			newBlock.NumCalculations = i
			break
		}
		i++ // increase by a 1 so a new nonce can be generated each time
	}

	newBlock.NumCalculations = i
	return newBlock
}

// checks if a hash is valid depending on the difficulty
// a hash is considered valid if the number of leading zero's is greater than
// equal to the difficulty
func isHashValid(hash string, difficulty int) bool {
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hash, prefix)
}
