package xerr

var message map[uint32]string

func init() {
	message = make(map[uint32]string)
	message[OK] = "ok"
	message[SERVER_COMMON_ERROR] = "system error, try it later"
	message[REUQEST_PARAM_ERROR] = "parameter error"
	message[INVALID_EMAIL_ERROR] = "invalid email"
	message[TOO_MANY_REQUEST_ERROR] = "request too frequently"
	message[INVALID_BTCP2TRADDRESS_ERROR] = "invalid bitcoin P2TR address, the address must be prefix with bc1p"
}

func MapErrMsg(errcode uint32) string {
	if msg, ok := message[errcode]; ok {
		return msg
	} else {
		return "unknow error"
	}
}

func IsCodeErr(errcode uint32) bool {
	if _, ok := message[errcode]; ok {
		return true
	} else {
		return false
	}
}
