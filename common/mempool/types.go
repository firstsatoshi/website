package mempool

// Transaction

// {
// 	"txid": "15e10745f15593a899cef391191bdd3d7c12412cc4696b7bcb669d0feadc8521",
// 	"version": 1,
// 	"locktime": 0,
// 	"vin": [
// 	  {
// 		"txid": "1fdfed84588cb826b876cd761ecebcf1726453437f0a6826e82ed54b2807a036",
// 		"vout": 12,
// 		"prevout": {
// 		  "scriptpubkey": "76a9142d110e1702a73c56fb6ea709cd529ea00680114388ac",
// 		  "scriptpubkey_asm": "OP_DUP OP_HASH160 OP_PUSHBYTES_20 2d110e1702a73c56fb6ea709cd529ea006801143 OP_EQUALVERIFY OP_CHECKSIG",
// 		  "scriptpubkey_type": "p2pkh",
// 		  "scriptpubkey_address": "157HqdxdT8dxTjeRLVT5HPtFc1LH4CeuVC",
// 		  "value": 686860000
// 		},
// 		"scriptsig": "483045022100bcdf40fb3b5ebfa2c158ac8d1a41c03eb3dba4e180b00e81836bafd56d946efd022005cc40e35022b614275c1e485c409599667cbd41f6e5d78f421cb260a020a24f01210255ea3f53ce3ed1ad2c08dfc23b211b15b852afb819492a9a0f3f99e5747cb5f0",
// 		"scriptsig_asm": "OP_PUSHBYTES_72 3045022100bcdf40fb3b5ebfa2c158ac8d1a41c03eb3dba4e180b00e81836bafd56d946efd022005cc40e35022b614275c1e485c409599667cbd41f6e5d78f421cb260a020a24f01 OP_PUSHBYTES_33 0255ea3f53ce3ed1ad2c08dfc23b211b15b852afb819492a9a0f3f99e5747cb5f0",
// 		"is_coinbase": false,
// 		"sequence": 4294967295
// 	  }
// 	],
// 	"vout": [
// 	  {
// 		"scriptpubkey": "76a91472d52e2f5b88174c35ee29844cce0d6d24b921ef88ac",
// 		"scriptpubkey_asm": "OP_DUP OP_HASH160 OP_PUSHBYTES_20 72d52e2f5b88174c35ee29844cce0d6d24b921ef OP_EQUALVERIFY OP_CHECKSIG",
// 		"scriptpubkey_type": "p2pkh",
// 		"scriptpubkey_address": "1BUBQuPV3gEV7P2XLNuAJQjf5t265Yyj9t",
// 		"value": 1240000000
// 	  }
// 	],
// 	"size": 884,
// 	"weight": 3536,
// 	"fee": 20000,
// 	"status": {
// 	  "confirmed": true,
// 	  "block_height": 363348,
// 	  "block_hash": "0000000000000000139385d7aa78ffb45469e0c715b8d6ea6cb2ffa98acc7171",
// 	  "block_time": 1435754650
// 	}
//   }

type Transaction struct {
	Txid     string   `json:"txid"`
	Version  uint64   `json:"version"`
	Locktime uint64   `json:"locktime"`
	Vin      []TxIn   `json:"vin"`
	Vout     []TxOut  `json:"vout"`
	Size     uint64   `json:"size"`
	Weight   uint64   `json:"weight"`
	Fee      uint64   `json:"fee"`
	Status   TxStatus `json:"status"`
}

type TxStatus struct {
	Confirmed   bool   `json:"confirmed"`
	BlockHeight uint64 `json:"block_height"`
	BlockHash   string `json:"block_hash"`
	BlockTime   uint64 `json:"block_time"`
}

type TxIn struct {
	Txid       string `json:"txid"`
	Vout       uint64 `json:"vout"`
	Prevout    TxOut  `json:"prevout"`
	IsCoinbase bool   `json:"is_coinbase"`
	Sequence   uint64 `json:"sequence"`
}

type TxOut struct {
	Scriptpubkey        string `json:"scriptpubkey"`
	ScriptpubkeyAsm     string `json:"scriptpubkey_asm"`
	ScriptpubkeyType    string `json:"scriptpubkey_type"`
	ScriptpubkeyAddress string `json:"scriptpubkey_address"`
	Value               uint64 `json:"value"`
}

// Address
// https://mempool.space/zh/docs/api/rest#get-address
//
//	{
//		address: "1wiz18xYmhRX6xStj2b9t1rwWX4GKUgpv",
//		chain_stats: {
//		  funded_txo_count: 5,
//		  funded_txo_sum: 15007599040,
//		  spent_txo_count: 5,
//		  spent_txo_sum: 15007599040,
//		  tx_count: 7
//		},
//		mempool_stats: {
//		  funded_txo_count: 0,
//		  funded_txo_sum: 0,
//		  spent_txo_count: 0,
//		  spent_txo_sum: 0,
//		  tx_count: 0
//		}
//	  }
type AddressDetail struct {
	Address      string `json:"address"`
	ChainStats   Stats  `json:"chain_stats"`
	MempoolStats Stats  `json:"mempool_stats"`
}

type Stats struct {
	FundedTxoCount uint64 `json:"funded_txo_count"`
	FundedTxoSum   uint64 `json:"funded_txo_sum"`
	SpentTxoCount  uint64 `json:"spent_txo_count"`
	SpentTxoSum    uint64 `json:"spent_txo_sum"`
	TxCount        uint64 `json:"tx_count"`
}

// recommended
//
// {
//   fastestFee: 1,
//   halfHourFee: 1,
//   hourFee: 1,
//   economyFee: 1,
//   minimumFee: 1
// }

type RecommendedFee struct {
	FastestFee  uint64 `json:"fastestFee"`
	HalfHourFee uint64 `json:"halfHourFee"`
	HourFee     uint64 `json:"hourFee"`
	EconomyFee  uint64 `json:"economyFee"`
	MinimumFee  uint64 `json:"minimumFee"`
}

// UTXO
// [
//   {
//     "txid": "2e3972ab2ee9c4c36bab0e436be7f1ecb75f9f9bb4a36e70d5e322be36bef2ee",
//     "vout": 1,
//     "status": {
//       "confirmed": true,
//       "block_height": 792060,
//       "block_hash": "00000000000000000004b0e40b16b7c64397257407cd4dfb8c70aa6cd836a7b2",
//       "block_time": 1685426180
//     },
//     "value": 665893692
//   },
//   {
//     "txid": "a867b953fdb44e4eab42f630d779242e46baca1b5ec957adc8d75e3a4a3ddcd8",
//     "vout": 1,
//     "status": {
//       "confirmed": true,
//       "block_height": 791950,
//       "block_hash": "0000000000000000000074a796c4793cf7f0182c53dde6e81c8b109e9bd33bb5",
//       "block_time": 1685369281
//     },
//     "value": 671732724
//   }
// ]

type UTXO struct {
	Txid     string   `json:"txid"`
	Vout     uint64   `json:"vout"`
	TxStatus TxStatus `json:"status"`
	Value    uint64   `json:"value"`
}
