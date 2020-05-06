package mango

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/x554462/gin-example/middleware/mango/dao"
	"github.com/x554462/gin-example/middleware/mango/library/exception"
	"github.com/x554462/gin-example/middleware/mango/library/logging"
	"github.com/x554462/gin-example/middleware/mango/library/utils"
	"github.com/x554462/gin-example/middleware/mango/validator"
	"net/http"
)

const DefaultKey = "middleware/mango"

type Controller struct {
	ginCtx           *gin.Context
	session          *session
	daoSession       *dao.DaoSession
	responseFinished bool
	firstPanicOffset int
}

type Response struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
}

func New() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := newSession(c.Request, c.Writer)
		defer session.expiry()
		daoSession := dao.GetDaoSession(c.Request.Context())
		defer daoSession.Close()
		ctrl := &Controller{ginCtx: c, daoSession: daoSession, session: session}
		defer func() {
			if v := recover(); v != nil {
				if err, ok := v.(error); ok {
					for _, re := range []error{
						exception.ValidateError,
						exception.ResponseError,
						exception.UnauthorizedError,
						exception.ModelNotFoundError,
						exception.ModelRuntimeError,
						exception.JsonRuntimeError,
						exception.RuntimeError,
						exception.LoadConfigError,
						exception.RootError,
					} {
						if errors.Is(err, re) {
							if werr, ok := err.(exception.ErrorWrapInterface); ok {
								if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut {
									ctrl.JsonResponseWithMsg(nil, werr.GetMsg(), werr.GetCode())
								} else {
									ctrl.Echo(fmt.Sprintf("%d:%s", werr.GetCode(), werr.GetMsg()))
								}
								logging.ErrorWithPrefix(werr.GetFileLine(), werr)
							}
							return
						}
					}
				}
				ctrl.Echo(fmt.Sprintf("%v", v))
			}
		}()
		c.Set(DefaultKey, ctrl)
		c.Next()
		ctrl.JsonResponse(nil)
	}
}

func Default(c *gin.Context) *Controller {
	return c.MustGet(DefaultKey).(*Controller)
}

func (ctrl *Controller) GetDaoSession() *dao.DaoSession {
	return ctrl.daoSession
}

func (ctrl *Controller) GetPar(key string) validator.ValueInterface {
	if v, ok := ctrl.ginCtx.Params.Get(key); ok {
		return validator.NewValue(v)
	}
	return validator.NewNil(true)
}

func (ctrl *Controller) GetQuery(key string, must bool) validator.ValueInterface {
	if v, ok := ctrl.ginCtx.GetQuery(key); ok {
		return validator.NewValue(v)
	}
	return validator.NewNil(must)
}

func (ctrl *Controller) DefaultQuery(key string, defaultValue string) validator.ValueInterface {
	if v, ok := ctrl.ginCtx.GetQuery(key); ok {
		return validator.NewValue(v)
	}
	return validator.NewValue(defaultValue).NoValidate()
}

func (ctrl *Controller) GetForm(key string, must bool) validator.ValueInterface {
	if v, ok := ctrl.ginCtx.GetPostForm(key); ok {
		return validator.NewValue(v)
	}
	return validator.NewNil(must)
}

func (ctrl *Controller) DefaultForm(key string, defaultValue string) validator.ValueInterface {
	if v, ok := ctrl.ginCtx.GetPostForm(key); ok {
		return validator.NewValue(v)
	}
	return validator.NewValue(defaultValue).NoValidate()
}

func (ctrl *Controller) ParsePost(v interface{}) {
	binding.Validator = validator.NewValidator()
	err := ctrl.ginCtx.Bind(v)
	if err != nil {
		exception.ThrowMsgWithCallerDepth(err.Error(), exception.ValidateError, 3)
	}
}

func (ctrl *Controller) JsonResponse(data interface{}) {
	if ctrl.responseFinished {
		return
	}
	ctrl.responseFinished = true
	response := utils.JsonEncode(Response{
		Data:    data,
		Message: "ok",
		Code:    200,
	})
	EchoResponse(ctrl.ginCtx.Writer, response)
}

func (ctrl *Controller) JsonResponseWithMsg(data interface{}, message string, code int) {
	if ctrl.responseFinished {
		return
	}
	ctrl.EndRequest()
	response := utils.JsonEncode(Response{
		Data:    data,
		Message: message,
		Code:    code,
	})
	EchoResponse(ctrl.ginCtx.Writer, response)
}

func (ctrl *Controller) EndRequest() {
	ctrl.responseFinished = true
}

func (ctrl *Controller) Echo(response string) {
	ctrl.EndRequest()
	EchoResponse(ctrl.ginCtx.Writer, response)
}

func (ctrl *Controller) GetSession() *session {
	return ctrl.session
}

func EchoResponse(writer http.ResponseWriter, response string) {
	_, _ = writer.Write([]byte(response))
}
