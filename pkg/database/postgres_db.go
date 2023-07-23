package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"strconv"
)

type DBConfig struct {
	Host       string
	Port       string
	User       string
	Password   string
	DBName     string
	DriverName string
}

func NewDB(conf DBConfig) *sql.DB {
	portInt, err := strconv.Atoi(conf.Port)
	if err != nil {
		log.Fatalf("error occurs while conver string to int: %v", err)
	}

	dataSourceName := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
		conf.User, conf.Password, conf.DBName, conf.Host, portInt)

	db, err := sql.Open(conf.DriverName, dataSourceName)
	if err != nil {
		log.Fatalf("error occures while connecting to database: %v", err)
	}
	return db
}
