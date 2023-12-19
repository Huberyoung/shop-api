package initialize

import (
	"github.com/gin-gonic/gin"
	"shop-api/user-web/global"

	"shop-api/user-web/router"
)

func InitRouter() *gin.Engine {
	gin.SetMode(global.ServerSetting.RunMode)

	e := gin.Default()
	v1Group := e.Group("/u/v1")
	router.InitUser(v1Group)
	return e
}
