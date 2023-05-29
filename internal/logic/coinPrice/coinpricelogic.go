package coinPrice

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/fantopia-dev/website/internal/svc"
	"github.com/fantopia-dev/website/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CoinPriceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCoinPriceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CoinPriceLogic {
	return &CoinPriceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CoinPriceLogic) CoinPrice() (*types.CoinPriceResp, error) {
	key := "bitcion-price-usd"
	timeout := 3600

	v, err := l.svcCtx.Redis.Get(key)
	if err == nil {
		price, err := strconv.ParseFloat(v, 64)
		if err == nil {
			return &types.CoinPriceResp{
				BtcPriceUsd: float64(int32(price)),
			}, nil
		}
	}

	for i := 0; i < 5; i++ {

		type DataItem struct {
			PriceUsd string `json:"priceUsd"`
		}
		type Resp struct {
			Data DataItem `json:"data"`
		}

		rsp, err := http.Get("https://api.coincap.io/v2/assets/bitcoin")
		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		if rsp.StatusCode != http.StatusOK {
			time.Sleep(time.Second)
			continue
		}

		body, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		var r Resp
		if err = json.Unmarshal(body, &r); err != nil {
			time.Sleep(time.Second)
			continue
		}

		l.svcCtx.Redis.Setex(key, r.Data.PriceUsd, timeout)

		price, err := strconv.ParseFloat(v, 64)
		if err == nil {
			return &types.CoinPriceResp{
				BtcPriceUsd: float64(int32(price)),
			}, nil
		}
	}

	return nil, fmt.Errorf("coincap api error")
}
