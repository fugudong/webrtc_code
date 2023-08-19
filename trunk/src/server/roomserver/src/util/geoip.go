// 获取对应IP的国家和城市
package util

import (
	"github.com/oschwald/geoip2-golang"
	"net"
)

// 保留 国家和省份
//
// 增加ice地址
// https://ip.cn/index.php 查询网址
type GeoIpMap struct {
	db *geoip2.Reader
}
var GeoIpTable *GeoIpMap

func  GeoIpMapNew(path string) *GeoIpMap  {
	var geoIpMap *GeoIpMap
	db, err := geoip2.Open(path)
	if err != nil {
		return  nil
	}
	geoIpMap = &GeoIpMap{db: db}
	return geoIpMap
}

// 返回国家
func (geoIpMap *GeoIpMap)GeoIpGetCountry(ip string) string {
	ip_ := net.ParseIP(ip)
	record, err := geoIpMap.db.City(ip_)
	if err != nil {
		return ""		// 返回空
	}
	return  record.Country.IsoCode
}

// 返回省 / 州
func (geoIpMap *GeoIpMap)GeoIpGetSubdivisions(ip string) string {
	ip_ := net.ParseIP(ip)
	record, err := geoIpMap.db.City(ip_)
	if err != nil {
		return ""		// 返回空
	}
	if len(record.Subdivisions) == 0 {

		return ""
	}
	return  record.Subdivisions[0].Names["en"]
}

// 返回城市
func (geoIpMap *GeoIpMap)GeoIpGetCity(ip string) string {
	ip_ := net.ParseIP(ip)
	record, err := geoIpMap.db.City(ip_)
	if err != nil {
		return ""		// 返回空
	}
	return  record.City.Names["pt-BR"]
}
