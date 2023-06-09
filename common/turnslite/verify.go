package turnslite

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/firstsatoshi/website/common/globalvar"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

func VeifyToken(ctx context.Context, token string, rds *redis.Redis) (bool, error) {

	SECRET_KEY := "0x4AAAAAAAFdlF0_97nz6ddK51stJbVThuU"

	if len(token) < 10 {
		return false, fmt.Errorf("invalid token")
	}

	logx.Infof("token: %v", token)
	h := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(h[:])
	v, err := rds.GetCtx(ctx, fmt.Sprintf("%v:%v", globalvar.TURNSTILE_TOKEN_PREFIX, string(tokenHash[:])))
	if err == nil && v == token {
		logx.Infof("token ooooooooooooooooooooooooooook")
		return true, nil
	}

	postForm := url.Values{}
	postForm.Set("secret", SECRET_KEY)
	postForm.Set("response", token)

	url := "https://challenges.cloudflare.com/turnstile/v0/siteverify"
	rsp, err := http.PostForm(url, postForm)

	if err != nil {
		logx.Errorf("error: %v", err.Error())
		return false, fmt.Errorf("%v", err.Error())
	}

	if rsp.StatusCode != http.StatusOK {
		logx.Errorf("statusCode: %v", rsp.StatusCode)
		return false, fmt.Errorf("statusCode %v", rsp.StatusCode)
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		logx.Errorf("error: %v", err.Error())
		return false, fmt.Errorf("errro: %v", err.Error())
	}

	type Resp struct {
		Success bool `json:"success"`
	}
	var r Resp
	if err = json.Unmarshal(body, &r); err != nil {
		logx.Errorf("json Unmarshal error: %v", err.Error())
		return false, fmt.Errorf("json Unmarshal error: %v", err.Error())
	}

	if !r.Success {
		logx.Errorf("valid failed  ")
		return false, fmt.Errorf("error: %v", string(body))
	}

	// set expire
	rds.SetexCtx(ctx, fmt.Sprintf("%v:%v", globalvar.TURNSTILE_TOKEN_PREFIX, tokenHash), token, 300)

	logx.Infof("===========turnslite token verify ok ===================")
	return true, nil
}
