package main

import (
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"log"
	"net"
)

func main() {
	db, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// If you are using strings that may be invalid, check that ip is not nil
	//ip := net.ParseIP("175.10.52.99")	// 中国 湖南长沙，解析正确
	ip := net.ParseIP("124.156.13.115")	// 印度 腾讯, 解析为IN
	//ip := net.ParseIP("35.173.220.203")	// 美国
	//ip := net.ParseIP("82.209.192.0")
	//ip := net.ParseIP("129.204.197.215")	// 广州，被解析为北京

	record, err := db.City(ip)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Portuguese (BR) city name: %v\n", record.City.Names["pt-BR"])	// changsha
	//fmt.Printf("English subdivision name: %v\n", record.Subdivisions[0].Names["en"]) // hunan Virginia
	fmt.Printf("Russian country name: %v\n", record.Country.Names["ru"])
	// 中国CN， 美国US，印度IN, 俄罗斯 BY
	fmt.Printf("ISO country code: %v\n", record.Country.IsoCode)	// CN
	fmt.Printf("Time zone: %v\n", record.Location.TimeZone)
	fmt.Printf("Coordinates: %v, %v\n", record.Location.Latitude, record.Location.Longitude)
	// Output:
	// Portuguese (BR) city name: Londres
	// English subdivision name: England
	// Russian country name: Великобритания
	// ISO country code: GB
	// Time zone: Europe/London
	// Coordinates: 51.5142, -0.0931
}