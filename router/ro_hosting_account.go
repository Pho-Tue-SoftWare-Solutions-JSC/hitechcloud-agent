package router

import (
	v2 "github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/app/api/v2"
	"github.com/gin-gonic/gin"
)

type HostingAccountRouter struct{}

func (s *HostingAccountRouter) InitRouter(Router *gin.RouterGroup) {
	accountRouter := Router.Group("accounts")
	baseApi := v2.ApiGroupApp.BaseApi

	{
		accountRouter.POST("", baseApi.CreateHostingAccount)
		accountRouter.GET("", baseApi.ListHostingAccounts)
		accountRouter.GET("/:username", baseApi.GetHostingAccount)
		accountRouter.PUT("/:username", baseApi.UpdateHostingAccount)
		accountRouter.DELETE("/:username", baseApi.TerminateHostingAccount)
		accountRouter.POST("/:username/suspend", baseApi.SuspendHostingAccount)
		accountRouter.POST("/:username/unsuspend", baseApi.UnsuspendHostingAccount)
		accountRouter.PUT("/:username/password", baseApi.ChangeHostingAccountPassword)
		accountRouter.GET("/:username/stats", baseApi.GetHostingAccountStats)
	}
}
