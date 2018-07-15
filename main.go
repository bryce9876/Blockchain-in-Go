package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"Blockchain/api"
	"Blockchain/blockchainhelpers"
	"Blockchain/model"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

func main() {

	// load enviroment from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		t := time.Now()

		// Create genesis block which is always the very first block
		genesisBlock := model.Block{}
		genesisBlock = model.Block{0, t.String(), 0, blockchainhelpers.CalculateHash(genesisBlock), "", model.Difficulty, "", -1}
		spew.Dump(genesisBlock)

		// add genesis block to blockchain
		model.Blockchain = append(model.Blockchain, genesisBlock)
	}()

	// run the web server
	log.Fatal(run())
}

// set up and run the server
func run() error {

	mux := api.MakeMuxRouter()

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
