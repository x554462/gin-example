package utils

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/x554462/gin-example/middleware/mango/library/exception"
)

type jsonOb struct {
	data []byte
}

func JsonParse(data []byte) *jsonOb {
	return &jsonOb{data: data}
}

func (jsonOb *jsonOb) Get(path ...interface{}) jsoniter.Any {
	return jsoniter.Get(jsonOb.data, path...)
}

func (jsonOb *jsonOb) MustGet(path ...interface{}) jsoniter.Any {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); !ok {
				exception.CheckError(err)
			}
		}
	}()
	return jsoniter.Get(jsonOb.data, path...).MustBeValid()
}
