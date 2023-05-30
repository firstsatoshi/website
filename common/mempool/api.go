package mempool

type MempoolApiClient struct {
	Host string
}

func NewMempoolApiClient(host string) *MempoolApiClient {
	return &MempoolApiClient{
		Host: host,
	}
}

// GetBlockHashByHeight Returns the hash of the block currently at :height.
// https://mempool.space/zh/docs/api/rest#get-block-height
func (m *MempoolApiClient) GetBlockHashByHeight(height uint64) (blockHash string, err error) {

	return "", nil
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


func (m *MempoolApiClient) GetAddressUTXOs(address string) ( utxos []UTXO, err error) {
	return []UTXO{}, nil
}
