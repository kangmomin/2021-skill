package conn

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type DBConfig struct {
	Driver   string `json:"driver"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

var DB = connectDb()

func connectDb() (db *sql.DB) {
	config := getConfig()
	db, err := sql.Open(config.Driver, config.User+":"+config.Password+"@/"+config.Database)

	if err != nil {
		panic(err.Error())
	}

	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(2)

	return db
}

func getConfig() DBConfig {
	var config DBConfig
	data, err := os.Open("config/db.json")
	if err != nil {
		panic(err.Error())
	}
	defer data.Close()

	byteValue, err := ioutil.ReadAll(data)
	if err != nil {
		panic(err.Error())
	}

	err = json.Unmarshal(byteValue, &config)

	if err != nil {
		panic(err.Error())
	}

	return config
}
