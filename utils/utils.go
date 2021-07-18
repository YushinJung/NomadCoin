package utils

import (
	"bytes"
	"encoding/gob"
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
