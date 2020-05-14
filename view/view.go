package view

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type TemplateStruct struct {
	Context *gin.Context
	Name string
	Obj map[string]interface{}
}

var Template TemplateStruct

func Assign(k string,i interface{})  {
	if Template.Obj == nil {
		Template.Obj = make(map[string]interface{})
	}
	Template.Obj[k] = i
}

//渲染方法
func Fetch(name string){
	c := Template.Context
	c.HTML(http.StatusOK, name, Template.Obj)
}
