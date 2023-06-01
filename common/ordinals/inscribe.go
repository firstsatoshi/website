package ordinals

import (
	"fmt"
	"io/ioutil"

	"github.com/firstsatoshi/website/common/mempool"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

func Inscribe(wifPrivKey string, feeRate int, inscriptionData []InscriptionData, onlyEstimate bool) (txid string, txids []string, fee int64, err error) {
	// netParams := &chaincfg.MainNetParams

	netParams := &chaincfg.RegressionNetParams
	btcApiClient := mempool.NewClient(netParams)
	wifKey, err := btcutil.DecodeWIF(wifPrivKey)
	if err != nil {
		return
	}
	utxoTaprootAddress, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(txscript.ComputeTaprootKeyNoScript(wifKey.PrivKey.PubKey())), netParams)
	if err != nil {
		return
	}
	unspentList, err := btcApiClient.ListUnspent(utxoTaprootAddress)
	// unspentList, err := ListUnspent(utxoTaprootAddress)
	// if err != nil {
	// 	return
	// }

	fmt.Printf("utxo size is %v\n", len(unspentList))

	// return

	// if len(unspentList) == 0 {
	// 	err = fmt.Errorf("no utxo for %s", utxoTaprootAddress)
	// 	return
	// }

	vinAmount := 0
	commitTxOutPointList := make([]*wire.OutPoint, 0)
	commitTxPrivateKeyList := make([]*btcec.PrivateKey, 0)
	for i := range unspentList {
		if unspentList[i].Output.Value < 10000 {
			continue
		}
		commitTxOutPointList = append(commitTxOutPointList, unspentList[i].Outpoint)
		commitTxPrivateKeyList = append(commitTxPrivateKeyList, wifKey.PrivKey)
		vinAmount += int(unspentList[i].Output.Value)
	}
	fmt.Printf("len(commitTxOutPointList) is %v\n", len(commitTxOutPointList))
	fmt.Printf("len(commitTxPrivateKeyList) is %v\n", len(commitTxPrivateKeyList))

	dataList := make([]InscriptionData, 0)

	// read image from filename
	imgBs, err := ioutil.ReadFile("../eagle-1.png")
	if err != nil {
		fmt.Printf("error:%v\n", err.Error())
		return
	}
	mint := InscriptionData{
		// ContentType: "image/jpeg",
		// ContentType: "image/gif",
		ContentType: "image/png",
		Body:        imgBs,
		Destination: utxoTaprootAddress.EncodeAddress(),
	}

	// mint := ord.InscriptionData{
	// 	ContentType: "text/plain;charset=utf-8",
	// 	Body:        []byte(fmt.Sprintf(`{"p":"brc-20","op":"%s","tick":"%s","amt":"%s"}`, gop, gtick, gamount)),
	// 	Destination: utxoTaprootAddress.EncodeAddress(),
	// }

	count := len(inscriptionData)

	for i := 0; i < count; i++ {
		dataList = append(dataList, mint)
	}

	request := inscriptionRequest{
		CommitTxOutPointList:   commitTxOutPointList,
		CommitTxPrivateKeyList: commitTxPrivateKeyList,
		CommitFeeRate:          int64(feeRate),
		FeeRate:                int64(feeRate),
		DataList:               dataList,
		SingleRevealTxOnly:     false,
	}

	tool, err := newInscriptionToolWithBtcApiClient(netParams, btcApiClient, &request)
	// tool, err := ord.NewInscriptionTool(netParams, client, &request)
	if err != nil {
		return
	}

	baseFee := tool.calculateFee()

	if onlyEstimate {
		fee = baseFee
		return
	}

	commitTxHash, revealTxHashList, _, _, err := tool.Inscribe()
	if err != nil {
		err = fmt.Errorf("send tx errr, %v", err)
		return
	}

	txid = commitTxHash.String()
	fmt.Println(txid)
	for i := range revealTxHashList {
		txids = append(txids, revealTxHashList[i].String())
		fmt.Println(revealTxHashList[i].String())
	}
	fmt.Printf("fee: %v\n", fee)

	return
}
