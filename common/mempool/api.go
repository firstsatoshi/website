package mempool

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

type MempoolApiClient struct {
	host   string
	client *resty.Client
}

func NewMempoolApiClient(host string) *MempoolApiClient {

	if host == "" {
		host = "https://mempool.space/api"
	}

	client := resty.New()
	// client.BaseURL = host
	return &MempoolApiClient{
		host:   host,
		client: client,
	}
}

// GetBlockHashByHeight Returns the hash of the block currently at :height.
// https://mempool.space/zh/docs/api/rest#get-block-height
func (m *MempoolApiClient) GetBlockHashByHeight(height uint64) (blockHash string, err error) {

	// "https://mempool.space/api/block-height/615615"
	url := fmt.Sprintf("%s/block-height/%v", m.host, height)
	resp, err := m.client.R().Get(url)
	if err != nil {
		return
	}

	if resp.StatusCode() != http.StatusOK {
		err = fmt.Errorf("code:%v,error:%v", resp.StatusCode(), resp.Body())
		return
	}

	blockHash = string(resp.Body())
	return
}

// GetBlockStatus Returns the confirmation status of a block.
// Available fields:
//
// in_best_chain : boolean, false for orphaned blocks,
// next_best : the hash of the next block, only available for blocks in the best chain
//
// https://mempool.space/zh/docs/api/rest#get-block-status
func (m *MempoolApiClient) GetBlockStatus(blockHash string) (isBestChain bool, height uint64, err error) {
	// "https://mempool.space/api/block/0000000000000000000065bda8f8a88f2e1e00d9a6887a43d640e52a4c7660f2/status"
	url := fmt.Sprintf("%s/block/%v/status", m.host, blockHash)
	resp, err := m.client.R().Get(url)
	if err != nil {
		return
	}
	if resp.StatusCode() != http.StatusOK {
		err = fmt.Errorf("code:%v,error:%v", resp.StatusCode(), resp.Body())
		return
	}

	body := string(resp.Body())
	isBestChain = gjson.Get(body, "in_best_chain").Bool()
	height = uint64(gjson.Get(body, "height").Int())
	return
}

// GetTipBlockHeight Returns the height of the last block.
// https://mempool.space/zh/docs/api/rest#get-block-tip-height
func (m *MempoolApiClient) GetTipBlockHeight() (height uint64, err error) {
	url := fmt.Sprintf("%s/blocks/tip/height", m.host)
	resp, err := m.client.R().Get(url)
	if err != nil {
		return
	}
	if resp.StatusCode() != http.StatusOK {
		err = fmt.Errorf("code:%v,error:%v", resp.StatusCode(), resp.Body())
		return
	}

	height, err = strconv.ParseUint(string(resp.Body()), 10, 64)
	return
}

// GetTipBlockHash Returns the hash of the last block.
// https://mempool.space/zh/docs/api/rest#get-block-tip-hash
func (m *MempoolApiClient) GetTipBlockHash() (hash string, err error) {
	url := fmt.Sprintf("%s/blocks/tip/hash", m.host)
	resp, err := m.client.R().Get(url)
	if err != nil {
		return
	}
	if resp.StatusCode() != http.StatusOK {
		err = fmt.Errorf("code:%v,error:%v", resp.StatusCode(), resp.Body())
		return
	}

	hash = string(resp.Body())
	return
}

// GetBlockTansactionIDs Returns a list of all txids in the block.
// https://mempool.space/zh/docs/api/rest#get-block-transaction-ids
func (m *MempoolApiClient) GetBlockTansactionIDs(blockHash string) (txids []string, err error) {
	url := fmt.Sprintf("%s/block/%v/txids", m.host, blockHash)
	resp, err := m.client.R().SetResult([]string{}).Get(url)
	if err != nil {
		return
	}
	if resp.StatusCode() != http.StatusOK {
		err = fmt.Errorf("code:%v,error:%v", resp.StatusCode(), resp.Body())
		return
	}

	txids = resp.Result().([]string)
	return
}

// GetTansaction Returns details about a transaction.
//
//	https://mempool.space/zh/docs/api/rest#get-transaction
func (m *MempoolApiClient) GetTansaction(txid string) (tx Transaction, err error) {

	// https://mempool.space/api/tx/15e10745f15593a899cef391191bdd3d7c12412cc4696b7bcb669d0feadc8521"
	url := fmt.Sprintf("%s/tx/%s", m.host, txid)
	resp, err := m.client.R().SetResult(Transaction{}).Get(url)
	if err != nil {
		return
	}
	if resp.StatusCode() != http.StatusOK {
		err = fmt.Errorf("code:%v,error:%v", resp.StatusCode(), resp.Body())
		return
	}

	tx = resp.Result().(Transaction)

	return
}

// GetAddressDetails Returns details about an address.
//
//	https://mempool.space/zh/docs/api/rest#get-address
func (m *MempoolApiClient) GetAddressDetails(address string) (addressDetails []AddressDetail, err error) {

	//  "https://mempool.space/api/address/1wiz18xYmhRX6xStj2b9t1rwWX4GKUgpv"
	url := fmt.Sprintf("%s/address/%s", m.host, address)
	resp, err := m.client.R().SetResult(Transaction{}).Get(url)
	if err != nil {
		return
	}
	if resp.StatusCode() != http.StatusOK {
		err = fmt.Errorf("code:%v,error:%v", resp.StatusCode(), resp.Body())
		return
	}

	addressDetails = resp.Result().([]AddressDetail)

	return
}

// GetAddressMempool Get unconfirmed transaction history for the specified address/scripthash.
// Returns up to 50 transactions (no paging).
//
//	https://mempool.space/zh/docs/api/rest#get-address-transactions-mempool
func (m *MempoolApiClient) GetAddressMempoolTxs(address string) (mempoolTxs []Transaction, err error) {

	// "https://mempool.space/api/address/1wiz18xYmhRX6xStj2b9t1rwWX4GKUgpv/txs/mempool"
	url := fmt.Sprintf("%s/address/%s/txs/mempool", m.host, address)
	resp, err := m.client.R().SetResult(Transaction{}).Get(url)
	if err != nil {
		return
	}
	if resp.StatusCode() != http.StatusOK {
		err = fmt.Errorf("code:%v,error:%v", resp.StatusCode(), resp.Body())
		return
	}

	mempoolTxs = resp.Result().([]Transaction)

	return
}

// GetRecommendedFees Returns our currently suggested fees for new transactions.
// https://mempool.space/zh/docs/api/rest#get-recommended-fees
func (m *MempoolApiClient) GetRecommendedFees(address string) (recommendedFee RecommendedFee, err error) {

	// "https://mempool.space/api/v1/fees/recommended"
	url := fmt.Sprintf("%s/v1/fees/recommended", m.host)
	resp, err := m.client.R().SetResult(RecommendedFee{}).Get(url)
	if err != nil {
		return
	}
	if resp.StatusCode() != http.StatusOK {
		err = fmt.Errorf("code:%v,error:%v", resp.StatusCode(), resp.Body())
		return
	}

	recommendedFee = resp.Result().(RecommendedFee)

	return
}

// GetAddressUTXOs Get the list of unspent transaction outputs associated with the address/scripthash.
// https://mempool.space/zh/docs/api/rest#get-address-utxo
func (m *MempoolApiClient) GetAddressUTXOs(address string) (utxos []UTXO, err error) {

	// "https://mempool.space/api/address/1KFHE7w8BhaENAswwryaoccDb6qcT6DbYY/utxo"
	url := fmt.Sprintf("%s/address/%s/utxo", m.host, address)
	resp, err := m.client.R().SetResult([]UTXO{}).Get(url)
	if err != nil {
		return
	}
	if resp.StatusCode() != http.StatusOK {
		err = fmt.Errorf("code:%v,error:%v", resp.StatusCode(), resp.Body())
		return
	}

	utxos = resp.Result().([]UTXO)

	return
}
