package db

import (
	"github.com/gocraft/dbr"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var Conn *dbr.Connection

func init() {
	var err error
	dsn := "root@tcp/roadster?parseTime=true"
	Conn, err = dbr.Open("mysql", dsn, &dbr.NullEventReceiver{})
	if err != nil {
		log.Fatal(err)
	}
}
