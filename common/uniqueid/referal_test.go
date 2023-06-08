package uniqueid_test

import (
	"testing"

	"github.com/firstsatoshi/website/common/uniqueid"
)

func TestGetReferalCodeById(t *testing.T) {

	id := int64(1000)
	code := uniqueid.GetReferalCodeById(id)
	t.Logf("code : %v", code)
	rid := uniqueid.GetIdByReferalCode(code)
	if rid != id {
		t.FailNow()
	}
	t.Logf("id : %v", rid)

}
