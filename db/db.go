package db

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

var DB *sql.DB

func Init(database string, maxidle, maxopen int) {
	var err error
	DB, err = sql.Open("mysql", database)
	if err != nil {
		log.Fatalln("open db fail:", err)
	}
	DB.SetMaxIdleConns(maxidle)
	DB.SetMaxOpenConns(maxopen)
	err = DB.Ping()
	if err != nil {
		log.Fatalln("ping db fail:", err)
	}
}
