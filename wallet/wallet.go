package wallet

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/YushinJung/NomadCoin/utils"
)

const (
	signature     string = "4208c63b8c2c6f930589b3a6781c23d7a59f8aab645f3d1080b4c58400fcd221ee3138f130fb345b831715bd9db688cbde739aaaafc4a4358da0962a5a9a36c8"
	privateKey    string = "30770201010420ce106e3456cf6f62ebcff0de3111445220d16397b7ad4564157f9b2149bc386aa00a06082a8648ce3d030107a14403420004e42f74822101f04e42dc8f2ef92fb10fb89ce949a9cb5c6936ac8685f1024e2d4266863cf00c7adc096db7e95e4641bee2bd6a6291619f5d9e68024de67bb007"
	hashedMessage string = "f3edaa4cf9ec6721ab8c65b1227637628261bc8d02fca2fc2609ab9636c4c503"
)

func Start() {
	// load private key
	privateByte, err := hex.DecodeString(privateKey)
	utils.HandleErr(err)
	// byte을 바로 넘겨줄 수도 있지만, 해당 function은 제대로된 형태인지 확인 안하고,
	// private key로 바꾸기 때문에, file이 잘 못 되어도 진행이 된다.
	// x509.ParseECPrivateKey([]byte(privateKey))
	restoredKey, err := x509.ParseECPrivateKey(privateByte)
	utils.HandleErr(err)

	// restore signature
	sigBytes, err := hex.DecodeString(signature)
	utils.HandleErr(err)
	rBytes := sigBytes[:len(sigBytes)/2]
	sBytes := sigBytes[len(sigBytes)/2:]

	var bigR, bigS = big.Int{}, big.Int{}
	bigR.SetBytes(rBytes)
	bigS.SetBytes(sBytes)

	// restore hash
	hashBytes, err := hex.DecodeString(hashedMessage)
	utils.HandleErr(err)

	// verify
	ok := ecdsa.Verify(&restoredKey.PublicKey, hashBytes, &bigR, &bigS)
	fmt.Println(ok)
}
