# Blockchain in Go

Just a barebones blockchain in Go that can be run on a single server. 

Allows both GET and POST requests.
- GET requests to see the blockchain with all its individual blocks, each block containing informatoin like it's hash, value,
nonce etc.
- POST requests with a given value and valid passoword will have a new block with a valid nonce calculated on the serverside. 
The password will be hashed and compared against the hashed password server side. Currently the password is preset to "testpassword". Note it is realised that it is usual for nonce calculaiton to be done client site and verified server side, but this is just a proof of concept. 


## How to use it
- Make sure you have Go installed
- Clone repo
- Make a file called .env, and put the address of the port you want to run the server on in this file e.g. ADDR=8080
- Then just run it using "go run main.go"
