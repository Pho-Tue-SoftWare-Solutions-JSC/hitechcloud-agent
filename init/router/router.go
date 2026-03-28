package router

import (
	v2 "github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/app/api/v2"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/global"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/i18n"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/middleware"
	rou "github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/router"
	"github.com/gin-gonic/gin"
)

var (
	Router *gin.Engine
)

func Routers() *gin.Engine {
	Router = gin.Default()
	Router.Use(i18n.UseI18n())

	// Public health check (no auth required)
	Router.GET("/api/v1/health", v2.ApiGroupApp.BaseApi.CheckHealth)

	// Private API group with API Key auth
	PrivateGroup := Router.Group("/api/v1")
	if global.IsMaster {
		// Master mode: no auth (unix socket)
	} else {
		PrivateGroup.Use(middleware.ApiKeyAuth())
	}
	PrivateGroup.Use(middleware.OperationResolveMeta())
	for _, router := range rou.RouterGroupApp {
		router.InitRouter(PrivateGroup)
	}

	// Keep v2 routes for backward compatibility during migration
	V2Group := Router.Group("/api/v2")
	if global.IsMaster {
	} else {
		V2Group.Use(middleware.ApiKeyAuth())
	}
	V2Group.Use(middleware.OperationResolveMeta())
	for _, router := range rou.RouterGroupApp {
		router.InitRouter(V2Group)
	}

	return Router
}
