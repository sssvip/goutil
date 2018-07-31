package webserver

import "github.com/gin-gonic/gin"

func NewDebugWebServer() *gin.Engine {
	return gin.Default()
}

func NewWebServer() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	return gin.Default()
}
