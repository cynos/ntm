package tools

import (
	"database/sql"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	Host            string
	Port            string
	Name            string
	User            string
	Password        string
	ApplicationName string
	ConnectTimeout  int
	MaxOpenConn     int
	MaxIdleConn     int
}

func (c *DBConfig) dbinfo() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=disable password=%s connect_timeout=%d application_name=%s",
		c.Host, c.Port, c.User, c.Name, c.Password, c.ConnectTimeout, c.ApplicationName,
	)
}

func DBClient(config DBConfig) (*gorm.DB, *sql.DB) {
	db, err := gorm.Open(postgres.Open(config.dbinfo()), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	genericdb, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	genericdb.SetMaxIdleConns(config.MaxIdleConn)
	genericdb.SetMaxOpenConns(config.MaxOpenConn)

	return db, genericdb
}
