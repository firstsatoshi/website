package xerr

var message map[uint32]string

func init() {
	message = make(map[uint32]string)
	message[OK] = "ok"
	message[SERVER_COMMON_ERROR] = "System error, please try again later"
	message[REUQEST_PARAM_ERROR] = "Invalid request parameter error, please check it"
	message[INVALID_EMAIL_ERROR] = "Invalid email"
	message[TOO_MANY_REQUEST_ERROR] = "Too many requests, please try it later"
	message[INVALID_BTCP2TRADDRESS_ERROR] = "Invalid Bitcoin Taproot(P2TR) inscription receive address, it must starts with 'bc1p'"
	message[FEERATE_TOO_SMALL_ERROR] = "Fee rate is not safe, too small or too large, please set a proper fee rate and try again"
	message[COUNT_EXCEED_PER_ORDER_LIMIT_ERROR] = "Sorry mint count exceed order limit 50, please reduce count and try again"
	message[EVENT_NOT_EXISTS_ERROR] = "event does not exists"
	message[AVAILABLE_COUNT_IS_NOT_ENOUGH] = "Sorry, current avialable count is not enough."
	message[EXCEED_MINT_LIMIT_ERROR] = "Sorry, your total mint count exceeds the mint limit, please check it and try again."
	message[ONLY_WHITELIST_ERROR] = "Only allow whitelist to mint"
	message[INVALID_TOKEN_ERROR] = "Token expired, please refresh and try again."
	message[DB_ERROR] = "system error"
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
