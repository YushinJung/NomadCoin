package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"

	"github.com/YushinJung/NomadCoin/utils"
)

const (
	signature     string = "4208c63b8c2c6f930589b3a6781c23d7a59f8aab645f3d1080b4c58400fcd221ee3138f130fb345b831715bd9db688cbde739aaaafc4a4358da0962a5a9a36c8"
	privateKey    string = "30770201010420ce106e3456cf6f62ebcff0de3111445220d16397b7ad4564157f9b2149bc386aa00a06082a8648ce3d030107a14403420004e42f74822101f04e42dc8f2ef92fb10fb89ce949a9cb5c6936ac8685f1024e2d4266863cf00c7adc096db7e95e4641bee2bd6a6291619f5d9e68024de67bb007"
	hashedMessage string = "f3edaa4cf9ec6721ab8c65b1227637628261bc8d02fca2fc2609ab9636c4c503"
)

func Start() {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	// privateKey 에 public key 가 들어있는 구조
	keyAsBytes, err := x509.MarshalECPrivateKey(privateKey)

	fmt.Printf("%x\n", keyAsBytes)

	utils.HandleErr(err)

	byteHash, err := hex.DecodeString(hashedMessage)

	utils.HandleErr(err)

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, byteHash)

	signature := append(r.Bytes(), s.Bytes()...)

	fmt.Printf("%x\n", signature)

	fmt.Println(r.Bytes(), s.Bytes())

	utils.HandleErr(err)

	ok := ecdsa.Verify(&privateKey.PublicKey, byteHash, r, s)

	fmt.Print(ok)
}
