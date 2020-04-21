package model

import (
	"github.com/x554462/gin-example/middleware/mango/library/exception"
)

var AdminRoleNotFoundError = exception.ModelNotFoundError

type AdminRole struct {
	Id         int    `db:"id"`
	Name       string `db:"name"`
	Describe   string `db:"describe"`
	CreateTime int    `db:"create_time"`
}

func (this *AdminRole) GetIndexValues() []interface{} {
	return []interface{}{this.Id}
}

func (this *AdminRole) InitModelInfo() (tableName string, indexFields []string, notFoundErr exception.ErrorWrap) {
	return "admin_role", []string{"id"}, exception.New("用户角色未找到", AdminRoleNotFoundError)
}
