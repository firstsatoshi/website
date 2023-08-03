package xerr

var message map[uint32]string

func init() {
	message = make(map[uint32]string)
	message[OK] = "ok"
	message[SERVER_COMMON_ERROR] = "System error, please try again later"
	message[REUQEST_PARAM_ERROR] = "Invalid request parameter error, please check it"
	message[INVALID_EMAIL_ERROR] = "Invalid email"
	message[TOO_MANY_REQUEST_ERROR] = "Too many requests, please try it later"
	message[INVALID_BTCADDRESS_ERROR] = "Invalid Bitcoin Taproot(P2TR) inscription receive address, it must starts with 'bc1p'"
	message[FEERATE_TOO_SMALL_ERROR] = "Fee rate is not safe, too small or too large, please set a proper fee rate and try again"
	message[COUNT_EXCEED_PER_ORDER_LIMIT_ERROR] = "Sorry mint file count exceed limit count, please reduce file count and try again"
	message[EVENT_NOT_EXISTS_ERROR] = "event does not exists"
	message[AVAILABLE_COUNT_IS_NOT_ENOUGH] = "Sorry, current avialable count is not enough."
	message[EXCEED_MINT_LIMIT_ERROR] = "Sorry, your total mint count exceeds the mint limit, please check it and try again."
	message[ONLY_WHITELIST_ERROR] = "Only allow whitelist to mint"
	message[INVALID_TOKEN_ERROR] = "Token expired, please refresh and try again."
	message[GET_BRC20_INFO_ERROR] = "Get BRC20 info error."
	message[TOO_LARGE_DATA_ERROR] = "inscribe data total size exceed limit size. please reduce data size and try again."
	message[INVALID_NAME_TYPE_ERROR] = "invalid name type error."
	message[INVALID_NAME_ERROR] = "invalid name error."
	message[BITFISH_HAS_BEEN_BORN_ERROR] = "Sorry, this fish is duplicated with existed fishes, please change fish's elements and try again."
	message[DUPLICATED_BITFISH_IN_ONE_ORDER_ERROR] = "Sorry, there are duplicated fishes in this order, please your order and try again."

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
