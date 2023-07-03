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
