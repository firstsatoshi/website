package response

import (
	"net/http"

	"github.com/firstsatoshi/website/xerr"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type Body struct {
	Code uint32      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Response(w http.ResponseWriter, resp interface{}, err error) {
	var body Body
	if err != nil {

		//错误返回
		errcode := xerr.SERVER_COMMON_ERROR
		errmsg := "服务器开小差啦，稍后再来试一试"

		causeErr := errors.Cause(err)                // err类型
		if e, ok := causeErr.(*xerr.CodeError); ok { //自定义错误类型
			//自定义CodeError
			errcode = e.GetErrCode()
			errmsg = e.GetErrMsg()
		}
		body.Code = errcode
		body.Msg = errmsg
		body.Data = nil

		logx.Errorf("%v", err)

	} else {
		body.Code = xerr.OK // 0
		body.Msg = "ok"
		body.Data = resp
	}
	httpx.OkJson(w, body)
}
