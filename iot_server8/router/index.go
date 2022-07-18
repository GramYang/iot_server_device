package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil) //禁用警告信息
	r.Use(cors.Default())
	loginRouter(r, "")
	return r
}
