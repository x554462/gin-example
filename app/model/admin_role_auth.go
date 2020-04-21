package model

import "github.com/x554462/gin-example/middleware/mango/library/exception"

var AdminRoleAuthNotFoundError = exception.ModelNotFoundError

type AdminRoleAuth struct {
	Id     int `db:"db"`
	RoleId int `db:"role_id"`
	AuthId int `db:"auth_id"`
}

func (this *AdminRoleAuth) GetIndexValues() []interface{} {
	return []interface{}{this.Id}
}

func (this *AdminRoleAuth) InitModelInfo() (tableName string, indexFields []string, notFoundErr exception.ErrorWrap) {
	return "admin_role_auth", []string{"id"}, exception.New("角色权限未找到", AdminRoleAuthNotFoundError)
}
