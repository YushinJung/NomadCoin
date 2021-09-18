package blockchain

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/YushinJung/NomadCoin/utils"
	"github.com/YushinJung/NomadCoin/wallet"
)

const (
	minerReward int = 50
)

type mempool struct {
	Txs []*Tx
	m   sync.Mutex
}

var m *mempool = &mempool{}
var memOnce sync.Once

func Mempool() *mempool {
	memOnce.Do(func() {
		m = &mempool{}
	})
	return m
}

type Tx struct {
	ID        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txins"`
	TxOuts    []*TxOut `json:"txouts"`
}

type TxIn struct { // 특정 txid 에서 어떤 index 의 output을 사용하는지 알려줌.
	TxID      string `json:"txid"`
	Index     int    `json:"index"`
	Signature string `json:"signature"`
}

type TxOut struct {
	Address string `json:"address"`
	Amount  int    `json:"amount"`
}

type UTxOut struct {
	TxID   string `json:"txid"`
	Index  int    `json:"index"`
	Amount int    `json:"amount"`
}

func (t *Tx) getId() {
	t.ID = utils.Hash(t)
}

func (t *Tx) sign() {
	for _, txIn := range t.TxIns {
		txIn.Signature = wallet.Sign(t.ID, wallet.Wallet())
	}
}

func validate(tx *Tx) bool {
	// validate person who makes transaction
	// which means validating the output of the transaction
	valid := true
	// input에 참조된 output이 우리가 소유하고 있음을 확인하고 싶음.
	for _, txIn := range tx.TxIns {
		// txIn 의 TxID는 사용한 transaction output을 만든 transaction id 임.
		prevTx := FindTx(Blockchain(), txIn.TxID)
		if prevTx == nil {
			valid = false
			break
		}
		address := prevTx.TxOuts[txIn.Index].Address //public key
		valid = wallet.Verify(txIn.Signature, tx.ID, address)
		// tx의 signature 와 ID를 address(public key)로 증명
		if !valid {
			break
		}
	}
	return valid
}

func isOnMempool(uTxOut *UTxOut) bool {
	exists := false
Outer:
	for _, tx := range Mempool().Txs {
		for _, input := range tx.TxIns {
			if input.TxID == uTxOut.TxID && input.Index == uTxOut.Index {
				exists = true
				break Outer
			}
		}
	}
	return exists
}

func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"", -1, "COINBASE"},
	}
	txOuts := []*TxOut{
		{address, minerReward},
	}
	tx := Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return &tx
}

var ErrorNoMoney = errors.New("not enogh money")
var ErrorNotValid = errors.New("Tx not valid")

func makeTx(from, to string, amount int) (*Tx, error) {
	if BalanceByAddress(from, Blockchain()) < amount {
		return nil, ErrorNoMoney
	}
	var txOuts []*TxOut
	var txIns []*TxIn
	total := 0
	uTxOuts := UTxOutsByAddress(from, Blockchain())
	for _, uTxOut := range uTxOuts {
		if total >= amount { // if we have enought amount break
			break
		}
		// still "from" is not secured yet
		txIn := &TxIn{uTxOut.TxID, uTxOut.Index, from}
		txIns = append(txIns, txIn)
		total += uTxOut.Amount
	}
	// for "from" address
	if change := total - amount; change != 0 {
		changeTxOut := &TxOut{from, change}
		txOuts = append(txOuts, changeTxOut)
	}
	// for "to" address
	txOut := &TxOut{to, amount}
	txOuts = append(txOuts, txOut)
	tx := &Tx{
		ID:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	tx.sign()
	valid := validate(tx)
	if !valid {
		return nil, ErrorNotValid
	}
	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) (*Tx, error) {
	// be called by API
	// for whom to send and amount to send
	// if transaction cannot be made error will be called
	tx, err := makeTx(wallet.Wallet().Address, to, amount) // Yushin 대신 address가 들어가게 될 것
	if err != nil {
		return nil, err
	}
	m.Txs = append(m.Txs, tx)
	return tx, nil
}

func (m *mempool) TxToConfirm() []*Tx {
	coinbase := makeCoinbaseTx(wallet.Wallet().Address) // coin 채굴 시
	txs := m.Txs                                        // mempool의 모든 transaction을
	txs = append(txs, coinbase)                         // 하나로 합쳐서 전달
	m.Txs = nil                                         // mempool은 비우자
	return txs
}

func StatusMempool(rw http.ResponseWriter) {
	m.m.Lock()
	defer m.m.Unlock()
	utils.HandleErr(json.NewEncoder(rw).Encode(m.Txs))
}

func (m *mempool) AddPeerTx(tx *Tx) {
	m.m.Lock()
	defer m.m.Unlock()
	m.Txs = append(m.Txs, tx)
}
