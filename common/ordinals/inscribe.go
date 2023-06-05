package ordinals

import (
	"fmt"

	"github.com/firstsatoshi/website/common/mempool"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

func Inscribe(changeAddress string, wifPrivKey string, netParams *chaincfg.Params, feeRate int,
	inscriptionData []InscriptionData, onlyEstimate bool) (commitTxid string, revealsTxids []string, fee int64, change int64, err error) {

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

	// TODO:  multiple utxo ?
	// collect all of UTXOs
	logx.Infof("utxo size is %v\n", len(unspentList))
	vinAmount := 0
	commitTxOutPointList := make([]*wire.OutPoint, 0)
	commitTxPrivateKeyList := make([]*btcec.PrivateKey, 0)
	for i := range unspentList {
		if unspentList[i].Output.Value < 1000 {
			continue
		}
		commitTxOutPointList = append(commitTxOutPointList, unspentList[i].Outpoint)
		commitTxPrivateKeyList = append(commitTxPrivateKeyList, wifKey.PrivKey)
		vinAmount += int(unspentList[i].Output.Value)
	}
	logx.Infof("len(commitTxOutPointList) is %v\n", len(commitTxOutPointList))
	logx.Infof("len(commitTxPrivateKeyList) is %v\n", len(commitTxPrivateKeyList))

	request := inscriptionRequest{
		ChangeAddress:          changeAddress,
		CommitTxOutPointList:   commitTxOutPointList,
		CommitTxPrivateKeyList: commitTxPrivateKeyList,
		CommitFeeRate:          int64(feeRate),
		FeeRate:                int64(feeRate),
		DataList:               inscriptionData,
		SingleRevealTxOnly:     false,
	}

	tool, err := newInscriptionToolWithBtcApiClient(netParams, btcApiClient, &request)
	if err != nil {
		return
	}

	change = tool.changeSat
	fee = tool.calculateFee()
	if onlyEstimate {
		return
	}

	commitTxHash, revealTxHashList, _, _, err := tool.Inscribe()
	if err != nil {
		err = fmt.Errorf("send tx errr, %v", err)
		return
	}

	commitTxid = commitTxHash.String()
	logx.Info("================================================")
	logx.Infof("commitTxid: %v", commitTxid)

	for i := range revealTxHashList {
		revealsTxids = append(revealsTxids, revealTxHashList[i].String())
		logx.Infof("revealTxid[%v]: %v", i, revealTxHashList[i].String())
	}
	logx.Infof("fee: %v", fee)
	logx.Infof("change: %v", change)
	logx.Info("================================================")

	return
}
