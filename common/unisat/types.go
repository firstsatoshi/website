package unisat

type Brc20Info struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`  // if ticker is not exists , code is -1
	Data struct {
		Ticker        string `json:"ticker"`
		Limit         string `json:"limit"`
		Max           string  `json:"max"`
		Minted        string `json:"minted"`
		Decimal       int    `json:"decimal"`
		InscriptionId string `json:"inscriptionId"`
	} `json:"data"`
}
