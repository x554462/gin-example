package model

import "github.com/x554462/gin-example/middleware/mango/library/exception"

var AdminUserNotFoundError = exception.ModelNotFoundError

type AdminUser struct {
	Id       int    `db:"id"`
	Account  string `db:"account"`
	Passport string `db:"passport"`
	Name     string `db:"name"`
	RoleId   int    `db:"role_id"`
}

func (this *AdminUser) GetIndexValues() []interface{} {
	return []interface{}{this.Id}
}

func (this *AdminUser) InitModelInfo() (tableName string, indexFields []string, notFoundErr exception.ErrorWrap) {
	return "admin_user", []string{"id"}, exception.New("用户未找到", AdminUserNotFoundError)
}

func (this *AdminUser) VerifyPassport(passport string) bool {
	return this.Passport == passport
}
