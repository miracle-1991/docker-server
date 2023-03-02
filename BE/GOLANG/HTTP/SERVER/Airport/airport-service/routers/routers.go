package routers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"xiaolong.ji.com/airport/airport-service/service/today"
)

func InitRouters() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.Default())

	r.POST("/today", today.GetSchedulesForToday)
	return r
}
