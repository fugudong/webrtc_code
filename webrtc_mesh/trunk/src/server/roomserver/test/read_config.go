package main

import (
	"github.com/Unknwon/goconfig"
	"log"
)

//1. 读取配置文件
func readConfiguration()  {
	cfg, err := goconfig.LoadConfigFile("config.ini")
	if err != nil {
		log.Println("读取配置文件失败[config.ini]")
		return
	}
	username, err := cfg.GetValue("mysql", "username")
	log.Println("username = ", username)
	password, err := cfg.GetValue("mysql", "password")
	log.Printf("password = %s", password)
	mysql, err := cfg.GetSection("mysql")

	log.Println("mysql = ", mysql)
}

func main() {
	log.Println("开始干活了.....")
	readConfiguration()
}