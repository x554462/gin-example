# gin-example
* gin框架开发示例

# 介绍
* 封装底层Dao操作类，支持主从数据库操作
* 不包含第三方orm，使用[gendry](https://github.com/didi/gendry)辅助操作数据库
* 使用redis封装session管理
* error统一使用exception管理
* 支持swagger文档
* 使用[endless](https://github.com/fvbock/endless)优雅重启

# Dao抽象层

## 介绍
* dao抽象层为对常用的数据库操作的封装，避免在dao层书写重复的代码，仅通过简单的组合和两行代码就可以进行查询数据库的操作，如：
```go
	adminUserD := NewAdminAuthDao(dao.GetDaoSession())
	// 根据主键查询
	// 第一个参数为是否加for update锁，第二个参数为可变长参数，填入查询的主键值
	// 查询得到id为1用户
	user := adminUserD.Select(false, 1).(*model.AdminUser)
	fmt.Println(user)
```
## 具体用例
1、创建model文件，并实现GetIndexValues和InitModelInfo方法，使该model是ModelInterface的实现
```go
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
```

2、创建model的dao文件，并引用抽象Dao

```go
package dao

import (
	"github.com/x554462/gin-example/app/model"
	"github.com/x554462/gin-example/middleware/mango/dao"
)

type AdminUserDao struct {
	dao.Dao
}

func NewAdminUserDao(ds *dao.DaoSession) *AdminUserDao {
	return ds.GetDao(&model.AdminUser{}, &AdminUserDao{}).(*AdminUserDao)
}
```

3、调用查询方法示例
```go
	// 获得AdminUser的操作Dao
	adminUserD := NewAdminAuthDao(dao.GetDaoSession())
	// 根据主键查询
	// 第一个参数为是否加for update锁，第二个参数为可变长参数，填入查询的主键值
	// 查询得到id为1用户
	user := adminUserD.Select(false, 1).(*model.AdminUser)
	fmt.Println(user)

	// 根据where条件查询单条记录
	// 第一个参数为是否使用从库查询，第二个参数为查询条件
	// 去从库查询单个account为admin的用户
	user = adminUserD.SelectOne(true, map[string]interface{}{"account":"admin"}).(*model.AdminUser)
	fmt.Println(user)

	// 根据where条件查询多条记录
	// 第一个参数为是否使用从库查询，第二个参数为查询条件
	// 去从库查询单个name为abc123的用户
	users := adminUserD.SelectMulti(true, map[string]interface{}{"name":"abc123"})
	for _, v := range users {
		user = v.(*model.AdminUser)
		fmt.Println(user)
	}
```

# exception封装

## 介绍
* 在编程过程中，可以经常遇到一些可以直接报错给用户。标准的go编程中，推荐我们应在每个方法返回错误的消息，然后在外层根据错误提示用户。这类方法在调用层数不多的情况下是很好用的，但是在一些调用层数过深的情况下，每一级调用都返回错误处理，显然加大了我们编程的心智负担。所以exception提供了一种统一的处理方式，对于可以直接提示给用户的错误提示，通过直接panic来中断程序的执行，并提示对应的错误信息给到用户。

示例：
```go
	// 获得AdminUser的操作Dao
	adminUserD := NewAdminAuthDao(dao.GetDaoSession())
	
    // 封装来TryCatch方法
    // 第一个参数：执行的方法体，相当于try
    // 第二个参数：报错后执行的方法， 相当于catch
    // 第三个参数：是一个可变长的错误类型，用来判断try中报错的类型是否在预定义的类型里面
	exception.TryCatch(func() {
		// 根据主键查询
		// 第一个参数为是否加for update锁，第二个参数为可变长参数，填入查询的主键值
		// 查询得到id为1用户
		user := adminUserD.Select(false, 1).(*model.AdminUser)
        // 有可能会报记录未找到
		fmt.Println(user)
	}, func(err error) {
        // 捕获错误
		fmt.Println(err)
	}, model.AdminUserNotFoundError)
```