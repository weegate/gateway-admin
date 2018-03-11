//@author wuyong
//@date   2018/1/9
//@desc

package lib

type Error struct {
	Status  int         `json:"status"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

var errors = map[string]Error{
	"REQUEST_PARAMS_VALIDATE_ERROR": {50000, "request params validate is error! please check request params", nil},

	"REQUEST_PARAMS_ID_ERROR": {40000, "request params id is error!", nil},
	"REQUEST_SERVER_ERROR":    {41001, "ops,request tripartite server error!", nil},

	"UNAUTHORIZED":         {30000, "this user unauthorized,please check it!", nil},
	"UNKNOWN_METHOD":       {30001, "unknown method!", nil},
	"UNKNOWN_TYPE":         {30002, "unknown type!", nil},
	"UNKNOWN_JSON_FORMAT":  {30003, "unknown json format!", nil},
	"JSON_UNMARSHAL_ERROR": {30003, "json unmarshal error!", nil},

	"DB_OP_ERROR":      {20000, "ops, DB error!", nil},
	"ZK_OP_ERROR":      {21000, "ops, ZK error!", nil},
	"ZK_NO_PATH_ERROR": {21001, "ops, ZK no node path !", nil},

	"DEFAULT": {10000, "unknown error!", nil},
}

func GetErrorByName(name string) Error {
	if error, ok := errors[name]; ok {
		return error
	} else {
		return errors["DEFAULT"]
	}
}

func GetSuccess(data interface{}) Error {
	return Error{0, "", data}
}
