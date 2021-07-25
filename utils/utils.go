package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

func HandleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func ToBytes(i interface{}) []byte {
	// get any type as interface
	var aBuffer bytes.Buffer
	encoder := gob.NewEncoder(&aBuffer)
	// save the result to block buffer
	HandleErr(encoder.Encode(i))
	return aBuffer.Bytes()
}

func FromBytes(i interface{}, data []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	HandleErr(decoder.Decode(i))
}

func Hash(i interface{}) string {
	s := fmt.Sprintf("%v", i)
	// v is default formatter for interface
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash)
}
