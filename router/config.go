package router

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
		JwtHandler  
}

type (
	JwtHandler func(relativePath string, handlers []gin.HandlerFunc) (gin.IRoutes)
)




