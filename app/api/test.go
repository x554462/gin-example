package api

import (
	"github.com/gin-gonic/gin"
	"github.com/x554462/gin-example/app/dao"
	"github.com/x554462/gin-example/middleware/mango"
)
// @Summary 测试get接口
// @Produce  json
// @Param id path int true "id"
// @Param data query int true "data"
// @Success 200 {object} mango.Response
// @Router /api/{id} [get]
func TestGet(c *gin.Context) {

	ctrl := mango.Default(c)
	data := ctrl.GetQuery("data", true).Name("data").Min(2).Max(10).Int64()
	id := ctrl.GetPar("id").Name("id").Int()
	ctrl.JsonResponse(map[string]interface{}{"data":data, "sss": id})
}

// @Summary 测试post接口
// @Produce  json
// @Param id path int true "id"
// @Param data query int true "data"
// @Param data body string true "Data"
// @Success 200 {object} mango.Response
// @Router /api/{id} [post]
func TestPost(c *gin.Context) {

	ctrl := mango.Default(c)
	data := ctrl.GetQuery("data", true).Name("data").Min(2).Max(10).Int64()
	id := ctrl.GetPar("id").Name("id").Int()

	//logging.Error(exception.RootError)
	var t struct {
		Data string `json:"data" validate:"integer=as,1,10"`
	}
	ctrl.ParsePost(&t)
	testD := dao.NewTestDao(ctrl.GetDaoSession())
	//b := testD.Insert(map[string]interface{}{"name":t.Data})
	//testD.Update(b, map[string]interface{}{"name":"545676"})
	testD.GetDaoSession().BeginTransaction()
	b := testD.Select(true, t.Data)
	//testD.GetDaoSession().SubmitTransaction()
	b = testD.Select(true, t.Data)
	b = testD.Select(true, t.Data)
	b = testD.Select(true, t.Data)
	b = testD.Select(true, t.Data)
	testD.GetDaoSession().SubmitTransaction()


	ctrl.JsonResponse(map[string]interface{}{"123": data, "sss": id, "db":b})
}
