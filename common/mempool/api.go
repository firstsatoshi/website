package mempool

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type MempoolApiClient struct {
	host   string
	client *resty.Client
}

func NewMempoolApiClient(host string) *MempoolApiClient {

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
	return false, 0, nil
}

// GetTipBlockHeight Returns the height of the last block.
// https://mempool.space/zh/docs/api/rest#get-block-tip-height
func (m *MempoolApiClient) GetTipBlockHeight() (height uint64, err error) {
	return 0, nil
}

// GetTipBlockHash Returns the hash of the last block.
// https://mempool.space/zh/docs/api/rest#get-block-tip-hash
func (m *MempoolApiClient) GetTipBlockHash() (hash string, err error) {
	return "", nil
}

// GetBlockTansactionIDs Returns a list of all txids in the block.
// https://mempool.space/zh/docs/api/rest#get-block-transaction-ids
func (m *MempoolApiClient) GetBlockTansactionIDs(blackHash string) (txids []string, err error) {
	return []string{}, nil
}

// GetTansaction Returns details about a transaction.
//
//	https://mempool.space/zh/docs/api/rest#get-transaction
func (m *MempoolApiClient) GetTansaction(txid string) (txs []Transaction, err error) {
	return []Transaction{}, nil
}

// GetAddressDetails Returns details about an address.
//
//	https://mempool.space/zh/docs/api/rest#get-address
func (m *MempoolApiClient) GetAddressDetails(address string) (addressDetails []AddressDetail, err error) {
	return []AddressDetail{}, nil
}

// GetAddressMempool Get unconfirmed transaction history for the specified address/scripthash.
// Returns up to 50 transactions (no paging).
//
//	https://mempool.space/zh/docs/api/rest#get-address-transactions-mempool
func (m *MempoolApiClient) GetAddressMempool(address string) (mempoolTxs []Transaction, err error) {
	return []Transaction{}, nil
}

// GetRecommendedFees Returns our currently suggested fees for new transactions.
// https://mempool.space/zh/docs/api/rest#get-recommended-fees
func (m *MempoolApiClient) GetRecommendedFees(address string) (recommendedFee RecommendedFee, err error) {
	return RecommendedFee{}, nil
}

// GetAddressUTXOs Get the list of unspent transaction outputs associated with the address/scripthash.
// https://mempool.space/zh/docs/api/rest#get-address-utxo
func (m *MempoolApiClient) GetAddressUTXOs(address string) (utxos []UTXO, err error) {
	return []UTXO{}, nil
}
