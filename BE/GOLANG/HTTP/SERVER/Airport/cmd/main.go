package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"xiaolong.ji.com/airport/airport-service/routers"
	_ "xiaolong.ji.com/airport/pkg/config"
)

func main() {
	gin.SetMode(gin.DebugMode)
	route := routers.InitRouters()
	endPoint := fmt.Sprintf(":%d", 8006)
	server := &http.Server{
		Addr:    endPoint,
		Handler: route,
	}
	log.Printf("start http server, listening %s", endPoint)
	_ = server.ListenAndServe()
}
