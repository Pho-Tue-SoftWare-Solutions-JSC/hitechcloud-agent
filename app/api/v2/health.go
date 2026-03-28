package v2

import (
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/app/api/v2/helper"
	"github.com/gin-gonic/gin"
)

func (b *BaseApi) CheckHealth(c *gin.Context) {
	helper.Success(c)
}
