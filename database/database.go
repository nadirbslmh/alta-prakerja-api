package database

import (
	"database/sql"
	"fmt"
	"gugcp/utils"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		utils.GetConfig("DB_USERNAME"),
		utils.GetConfig("DB_PASSWORD"),
		utils.GetConfig("DB_HOST"),
		utils.GetConfig("DB_PORT"),
		utils.GetConfig("DB_NAME"),
	)

	var err error

	DB, err = sql.Open("mysql", dsn)

	if err != nil {
		panic(err.Error())
	}

	log.Println("connected to the database")
}
