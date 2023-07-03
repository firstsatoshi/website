package unisat

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

type UnisatApiClient struct {
	host   string
	client *resty.Client
}

func NewUnisatApiClient() *UnisatApiClient {
	return &UnisatApiClient{
		client: resty.New(),
	}
}

// get brc20 info by unsat api
func (c *UnisatApiClient) GetBrc20Info(ticker string) (brc20Info *Brc20Info, err error) {

	url := fmt.Sprintf("https://unisat.io/brc20-api-v2/brc20/%v/info", ticker)
	resp, err := c.client.R().SetResult(&Brc20Info{}).Get(url)
	if err != nil {
		return
	}

	if resp.StatusCode() != http.StatusOK {
		err = fmt.Errorf("code:%v,error:%v", resp.StatusCode(), string(resp.Body()))
		return
	}

	info := resp.Result().(*Brc20Info)
	if info == nil {
		err = fmt.Errorf("code:%v,error:%v", resp.StatusCode(), string(resp.Body()))
		return
	}

	brc20Info = info
	return
}

// get names
func (c *UnisatApiClient) CheckNames(nameType string, names []string) (results map[string]bool, err error) {
	url := fmt.Sprintf("https://unisat.io/brc20-api-v2/inscriptions/category/%s/existence", nameType)
	namesJson, err := json.Marshal(map[string][]string{
		"names": names,
	})
	if err != nil {
		return nil, err
	}

	resp, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(namesJson).Post(url)
	if err != nil {
		return
	}

	if resp.StatusCode() != http.StatusOK {
		err = fmt.Errorf("code:%v,error:%v", resp.StatusCode(), string(resp.Body()))
		return
	}

	body := string(resp.Body())
	fmt.Printf("==%v", body)
	result := gjson.Get(body, "data").Map()

	results = make(map[string]bool, 0)
	for k, v := range result {
		results[k] = v.String() == "confirmed"
	}
	return
}
