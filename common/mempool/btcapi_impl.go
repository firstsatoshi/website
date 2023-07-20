package mempool

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/firstsatoshi/website/common/btcapi"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
)

// implement  BTCAPIClient interface
var _ btcapi.BTCAPIClient = &MempoolApiClient{}

// ListUnspent get utox of address
func (m *MempoolApiClient) ListUnspent(address btcutil.Address) ([]*btcapi.UnspentOutput, error) {

	utxos, err := m.GetAddressUTXOs(address.EncodeAddress())
	if err != nil {
		return []*btcapi.UnspentOutput{}, err
	}

	unspentOutputs := make([]*btcapi.UnspentOutput, 0)
	for _, utxo := range utxos {
		txHash, err := chainhash.NewHashFromStr(utxo.Txid)
		if err != nil {
			return nil, err
		}
		pkScript, err := txscript.PayToAddrScript(address)
		if err != nil {
			return nil, err
		}
		unspentOutputs = append(unspentOutputs, &btcapi.UnspentOutput{
			Outpoint: wire.NewOutPoint(txHash, uint32(utxo.Vout)),
			// Output:   wire.NewTxOut(int64(utxo.Value), address.ScriptAddress()),
			Output: wire.NewTxOut(int64(utxo.Value), pkScript),
		})
	}
	return unspentOutputs, nil
}

// GetRawTransaction  https://mempool.space/zh/docs/api/rest#get-transaction-raw
func (m *MempoolApiClient) GetRawTransaction(txHash *chainhash.Hash) (*wire.MsgTx, error) {
	// https://mempool.space/api/tx/15e10745f15593a899cef391191bdd3d7c12412cc4696b7bcb669d0feadc8521/raw

	url := fmt.Sprintf("%s/tx/%v/raw", m.host, txHash.String())
	resp, err := m.client.R().Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		err = fmt.Errorf("code:%v,error:%v", resp.StatusCode(), string(resp.Body()))
		return nil, err
	}
	tx := wire.NewMsgTx(wire.TxVersion)
	rawtx := resp.Body()
	if len(rawtx) == 0 {
		return nil, fmt.Errorf("empty rawtx binary data reponse")
	}
	if err := tx.Deserialize(bytes.NewReader(rawtx)); err != nil {
		return nil, err
	}
	return tx, nil
}

// BroadcastTx https://mempool.space/zh/docs/api/rest#post-transaction
func (m *MempoolApiClient) BroadcastTx(tx *wire.MsgTx) (txHash *chainhash.Hash, err error) {

	var buf bytes.Buffer
	if err := tx.Serialize(&buf); err != nil {
		logx.Errorf("tx.Serialize error: %v", err.Error())
		return nil, err
	}

	url := fmt.Sprintf("%s/tx", m.host)
	// fmt.Printf("url: %v", url)
	txHex := hex.EncodeToString(buf.Bytes())
	logx.Infof("BroadcastTx rawtx: %v", txHex)

	rand.Seed(time.Now().UnixMilli())
	txHash = nil
	txid := ""
	for iTry := 0; iTry < 120; iTry++ {
		logx.Infof(" BroadcastTx  try %v times , hexRawTx: %v", iTry, txHex)
		sleepSecs := iTry + 5
		time.Sleep(time.Second * time.Duration(sleepSecs))

		resp, err := m.client.R().SetBody([]byte(hex.EncodeToString(buf.Bytes()))).Post(url)
		if err != nil {
			continue
		}
		if resp.StatusCode() != http.StatusOK {
			continue
		}

		txid = string(resp.Body())
		if len(txid) != 64 {
			continue
		}
		_, e := hex.DecodeString(txid)
		if e != nil {
			continue
		}

		// everything is ok
		txHash, err := chainhash.NewHashFromStr(txid)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to parse tx hash, %s", string(txid)))
		}
		logx.Infof("broadcasd succeed, txid: %v", txHash)
		return txHash, nil
	}

	logx.Infof("BroadcastTx rawtx error %v ", err.Error())

	return nil, err
}

// BroadcastTx https://mempool.space/zh/docs/api/rest#post-transaction
func (m *MempoolApiClient) BroadcastTxHex(hexRawTx string) (txHash *chainhash.Hash, err error) {

	url := fmt.Sprintf("%s/tx", m.host)
	txHash = nil
	txid := ""
	rand.Seed(time.Now().UnixMilli())
	for iTry := 0; iTry < 120; iTry++ {
		logx.Infof(" BroadcastTx  try %v times , hexRawTx : %v", iTry, hexRawTx)
		sleepSecs := iTry + 5
		time.Sleep(time.Second * time.Duration(sleepSecs))

		resp, err := m.client.R().SetBody([]byte(hexRawTx)).Post(url)
		if err != nil {
			continue
		}
		if resp.StatusCode() != http.StatusOK {
			continue
		}

		txid = string(resp.Body())
		if len(txid) != 64 {
			continue
		}
		_, e := hex.DecodeString(txid)
		if e != nil {
			continue
		}

		// everything is ok
		txHash, err := chainhash.NewHashFromStr(txid)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to parse tx hash, %s", string(txid)))
		}
		logx.Infof("txid: %v", txHash)
		return txHash, nil
	}

	logx.Infof("BroadcastTx rawtx error %v ", err.Error())

	return nil, err
}
