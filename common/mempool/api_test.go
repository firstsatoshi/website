package mempool_test

import (
	"testing"

	"github.com/firstsatoshi/website/common/mempool"
)

func TestGetBlockHashByHeight(t *testing.T) {

	client := mempool.NewMempoolApiClient("")

	hash, err := client.GetBlockHashByHeight(615615)
	if err != nil {
		t.Fatalf("error:%v", err.Error())
	}

	if hash != "000000000000000000067bea442af50a91377ac796e63b8d284354feff4042b3" {
		t.Fatalf("expected 000000000000000000067bea442af50a91377ac796e63b8d284354feff4042b3")
	}
}

func TestGetBlockStatus(t *testing.T) {

	client := mempool.NewMempoolApiClient("")

	isBestChain, height, err := client.GetBlockStatus("0000000000000000000065bda8f8a88f2e1e00d9a6887a43d640e52a4c7660f2")
	if err != nil {
		t.Fatalf("error:%v", err.Error())
	}

	if isBestChain != true {
		t.Fatalf("expected:%v", true)
	}

	if height != 690557 {
		t.Fatalf("expected:%v", 690557)
	}

}

func TestGetTipBlockHeight(t *testing.T) {

	client := mempool.NewMempoolApiClient("")

	height, err := client.GetTipBlockHeight()
	if err != nil {
		t.Fatalf("error:%v", err.Error())
	}

	if height < 1000 {
		t.Fatalf("expected: positive block heigh")
	}

}

func TestGetTipBlockHash(t *testing.T) {

	client := mempool.NewMempoolApiClient("")

	hash, err := client.GetTipBlockHash()
	if err != nil {
		t.Fatalf("error:%v", err.Error())
	}

	if len(hash) != 64 {
		t.Fatalf("expected 64 hash %v", hash)
	}

}

func TestGetBlockTansactionIDs(t *testing.T) {

	client := mempool.NewMempoolApiClient("")

	txids, err := client.GetBlockTansactionIDs("000000000000000015dc777b3ff2611091336355d3f0ee9766a2cf3be8e4b1ce")
	if err != nil {
		t.Fatalf("error:%v", err.Error())
	}

	if len(txids) == 0 {
		t.Fatalf("expected non-empty txids ")
	}
}

func TestGetTansaction(t *testing.T) {

	client := mempool.NewMempoolApiClient("")

	txid := "15e10745f15593a899cef391191bdd3d7c12412cc4696b7bcb669d0feadc8521"
	tx, err := client.GetTansaction(txid)
	if err != nil {
		t.Fatalf("error:%v", err.Error())
	}

	if tx.Txid != txid {
		t.Fatalf("txid not matched")
	}

	if tx.Vout[0].ScriptpubkeyAddress != "1BUBQuPV3gEV7P2XLNuAJQjf5t265Yyj9t" {
		t.FailNow()
	}

	if tx.Vout[3].ScriptpubkeyAddress != "1wizSAYSbuyXbt9d8JV8ytm5acqq2TorC" {
		t.FailNow()
	}

	if tx.Vin[0].Txid != "1fdfed84588cb826b876cd761ecebcf1726453437f0a6826e82ed54b2807a036" {
		t.FailNow()
	}

	if tx.Vin[4].Txid != "1fdfed84588cb826b876cd761ecebcf1726453437f0a6826e82ed54b2807a036" {
		t.FailNow()
	}

}

func TestGetAddressDetails(t *testing.T) {

	client := mempool.NewMempoolApiClient("")

	addrDetails, err := client.GetAddressDetails("1wiz18xYmhRX6xStj2b9t1rwWX4GKUgpv")
	if err != nil {
		t.Fatalf("error:%v", err.Error())
	}

	if addrDetails.Address != "1wiz18xYmhRX6xStj2b9t1rwWX4GKUgpv" {
		t.FailNow()
	}

	if addrDetails.ChainStats.FundedTxoSum < 100 {
		t.FailNow()
	}

}

func TestGetAddressMempoolTxs(t *testing.T) {

	client := mempool.NewMempoolApiClient("")

	_, err := client.GetAddressMempoolTxs("1wiz18xYmhRX6xStj2b9t1rwWX4GKUgpv")
	if err != nil {
		t.Fatalf("error:%v", err.Error())
	}

}

func TestGetRecommendedFees(t *testing.T) {

	client := mempool.NewMempoolApiClient("")

	fee, err := client.GetRecommendedFees()
	if err != nil {
		t.Fatalf("error:%v", err.Error())
	}

	if fee.EconomyFee == 0 || fee.FastestFee == 0 {
		t.FailNow()
	}

}


func TestGetAddressUTXOs(t *testing.T) {

	client := mempool.NewMempoolApiClient("")

	utxos, err := client.GetAddressUTXOs("1KFHE7w8BhaENAswwryaoccDb6qcT6DbYY")
	if err != nil {
		t.Fatalf("error:%v", err.Error())
	}


	if len(utxos) == 0 {
		t.FailNow()
	}

	if utxos[0].Value == 0 {
		t.FailNow()
	}

}





