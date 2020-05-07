package controller

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type RestControllerInterface interface {
	Get(c *gin.Context)
	Show(c *gin.Context)
	Edit(c *gin.Context)
	Store(c *gin.Context)
	Delete(c *gin.Context)
}

type RestControllerStruct struct {
}

type returnData struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func Success(c *gin.Context, msg string, data interface{}) {
	var success returnData
	success = returnData{1, msg}
	c.JSON(http.StatusOK, success)
}
