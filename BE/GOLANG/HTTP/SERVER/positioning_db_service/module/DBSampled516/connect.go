package DBSampled516

import (
	"database/sql"
	"example.com/db_service/pkg/config"
	"fmt"
	"gorm.io/gorm/logger"
	"log"

	"github.com/prestodb/presto-go-client/presto"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB(cfg *config.Config) error {
	var prestoURI = fmt.Sprintf("https://%s:%s@%s:%s",
		cfg.Presto.UserName+";cloud=aws&mode=adhoc",
		cfg.Presto.Password,
		cfg.Presto.Host,
		cfg.Presto.Port)
	var prestoConfig = &presto.Config{
		PrestoURI: prestoURI,
		Catalog:   cfg.Presto.CatLog,
		Schema:    cfg.Presto.Schema,
	}
	var dsn, _ = prestoConfig.FormatDSN()
	sqldb, err := sql.Open("presto", dsn)
	if err != nil {
		log.Printf("failed to connect to db by sql, error is:%v\n", err.Error())
		return err
	}

	sqldb.SetMaxOpenConns(100)
	sqldb.SetMaxIdleConns(20)
	db, err = gorm.Open(mysql.New(mysql.Config{Conn: sqldb, SkipInitializeWithVersion: true}), &gorm.Config{DisableAutomaticPing: true})

	db.Logger = logger.Default.LogMode(logger.Info)
	return nil
}

func GetDB() *gorm.DB {
	return db
}
