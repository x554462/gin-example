package model

import "github.com/x554462/gin-example/middleware/mango/library/exception"

var AdminAuthNotFoundError = exception.ModelNotFoundError

const (
	AdminAuthTypeMenu int8 = 0
	AdminAuthTypeApi  int8 = 1
)

type AdminAuth struct {
	Id      int    `db:"id"`
	Name    string `db:"name"`
	Type    int8   `db:"type"`
	Path    string `db:"path"`
	Pid     int    `db:"pid"`
	SortNum int    `db:"sort_num"`
	Depth   int    `db:"depth"`
}

func (this *AdminAuth) GetIndexValues() []interface{} {
	return []interface{}{this.Id}
}

func (this *AdminAuth) InitModelInfo() (tableName string, indexFields []string, notFoundErr exception.ErrorWrap) {
	return "admin_auth", []string{"id"}, exception.New("路径未找到", AdminAuthNotFoundError)
}
