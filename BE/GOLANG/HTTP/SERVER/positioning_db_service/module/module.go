package module

import (
	"example.com/db_service/module/DBSampled516"
	"example.com/db_service/pkg/config"
	"fmt"
)

func Init(cfg *config.Config) error {
	err := DBSampled516.InitDB(cfg)
	if err != nil {
		fmt.Printf("Failed to init presto, if you don't use it, just ignore\n")
	}
	return nil
}
