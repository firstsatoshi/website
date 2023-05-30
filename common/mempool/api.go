package mempool

// curl -sSL "https://mempool.space/api/blocks/tip/hash"

type MempoolApiClient struct {
	Host string
}

func NewMempoolApiClient(host string) *MempoolApiClient {
	return &MempoolApiClient{
		Host: host,
	}
}

// curl -sSL "https://mempool.space/api/block-height/615615"
func (m *MempoolApiClient) GetBlockHashByHeight(height uint32) (blockHash string, err error) {
	// 000000000000000000067bea442af50a91377ac796e63b8d284354feff4042b3

	return "", nil
}

// curl -sSL "https://mempool.space/api/block/0000000000000000000065bda8f8a88f2e1e00d9a6887a43d640e52a4c7660f2/status"
func (m *MempoolApiClient) GetBlockStatus(blockHash string) (isBestChain bool, height uint32, err error) {

	/*
		{
			in_best_chain: true,
			height: 690557,
			next_best: "00000000000000000003a59a34c93e39e636c8cd23ead726fdc467fbed0b7c5a"
		}
	*/
	return false, 0, nil
}

// curl -sSL "https://mempool.space/api/blocks/tip/height"
func (m *MempoolApiClient) GetTipBlockHeight() (height uint32, err error) {
	return 0, nil
}

// curl -sSL "https://mempool.space/api/blocks/tip/hash"
func (m *MempoolApiClient) GetTipBlockHash() (hash string, err error) {
	return "", nil
}

// curl -sSL "https://mempool.space/api/block/000000000000000015dc777b3ff2611091336355d3f0ee9766a2cf3be8e4b1ce/txids"
func (m *MempoolApiClient) GetBlockTansactionIDs(blackHash string) (txids []string, err error) {
	return []string{}, nil
}

// curl -sSL "https://mempool.space/api/tx/15e10745f15593a899cef391191bdd3d7c12412cc4696b7bcb669d0feadc8521"
func (m *MempoolApiClient) GetTansaction(txid string) (txs []Transaction, err error) {
	return []Transaction{}, nil
}
