package webutil

type Response struct {
	Respcd  string      `json:"respcd"`
	Respmsg string      `json:"respmsg"`
	Resperr string      `json:"resperr"`
	Data    interface{} `json:"data"`
}

// from qfcommon
const (
	OK string = "0000"

	ERR_DB      string = "2000"
	ERR_RPC     string = "2001"
	ERR_SESSION string = "2002"
	ERR_DATA    string = "2003"
	ERR_IO      string = "2004"

	ERR_LOGIN string = "2100"
	ERR_PARAM string = "2101"
	ERR_USER  string = "2102"
	ERR_ROLE  string = "2103"
	ERR_PWD   string = "2104"

	ERR_REQUEST string = "2200"
	ERR_IP      string = "2201"
	ERR_MAC     string = "2202"

	ERR_NODATA    string = "2300"
	ERR_DATAEXIST string = "2301"

	ERR_UNKNOW string = "2400"
)

var (
	ErrMsg map[string]string = map[string]string{
		OK:            "请求成功",
		ERR_DB:        "数据库错误",
		ERR_RPC:       "内部服务错误",
		ERR_SESSION:   "用户未登陆",
		ERR_DATA:      "数据错误",
		ERR_IO:        "输入输出错误",
		ERR_LOGIN:     "登陆错误",
		ERR_PARAM:     "参数错误",
		ERR_USER:      "用户错误",
		ERR_ROLE:      "角色错误",
		ERR_PWD:       "密码错误",
		ERR_REQUEST:   "非法请求",
		ERR_IP:        "IP受限",
		ERR_MAC:       "校验mac错误",
		ERR_NODATA:    "无数据",
		ERR_DATAEXIST: "数据已经存在",
		ERR_UNKNOW:    "未知错误",
	}
)

func Success(data interface{}, msg string) Response {
	return Response{
		Respcd:  OK,
		Respmsg: msg,
		Resperr: ErrMsg[OK],
		Data:    data,
	}
}

func Error(code string, data interface{}, msg string) Response {
	resperr, ok := ErrMsg[code]
	if !ok {
		resperr = "Error"
	}
	return Response{
		Respcd:  code,
		Respmsg: msg,
		Resperr: resperr,
		Data:    data,
	}
}
