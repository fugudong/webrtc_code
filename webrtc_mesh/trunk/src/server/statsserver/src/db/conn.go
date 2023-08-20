package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/cihub/seelog"
)

var (
	dbConn *sql.DB
	err    error
)

type MysqlConfig struct {
	UserName string			// 用户名
	Password string		// 密码
	DatabaseName string	// 库名
	Url string 			// 数据库地址
}

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

func ConnectDB(config *MysqlConfig) error{
	//	"root:123456@tcp(localhost:3306)/rtc_room_server?charset=utf8&parseTime=true&loc=Local")
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&loc=Local",
		config.UserName,
		config.Password,
		config.Url,
		config.DatabaseName)
	log.Infof("dataSourceName:%s", dataSourceName)
	dbConn, err = sql.Open("mysql", dataSourceName)
	return err
}

func DisconnectDB()  error {
	if dbConn != nil {
		return dbConn.Close()
	}
	return nil
}

