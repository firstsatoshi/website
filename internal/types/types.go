// Code generated by goctl. DO NOT EDIT.
package types

type JoinWaitListReq struct {
	Email       string `json:"email"`
	BtcAddress  string `json:"btcAddress"`
	ReferalCode string `json:"referalCode,optional"`
	Token       string `json:"token"`
}

type JoinWaitListResp struct {
	ReferalCode string `json:"referalCode"`
	Duplicated  bool   `json:"duplicated"`
}

type CreateOrderReq struct {
	EventId        int    `json:"eventId"`
	Count          int    `json:"count"`
	ReceiveAddress string `json:"receiveAddress"`
	FeeRate        int    `json:"feeRate"`
	Token          string `json:"token"`
}

type CreateOrderResp struct {
	OrderId        string `json:"orderId"`
	EventId        int    `json:"eventId"`
	Count          int    `json:"count"`
	DepositAddress string `json:"depositAddress"`
	ReceiveAddress string `json:"receiveAddress"`
	FeeRate        int    `json:"feeRate"`
	Bytes          int    `json:"bytes"`
	InscribeFee    int    `json:"inscribeFee"`
	ServiceFee     int    `json:"serviceFee"`
	Price          int    `json:"price"`
	Total          int    `json:"total"`
	CreateTime     string `json:"createTime"`
}

type QueryOrderReq struct {
	OrderId        string `json:"orderId,optional"`
	ReceiveAddress string `json:"receiveAddress,optional"`
	DepositAddress string `json:"depositAddress,optional"`
}

type NftDetail struct {
	Txid        string `json:"txid"`
	TxConfirmed bool   `json:"confirmed"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	ImageUrl    string `json:"imageUrl"`
}

type Order struct {
	OrderId          string      `json:"orderId"`
	EventId          int         `json:"eventId"`
	DepositAddress   string      `json:"depositAddress"`
	Count            int         `json:"count"`
	Total            int         `json:"total"`
	FeeRate          int         `json:"feeRate"`
	ReceiveAddress   string      `json:"receiveAddress"`
	OrderStatus      string      `json:"orderStatus"`
	PayTxid          string      `json:"payTxid"`
	PayTime          string      `json:"paytime"`
	PayConfirmedTime string      `json:"payConfirmedTime"`
	NftDetails       []NftDetail `json:"nftDetails"`
	CreateTime       string      `json:"createTime"`
}

type BlindboxEvent struct {
	EventId           int    `json:"eventId"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	PriceBtcSats      int    `json:"priceBtcSats"`
	PriceUsd          int    `json:"priceUsd"`          // usd价格
	PaymentToken      string `json:"paymentToken"`      // 支付币种
	AverageImageBytes int    `json:"averageImageBytes"` // 平均图片大小
	Supply            int    `json:"supply"`            // 总数量
	Avail             int    `json:"avail"`             // 可用数量
	MintLimit         int    `json:"mintLimit"`         // mint 限制
	Enable            bool   `json:"enable"`            // 是否开启
	OnlyWhiteist      bool   `json:"onlyWhitelist"`     // 仅白名单
	StartTime         string `json:"startTime"`         // 开始时间
	EndTime           string `json:"endTime"`           // 结束时间
}

type CoinPriceResp struct {
	BtcPriceUsd float64 `json:"btcPriceUsd"`
}

type QueryGalleryListReq struct {
	CurPage  int    `json:"curPage"`
	PageSize int    `json:"pageSize"`
	Category string `json:"category"`
}

type NFT struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageUrl    string `json:"imageUrl"`
}

type QueryGalleryListResp struct {
	Category string `json:"category"`
	CurPage  int    `json:"curPage"`
	Total    int    `json:"totalPage"`
	PageSize int    `json:"pageSize"`
	NFTs     []NFT  `json:"nfts"`
}

type CheckWhitelistReq struct {
	ReceiveAddress string `json:"receiveAddress"`
}

type CheckWhitelistResp struct {
	IsWhitelist bool `json:"isWhitelist"`
}

type QueryAddressReq struct {
	EventId        int    `json:"eventId"`
	ReceiveAddress string `json:"receiveAddress"`
}

type QueryAddressResp struct {
	EventId            int  `json:"eventId"`
	IsWhitelist        bool `json:"isWhitelist"`
	EventMintLimit     int  `json:"eventMintLimit"`
	CurrentOrdersTotal int  `json:"currentOrdersTotal"`
}
