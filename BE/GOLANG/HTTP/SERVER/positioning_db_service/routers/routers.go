package routers

import (
	"example.com/db_service/service/booking"
	"example.com/db_service/service/driver"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// InitRouter initialize routing information
func InitRouters() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.Default())

	r.POST("/driver", driver.GetDriverGPS)
	r.POST("/booking", booking.GetBookingGPS)
	return r
}
