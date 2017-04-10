package main

import (
	"database/sql"
	"log"
	//"strconv"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

func connsql(dbhost, dbuser, dbpw, dbname string) {
	db, err = sql.Open("mysql", dbuser+":"+dbpw+"@"+dbhost+"/"+dbname)
	if err != nil {
		log.Fatalf("打开数据库失败: %s\n", err)
	}
}
