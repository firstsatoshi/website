package ordinals

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/firstsatoshi/website/common/btcapi"
	"github.com/firstsatoshi/website/common/mempool"
	"github.com/zeromicro/go-zero/core/logx"
)

// EstimateFee estimates the fee for a given data size and inscription count by
// using  dummy transaction.
func EstimateFee(netParams *chaincfg.Params, feeRate int, noChange bool, inscriptionData []InscriptionData, revealValueSats int64) (fee int64, change int64, err error) {
	if len(inscriptionData) == 0 {
		err = fmt.Errorf("===empty inscriptionData===")
		return
	}

	// initial
	fee = 0
	change = 0
	err = nil
	logx.Infof("start inscribing.........")

	btcApiClient := mempool.NewClient(netParams)

	// DUMMY address
	tmpAddr := "bc1p0ftnthhe6gsthnhd6mswg96aukn888tzrqldz0wkmeeewpr4lkus0vqflq"
	tmpWif := "L2EF51PPPYHkDYVRhyHjzjsNznaHxvy4C8t6GdrDRn5udABQsCYC"
	if netParams.Name == "testnet3" {
		tmpAddr = "tb1p0ftnthhe6gsthnhd6mswg96aukn888tzrqldz0wkmeeewpr4lkuscykx90"
		tmpWif = "cSbEXvPEpbz1Nyxh6P6sN4NSd1shdP4kGB2ZP4JivtjusuFYuTZc"
	}
	wif, err := btcutil.DecodeWIF(tmpWif)
	utxoTaprootAddress, err := btcutil.DecodeAddress(tmpAddr, netParams)

	logx.Infof("start get utxos of address %v", utxoTaprootAddress)

	// make DUMMY UTXOs
	unspentList := make([]*btcapi.UnspentOutput, 0)
	if true {
		txHash, er := chainhash.NewHashFromStr("15e10745f15593a899cef391191bdd3d7c12412cc4696b7bcb669d0feadc8521")
		if er != nil {
			return
		}
		pkScript, er := txscript.PayToAddrScript(utxoTaprootAddress)
		if err != nil {
			return
		}
		unspentList = append(unspentList, &btcapi.UnspentOutput{
			Outpoint: wire.NewOutPoint(txHash, uint32(0)),
			Output:   wire.NewTxOut(10_0000_0000, pkScript),
		})
	}

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
		commitTxPrivateKeyList = append(commitTxPrivateKeyList, wif.PrivKey)
		vinAmount += int(unspentList[i].Output.Value)
	}
	logx.Infof("len(commitTxOutPointList) is %v\n", len(commitTxOutPointList))
	logx.Infof("len(commitTxPrivateKeyList) is %v\n", len(commitTxPrivateKeyList))

	if len(commitTxOutPointList) == 0 || len(commitTxPrivateKeyList) == 0 {
		err = fmt.Errorf("empty commitTxOutPointList or empty commitTxPrivateKeyList")
		return
	}

	request := inscriptionRequest{
		ChangeAddress:          tmpAddr,
		CommitTxOutPointList:   commitTxOutPointList,
		CommitTxPrivateKeyList: commitTxPrivateKeyList,
		CommitFeeRate:          int64(feeRate),
		FeeRate:                int64(feeRate),
		DataList:               inscriptionData,
		RevealOutValue:         revealValueSats,

		// must set false , fix: https://mempool.space/zh/testnet/tx/234e9d7998fcda471f596e5a2d06c311114addb5bbf97e12644cb7e391141523#vin=1
		SingleRevealTxOnly: false,

		// !NOTE: only for estimate fee
		OnlyEstimateFee: true,
		DummyTxOut:      unspentList[0].Output,
		NoChange: true,
	}

	logx.Infof("before newInscriptionToolWithBtcApiClient .........")
	tool, err := newInscriptionToolWithBtcApiClient(netParams, btcApiClient, &request)
	if err != nil {
		logx.Errorf("newInscriptionToolWithBtcApiClient error: %v", err.Error())
		return
	}
	logx.Infof("newInscriptionToolWithBtcApiClient ok")

	change = tool.changeSat
	// fee = tool.calculateFee()
	fee = tool.estimateFee

	logx.Infof("fee: %v", fee)
	// logx.Infof("change: %v", change)
	logx.Info("================================================")

	return
}
