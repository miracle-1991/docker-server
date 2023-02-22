package routers

import "github.com/gin-gonic/gin"

// InitRouter initialize routing information
func InitRouters() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.POST("/driver")
}
