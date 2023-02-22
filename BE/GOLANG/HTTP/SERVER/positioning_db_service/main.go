package main

import (
	"example.com/db_service/module/DBSampled516"
	"example.com/db_service/pkg/config"
	"example.com/db_service/routers"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	err = config.Init(pwd + "/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	err = DBSampled516.InitDB(config.GetConfig())
	if err != nil {
		log.Fatal(err)
	}

	gin.SetMode(gin.DebugMode)
	routersInit := routers.InitRouters()
	endPoint := fmt.Sprintf(":%d", 8003)
	server := &http.Server{
		Addr:    endPoint,
		Handler: routersInit,
	}
	log.Printf("[info] start http server listening %s", endPoint)
	server.ListenAndServe()
}
