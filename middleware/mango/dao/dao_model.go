package dao

import (
	"github.com/x554462/gin-example/middleware/mango/library/exception"
)

type ModelInterface interface {
	InitModelInfo() (tableName string, indexFields []string, notFoundErr exception.ErrorWrap)
	GetIndexValues() []interface{}
}
