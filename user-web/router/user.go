package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"shop-api/user-web/api"
	"shop-api/user-web/middlewares"
)

func InitUser(r *gin.RouterGroup) {

	zap.S().Infof("初始化路由:%s\n", "user")
	ur := r.Group("user")
	{
		ur.POST("/pwd_login", api.PasswordLogin)
		ur.POST("/create", api.CreateUser)
		ur.GET("/list", api.GetUserList).Use(middlewares.JWTAuth())
	}

}
