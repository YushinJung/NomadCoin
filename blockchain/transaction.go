package blockchain

import (
	"errors"
	"time"

	"github.com/YushinJung/NomadCoin/utils"
)

const (
	minerReward int = 50
)

type mempool struct {
	Txs []*Tx
}

var Mempool *mempool = &mempool{}

type Tx struct {
	Id        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txins"`
	TxOuts    []*TxOut `json:"txouts"`
}

func (t *Tx) getId() {
	t.Id = utils.Hash(t)
}

type TxIn struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"COINBASE", minerReward},
	}
	txOuts := []*TxOut{
		{address, minerReward},
	}
	tx := Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return &tx
}

func makeTx(from, to string, amount int) (*Tx, error) {
	// user 의 전체 transaction output을 통해 balance를 확인하면 다음 transaction의 input이 된다.
	// 충분히 갖고 있는지 확인해보자
	if Blockchain().BalanceByAddress(from) < amount {
		return nil, errors.New("not enough money")
	}
	// 모든 output을 고민할 필요 없이 전달할 amount 까지만 확인하면 된다.
	var txIns []*TxIn
	var txOuts []*TxOut

	total := 0
	oldTxOuts := Blockchain().TxOutsByAddress(from)
	for _, txOut := range oldTxOuts {
		if total > amount {
			break
		}
		txIn := &TxIn{txOut.Owner, txOut.Amount}
		txIns = append(txIns, txIn)
		total += txOut.Amount
	}
	change := total - amount
	if change != 0 {
		changeTxOut := &TxOut{
			Owner:  from,
			Amount: change,
		}
		txOuts = append(txOuts, changeTxOut)
	}
	txOut := &TxOut{
		Owner:  to,
		Amount: amount,
	}
	txOuts = append(txOuts, txOut)
	tx := &Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) error {
	// be called by API
	// for whom to send and amount to send
	// if transaction cannot be made error will be called
	tx, err := makeTx("Yushin", to, amount) // Yushin 대신 address가 들어가게 될 것
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}
