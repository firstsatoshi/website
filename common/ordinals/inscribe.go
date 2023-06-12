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
	inscriptionData []InscriptionData, revealValueSats int64, onlyEstimate bool) (
	commitTxid string, revealsTxids []string, fee int64, change int64,
	orderBroadcastAtom *OrderBroadcastAtom, err error,
) {

	// initial
	commitTxid = ""
	revealsTxids = make([]string, 0)
	fee = 0
	change = 0
	err = nil
	logx.Infof("start inscribing.........")

	btcApiClient := mempool.NewClient(netParams)
	wifKey, err := btcutil.DecodeWIF(wifPrivKey)
	if err != nil {
		return
	}
	utxoTaprootAddress, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(txscript.ComputeTaprootKeyNoScript(wifKey.PrivKey.PubKey())), netParams)
	if err != nil {
		return
	}

	logx.Infof("start get utxos of address %v", utxoTaprootAddress)
	unspentList, err := btcApiClient.ListUnspent(utxoTaprootAddress)
	if len(unspentList) == 0 {
		err = fmt.Errorf("empty utxos")
		return
	}

	// TODO:  multiple utxo ?
	// collect all of UTXOs
	logx.Infof("%v utxo size is %v\n", utxoTaprootAddress, len(unspentList))
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

	if len(commitTxOutPointList) == 0 || len(commitTxPrivateKeyList) == 0 {
		err = fmt.Errorf("empty commitTxOutPointList or empty commitTxPrivateKeyList")
		return
	}

	request := inscriptionRequest{
		ChangeAddress:          changeAddress,
		CommitTxOutPointList:   commitTxOutPointList,
		CommitTxPrivateKeyList: commitTxPrivateKeyList,
		CommitFeeRate:          int64(feeRate),
		FeeRate:                int64(feeRate),
		DataList:               inscriptionData,
		RevealOutValue:         revealValueSats,

		// must set false , fix: https://mempool.space/zh/testnet/tx/234e9d7998fcda471f596e5a2d06c311114addb5bbf97e12644cb7e391141523#vin=1
		SingleRevealTxOnly: false,
	}

	logx.Infof("before newInscriptionToolWithBtcApiClient .........")
	tool, err := newInscriptionToolWithBtcApiClient(netParams, btcApiClient, &request)
	if err != nil {
		return
	}
	logx.Infof("newInscriptionToolWithBtcApiClient ok")

	change = tool.changeSat
	fee = tool.calculateFee()
	if onlyEstimate {
		return
	}

	logx.Infof("before tool.Inscribe .....")

	commitTxHash, revealTxHashList, _, _, orderBroadcastAtom, err := tool.Inscribe()
	if err != nil {
		logx.Errorf("inscribe error: %v", err.Error())
		err = fmt.Errorf("send tx errr, %v", err)
		return
	}

	commitTxid = commitTxHash.String()
	logx.Info("==================Inscribe ok==============================")
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
