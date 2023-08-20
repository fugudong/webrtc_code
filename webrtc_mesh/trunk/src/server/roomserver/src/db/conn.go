package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var (
	dbConn *sql.DB
	err    error
)
//go语言中init函数用于包(package)的初始化，该函数是go语言的一个重要特性，
//func init() {
//	fmt.Println("Entering conn.go init function...")
//	//db, err := sql.Open("mysql", "用户名:密码@tcp(IP:端口)/数据库?charset=utf8")
//	dbConn, err = sql.Open("mysql",
//		"root:123456@tcp(localhost:3306)/rtc_room_server?charset=utf8&parseTime=true&loc=Local")
//	if err != nil {
//		panic(err.Error())
//	}
//	fmt.Printf("dbConn +%v\n", dbConn)
//}



func ConnectDB(userName string, password string, dbName string, dbUrl string, ) error{
	//	"root:123456@tcp(localhost:3306)/rtc_room_server?charset=utf8&parseTime=true&loc=Local")
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&loc=Local",
		userName, password, dbUrl, dbName)
	fmt.Println("dataSourceName:", dataSourceName)
	dbConn, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		return err
	}
	dbConn.SetMaxOpenConns(2000)
	dbConn.SetMaxIdleConns(1000)
	err := dbConn.Ping()
	return err
}

func DisconnectDB()  error {
	if dbConn != nil {
		return dbConn.Close()
	}
	return nil
}

