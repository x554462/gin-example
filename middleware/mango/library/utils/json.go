package utils

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/x554462/gin-example/middleware/mango/library/exception"
)

func JsonEncode(data interface{}) string {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	b, err := json.Marshal(data)
	exception.CheckError(err)
	return string(b)
}

func JsonEncodeByte(data interface{}) []byte {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	b, err := json.Marshal(data)
	exception.CheckError(err)
	return b
}

func JsonDecode(str string, v interface{}) {
	JsonDecodeWithByte([]byte(str), v)
}

func JsonDecodeWithByte(data []byte, v interface{}) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	exception.CheckError(json.Unmarshal(data, &v))
}
