package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

type DataConfig struct {
	Host string `json:"Host"`
	User string `json:"User"`
	Pass string `json:"Pass"`
}

var db *sql.DB

func Setup() {
	var err error

	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	configuration := DataConfig{}
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Fatal(err.Error())
	}

	dataSourceName := fmt.Sprintf("%s:%s@/%s", configuration.User, configuration.Pass, configuration.Host)
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal("Ошибка: соединение с бд не удалось.")
	}
}
