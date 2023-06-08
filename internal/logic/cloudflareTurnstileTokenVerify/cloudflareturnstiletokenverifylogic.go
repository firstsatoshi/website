package cloudflareTurnstileTokenVerify

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/firstsatoshi/website/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type CloudflareTurnstileTokenVerifyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext

	form *url.Values
}

func NewCloudflareTurnstileTokenVerifyLogic(ctx context.Context, svcCtx *svc.ServiceContext, form *url.Values) *CloudflareTurnstileTokenVerifyLogic {
	return &CloudflareTurnstileTokenVerifyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		form:   form,
	}
}

func (l *CloudflareTurnstileTokenVerifyLogic) CloudflareTurnstileTokenVerify() (resp string, err error) {
	logx.Infof("==========CloudflareTurnstileTokenVerify======")

	token := l.form.Get("cf-turnstile-response")
	remoteIp := l.form.Get("CF-Connecting-IP")
	SECRET_KEY := "0x4AAAAAAAFdlF0_97nz6ddK51stJbVThuU"

	logx.Infof("token: %v", token)

	postForm := url.Values{}
	postForm.Set("secret", SECRET_KEY)
	postForm.Set("response", token)
	postForm.Set("remoteip", remoteIp)

	url := "https://challenges.cloudflare.com/turnstile/v0/siteverify"
	rsp, err := http.PostForm(url, postForm)

	if err != nil {
		logx.Errorf("error: %v", err.Error())
		return "", fmt.Errorf("%v", err.Error())
	}

	if rsp.StatusCode != http.StatusOK {
		logx.Errorf("statusCode: %v", rsp.StatusCode)
		return "", fmt.Errorf("statusCode %v", rsp.StatusCode)
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		logx.Errorf("error: %v", err.Error())
		return "", fmt.Errorf("errro: %v", err.Error())
	}

	type Resp struct {
		Success bool `json:"success"`
	}
	var r Resp
	if err = json.Unmarshal(body, &r); err != nil {
		logx.Errorf("json Unmarshal error: %v", err.Error())
		return "", fmt.Errorf("json Unmarshal error: %v", err.Error())
	}

	if !r.Success {
		logx.Errorf("valid failed  ")
		return "", fmt.Errorf("error: %v", string(body))
	}

	// set expire
	l.svcCtx.Redis.SetexCtx(l.ctx, token, "ok", 300)

	return "valid ok!\n" + string(body), nil
}
