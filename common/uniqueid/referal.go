package uniqueid

import (
	"fmt"
	"strconv"

	"github.com/btcsuite/btcd/btcutil/base58"
)

func GetReferalCodeById(id int64) string {
	newId := 107 + int(id*3+1)*19
	return base58.Encode([]byte(fmt.Sprint(newId)))
}

func GetIdByReferalCode(referalCode string) int64 {
	newIdStr := string(base58.Decode(referalCode))
	newId, err := strconv.ParseInt(newIdStr, 10, 64)
	if err != nil {
		return 0
	}
	return int64(int64(newId-107)/19-1) / 3
}
