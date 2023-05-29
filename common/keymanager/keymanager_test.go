package keymanager_test

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/fantopia-dev/website/common/keymanager"
)

// func TestCreateWallet(t *testing.T) {

// 	km, err := keymanager.NewKeyManagerFromSeed("xiexieniyoungqqcn@163.com")
// 	if err != nil {
// 		t.Fatalf("error: %v", err.Error())
// 	}

// 	addressIndex := uint32(1)
// 	k, err := km.GetKey(keymanager.PurposeBIP44, keymanager.CoinTypeBTC, 0, 0, addressIndex)
// 	if err != nil {
// 		t.Fatalf("error: %v", err.Error())
// 	}

// 	t.Logf("path : %v\n", k.GetPath())

// 	compressed := true
// 	wif, address, taprootBech32, err := k.GetWifKeyAndAddress(compressed, chaincfg.MainNetParams)
// 	if err != nil {
// 		t.Fatalf("error: %v", err.Error())
// 	}

// 	t.Logf("WIF: %v\n", wif)
// 	t.Logf("address: %v\n", address)
// 	t.Logf("taprootBech32: %v\n", taprootBech32)
// }

func TestGetWifKeyAndAddresss(t *testing.T) {

	km, err := keymanager.NewKeyManagerFromSeed("xiexieniyoungqqcn@163.com", chaincfg.MainNetParams)
	if err != nil {
		t.Fatalf("error: %v", err.Error())
	}

	wif, taprootBech32, err := km.GetWifKeyAndAddresss(0, 1, chaincfg.MainNetParams)
	if err != nil {
		t.Fatalf("error: %v", err.Error())
	}

	t.Logf("WIF: %v\n", wif)
	// t.Logf("address: %v\n", address)
	t.Logf("taprootBech32: %v\n", taprootBech32)
}
