package view

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Template struct {
	Context *gin.Context
	Name    string
	Obj     map[string]interface{}
}

func (t *Template) Assign(k string, i interface{}) Template {
	if t.Obj == nil {
		t.Obj = make(map[string]interface{})
	}
	t.Obj[k] = i
	return *t
}

//渲染方法
func (t *Template) Fetch(name string) {
	c := t.Context
	c.HTML(http.StatusOK, name, t.Obj)
}

func (t *Template) Error(error string) {
	c := t.Context
	t.Obj["error"] = error
	c.HTML(http.StatusOK, "error.html", t.Obj)
}
