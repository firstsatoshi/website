package ordinals_test

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/firstsatoshi/website/common/ordinals"
)

func TestEstimateFee(t *testing.T) {

	noChange := true
	fee, _, err := ordinals.EstimateFee(&chaincfg.MainNetParams, 3, noChange, []ordinals.InscriptionData{
		{
			ContentType: "text/plain;charset=UTF-8",
			// Body:        []byte("Hello, world!"),
			Body:        []byte("fsa"),
			Destination: "bc1p0ftnthhe6gsthnhd6mswg96aukn888tzrqldz0wkmeeewpr4lkus0vqflq",
		},
	}, 546)

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("fee %v", fee)
}
