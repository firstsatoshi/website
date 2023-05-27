package xerr

// 成功返回
const OK uint32 = 0

/**(前3位代表业务,后三位代表具体功能)**/

// 全局错误码
const SERVER_COMMON_ERROR uint32 = 100001
const REUQEST_PARAM_ERROR uint32 = 100002
const INVALID_EMAIL_ERROR uint32 = 100003
const TOO_MANY_REQUEST_ERROR uint32 = 100004
const INVALID_BTCP2TRADDRESS_ERROR uint32 = 100005
const FEERATE_TOO_SMALL_ERROR uint32 = 100006
const DB_ERROR uint32 = 111111


//用户模块
