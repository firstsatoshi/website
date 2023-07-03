package unisat_test

import (
	"testing"

	"github.com/firstsatoshi/website/common/unisat"
)

func TestUnisatApiClient_GetBrc20Info(t *testing.T) {

	c := unisat.NewUnisatApiClient()
	{

		brc20Info, err := c.GetBrc20Info("ordi")
		if err != nil {
			t.Fatalf("err:%v", err)
		}

		t.Logf("brc20Info:%+v", brc20Info)
	}

	{
		brc20Info, err := c.GetBrc20Info("zfsl")
		// brc20Info, err := c.GetBrc20Info("fsat")
		if err != nil {
			t.Fatalf("err:%v", err)
		}
		t.Logf("brc20Info:%+v", brc20Info)
	}

}

func TestUnisatApiClient_CheckNames(t *testing.T) {
	c := unisat.NewUnisatApiClient()

	r, err := c.CheckNames("sats", []string{"aaa.sats"})
	if err != nil {
		t.Fatalf("err:%v", err)
	}

	t.Logf("r:%+v", r)
}
