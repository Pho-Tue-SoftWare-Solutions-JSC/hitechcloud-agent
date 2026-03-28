package v2

import (
	"net/http"

	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/app/dto/request"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/app/service"
	"github.com/gin-gonic/gin"
)

var hostingAccountService = service.NewIHostingAccountService()

// CreateHostingAccount creates a new hosting account
func (b *BaseApi) CreateHostingAccount(c *gin.Context) {
	var req request.HostingAccountCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	if err := hostingAccountService.Create(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Account created successfully"})
}

// GetHostingAccount returns account info
func (b *BaseApi) GetHostingAccount(c *gin.Context) {
	username := c.Param("username")
	info, err := hostingAccountService.GetByUsername(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": info})
}

// ListHostingAccounts returns all hosting accounts
func (b *BaseApi) ListHostingAccounts(c *gin.Context) {
	accounts, err := hostingAccountService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": accounts})
}

// UpdateHostingAccount updates account package/limits
func (b *BaseApi) UpdateHostingAccount(c *gin.Context) {
	username := c.Param("username")
	var req request.HostingAccountUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	if err := hostingAccountService.Update(username, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Account updated successfully"})
}

// SuspendHostingAccount suspends an account
func (b *BaseApi) SuspendHostingAccount(c *gin.Context) {
	username := c.Param("username")
	if err := hostingAccountService.Suspend(username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Account suspended"})
}

// UnsuspendHostingAccount unsuspends an account
func (b *BaseApi) UnsuspendHostingAccount(c *gin.Context) {
	username := c.Param("username")
	if err := hostingAccountService.Unsuspend(username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Account unsuspended"})
}

// TerminateHostingAccount deletes an account and all resources
func (b *BaseApi) TerminateHostingAccount(c *gin.Context) {
	username := c.Param("username")
	if err := hostingAccountService.Terminate(username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Account terminated"})
}

// ChangeHostingAccountPassword changes account password
func (b *BaseApi) ChangeHostingAccountPassword(c *gin.Context) {
	username := c.Param("username")
	var req request.HostingAccountPassword
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	if err := hostingAccountService.ChangePassword(username, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Password changed"})
}

// GetHostingAccountStats returns resource usage stats
func (b *BaseApi) GetHostingAccountStats(c *gin.Context) {
	username := c.Param("username")
	stats, err := hostingAccountService.GetStats(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": stats})
}
