package blockchain

type block struct {
	// save data in the blcok only
	data     string
	hash     string // connected
	prevHash string
}

// hash function is one-way function
// if the data of hash is changed in the middle
// all of the hash will be changed after that hash

type blockchain struct {
	blocks []block
}

var b *blockchain // will not be shared

// singletone pattern
// If main function wants a blockchain
// it will return current blockchain instance without making another one.
func GetBlockchain() *blockchain {
	if b == nil {
		// check b is initialized
		b = &blockchain{}
	}
	return b
}
