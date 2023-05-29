package address_test

import (
	"testing"

	"github.com/fantopia-dev/website/common/address"
)

func TestCreateWallet(t *testing.T) {

	km, err := address.NewKeyManagerFromSeed("xiexieniyoungqqcn@163.com")
	if err != nil {
		t.Fatalf("error: %v", err.Error())
	}

	addressIndex := uint32(1)
	k, err := km.GetKey(address.PurposeBIP44, address.CoinTypeBTC, 0, 0, addressIndex)
	if err != nil {
		t.Fatalf("error: %v", err.Error())
	}

	t.Logf("path : %v\n", k.GetPath())

	compressed := true
	wif, address, taprootBech32, err := k.GetWifKeyAndAddress(compressed)
	if err != nil {
		t.Fatalf("error: %v", err.Error())
	}

	t.Logf("WIF: %v\n", wif)
	t.Logf("address: %v\n", address)
	t.Logf("taprootBech32: %v\n", taprootBech32)
}
