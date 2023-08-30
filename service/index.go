package service

import "github.com/gin-gonic/gin"

// GetIndex godoc
// @Tags 首页
// @Accept json
// @Produce json
// @Success 200 {string} pong
// @Router /index [get]
func GetIndex(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
