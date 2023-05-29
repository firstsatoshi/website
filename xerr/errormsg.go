package xerr

var message map[uint32]string

func init() {
	message = make(map[uint32]string)
	message[OK] = "ok"
	message[SERVER_COMMON_ERROR] = "system error, try it later"
	message[REUQEST_PARAM_ERROR] = "invalid request parameter error, please check it"
	message[INVALID_EMAIL_ERROR] = "invalid email"
	message[TOO_MANY_REQUEST_ERROR] = "request too frequently"
	message[INVALID_BTCP2TRADDRESS_ERROR] = "invalid bitcoin P2TR(Taproot) inscription receive address, it must be prefix with 'bc1p'"
	message[FEERATE_TOO_SMALL_ERROR] = "fee rate too small"
	message[COUNT_EXCEED_PER_ORDER_LIMIT_ERROR] = "inscribe count exceed per order limit"
	message[EVENT_NOT_EXISTS_ERROR] = "event does not exists"
	message[AVAILABLE_COUNT_IS_NOT_ENOUGH] = "avialable count is not enough"
	message[DB_ERROR] = "db error"
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
