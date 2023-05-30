package mempool

type Transaction struct {
	Txid     string   `json:"txid"`
	Version  uint     `json:"version"`
	Locktime uint     `json:"locktime"`
	Vin      []TxIn   `json:"vin"`
	Vout     []TxOut  `json:"vout"`
	Size     uint     `json:"size"`
	Weight   uint     `json:"weight"`
	Fee      uint     `json:"uint"`
	Status   TxStatus `json:"status"`
}

type TxStatus struct {
	Confirmed   bool   `json:"confirmed"`
	BlockHeight uint   `json:"block_height"`
	BlockHash   string `json:"block_hash"`
	BlockTime   uint   `json:"block_time"`
}

type TxIn struct {
	Txid       string `json:"txid"`
	Vout       uint   `json:"vout"`
	Prevout    TxOut  `json:"prevout"`
	IsCoinbase bool   `json:"is_coinbase"`
	Sequence   uint64 `json:"sequence"`
}

type TxOut struct {
	Scriptpubkey        string `json:"scriptpubkey"`
	ScriptpubkeyType    string `json:"scriptpubkey_type"`
	ScriptpubkeyAddress string `json:"scriptpubkey_address"`
	Value               uint64 `json:"value"`
}
