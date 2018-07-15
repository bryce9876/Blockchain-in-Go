package api

import (
	"Blockchain/blockchainhelpers"
	"Blockchain/model"
	"encoding/json"
	"io"
	"net/http"
	"sync"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
)

// create handlers
func MakeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	return muxRouter
}

// handles GET request by responding with a JSON representaion of the current blockchain
func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {

	// create a json representation of the current blockchain with indentations
	bytes, err := json.MarshalIndent(model.Blockchain, "", "  ")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, string(bytes))
}

// takes JSON payload as an input for heart rate (BPM)
func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var m model.Message

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
	var mutex = &sync.Mutex{}
	mutex.Lock()
	newBlock := blockchainhelpers.GenerateBlock(model.Blockchain[len(model.Blockchain)-1], m.BPM)
	mutex.Unlock()

	if blockchainhelpers.IsBlockValid(newBlock, model.Blockchain[len(model.Blockchain)-1]) {
		model.Blockchain = append(model.Blockchain, newBlock)
		spew.Dump(model.Blockchain)
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

// check that the hash of the input password matches the stored hash
func Authenticate(usrPassword string) bool {
	return blockchainhelpers.HashStr(usrPassword) != model.PasswordHash
}
