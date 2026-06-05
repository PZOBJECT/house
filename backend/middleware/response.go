package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "success",
		Data: data,
	})
}

func Error(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, Response{
		Code: 1,
		Msg:  msg,
		Data: nil,
	})
}

func ErrorWithStatus(c *gin.Context, httpStatus int, msg string) {
	c.JSON(httpStatus, Response{
		Code: 1,
		Msg:  msg,
		Data: nil,
	})
}
