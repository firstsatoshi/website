// Code generated by goctl. DO NOT EDIT.
package types

type JoinWaitListReq struct {
	EventId     int    `json:"eventId"`
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
	CreateTime     int64  `json:"createTime"`
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
	ContentType string `json:"contentType"`
	FileName    string `json:"fileName"`
	Inscription string `json:"inscription"`
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
	PayTime          int64       `json:"paytime"`
	PayConfirmedTime int64       `json:"payConfirmedTime"`
	NftDetails       []NftDetail `json:"nftDetails"`
	CreateTime       int64       `json:"createTime"`
}

type QueryBlindboxEventReq struct {
	EventId       int    `json:"eventId,optional"`
	EventStatus   string `json:"status,optional"` // "all" "active" "inactive"
	EventEndpoint string `json:"eventEndpoint,optional"`
}

type MintPlan struct {
	Title string `json:"title"`
	Plan  string `json:"plan"`
}

type BlindboxEvent struct {
	EventId              int        `json:"eventId"`
	Name                 string     `json:"name"`
	EventEndpoint        string     `json:"eventEndpoint"`
	Description          string     `json:"description"`
	Detail               string     `json:"detail"`
	AvatarImageUrl       string     `json:"avatarImageUrl"` // 头像图片url
	BackgroundImageUrl   string     `json:"backgroundImageUrl"`
	RoadmapDescription   string     `json:"roadmapDescription"`
	RoadmapList          []string   `json:"roadmapList"`
	WebsiteUrl           string     `json:"websiteUrl"`
	TwitterUrl           string     `json:"twitterUrl"`
	DiscordUrl           string     `json:"discordUrl"`
	ImagesList           []string   `json:"imagesList"`
	CurrentMintPlanIndex int        `json:"currentMintPlanIndex"` // 当前mint计划index
	MintPlanList         []MintPlan `json:"mintPlanList"`
	PriceBtcSats         int        `json:"priceBtcSats"`
	PriceUsd             int        `json:"priceUsd"`          // usd价格
	PaymentToken         string     `json:"paymentToken"`      // 支付币种
	AverageImageBytes    int        `json:"averageImageBytes"` // 平均图片大小
	Supply               int        `json:"supply"`            // 总数量
	Avail                int        `json:"avail"`             // 可用数量
	MintLimit            int        `json:"mintLimit"`         // mint 限制
	Active               bool       `json:"isActive"`          // 是否激活
	Display              bool       `json:"isDisplay"`         // 是否显示
	OnlyWhiteist         bool       `json:"onlyWhitelist"`     // 仅白名单
	CustomMint           bool       `json:"customMint"`        // 是否是自定义mint的项目，类似bitfish可以自定义合成
	StartTime            int64      `json:"startTime"`         // 开始时间
	EndTime              int64      `json:"endTime"`           // 结束时间
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
	EventId        int    `json:"eventId"`
	ReceiveAddress string `json:"receiveAddress"`
}

type CheckWhitelistResp struct {
	EventId     int  `json:"eventId"`
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

type FileUpload struct {
	FileName string `json:"fileName"`
	DataUrl  string `json:"dataUrl"` // dataURL schema
}

type CreateInscribeOrderReq struct {
	FileUploads    []FileUpload `json:"fileUploads"`
	ReceiveAddress string       `json:"receiveAddress"`
	FeeRate        int          `json:"feeRate"`
	Checksum       string       `json:"checksum"` // checksum of data and contentType and dataBytes and token
	Token          string       `json:"token"`
}

type CreateInscribeOrderResp struct {
	OrderId        string   `json:"orderId"`
	Count          int      `json:"count"`
	Filenames      []string `json:"filenames"`
	DepositAddress string   `json:"depositAddress"`
	ReceiveAddress string   `json:"receiveAddress"`
	FeeRate        int      `json:"feeRate"`
	Bytes          int      `json:"bytes"`
	InscribeFee    int      `json:"inscribeFee"`
	ServiceFee     int      `json:"serviceFee"`
	Total          int      `json:"total"`
	CreateTime     int64    `json:"createTime"`
}

type QueryBrc20Req struct {
	Ticker string `json:"ticker"`
}

type QueryBrc20Resp struct {
	IsExists      bool   `json:"isExists"`
	Ticker        string `json:"ticker"`
	Limit         int64  `json:"limit"`
	Max           int64  `json:"max"`
	Minted        int64  `json:"minted"`
	Decimal       int    `json:"decimal"`
	InscriptionId string `json:"inscriptionId"`
}

type CheckNamesReq struct {
	Type  string   `json:"type"` // sats, unisat
	Names []string `json:"names"`
}

type CheckNamesResp struct {
	Name     string `json:"name"`
	IsExists bool   `json:"isExists"`
}

type EstimateFeeReq struct {
	FileUploads    []FileUpload `json:"fileUploads"`
	FeeRate        int          `json:"feeRate"`
	ReceiveAddress string       `json:"receiveAddress"`
}

type EstimateFeeResp struct {
	TotalFee int64 `json:"totalFee"`
}

type CheckPathReq struct {
	Path string `json:"path"`
}

type CheckPathResp struct {
	IsExists bool `json:"isExists"`
}
