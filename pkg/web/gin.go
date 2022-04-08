package web

import (
	"github.com/gin-gonic/gin"
	"github.com/zldongly/email_server/pkg/errors"
	"net/http"
)

type response struct {
	Data    interface{} `json:"data"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
}

func ResponseHttp(c *gin.Context, data interface{}, err error) {
	var (
		code = http.StatusOK
		msg  = ""
	)
	if err != nil {
		ec := errors.AsCode(err)
		code = ec.Code()
		msg = ec.Message()
	}

	c.JSON(code, &response{
		Data:    data,
		Code:    code,
		Message: msg,
	})
}
