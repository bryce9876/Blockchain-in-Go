# Blockchain in Go

Just a barebones blockchain in Go that can be run on a single server. 

Allows both GET and POST requests.
- GET requests to see the blockchain with all its individual blocks, each block containing informatoin like it's hash, value,
nonce etc.
- POST requests with a given value and valid passoword will have a new block with a valid nonce calculated on the serverside. 
The password will be hashed and compared against the hashed password server side. Note it is realised that it is usual for 
nonce calculaiton to be done client site and verified server side, but this is just a proof of concept. 


## How to use it
Just run it using "go run main.go"

Make sure you have Go installed first
