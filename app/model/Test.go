package model

import (
	"github.com/x554462/gin-example/middleware/mango/library/exception"
)

var NotFoundError = exception.New("Test not found", exception.ModelNotFoundError)

type Test struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
}

func (t *Test) GetIndexValues() []interface{} {
	return []interface{}{t.Id}
}

func (t *Test) InitModelInfo() (tableName string, indexFields []string, notFoundErr exception.ErrorWrap) {
	return "test", []string{"id"}, NotFoundError
}
