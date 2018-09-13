package webserver

import (
	"github.com/gin-gonic/gin"
)

func NewWebServerWithPProf() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	//ginpprof.Wrap(router) //"github.com/DeanThompson/ginpprof"
	return router
}
