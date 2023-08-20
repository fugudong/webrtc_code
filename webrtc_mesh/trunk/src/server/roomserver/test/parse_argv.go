// 测试输入参数 int 参数，字符串参数
package main

import(
	"flag"
	"fmt"
)
var n = flag.Int("n",1,"number of page")
var configPath = flag.String("config","config.ini","config file path")
var port = flag.Int("port", 8090, "lisent port")
var tls_crt = flag.String("tls_crt","1_easywebrtc.com_bundle.crt","tls crt path")
var tls_key = flag.String("tls_key","2_easywebrtc.com_bundle.key","tls key path")
// Go实战--golang中使用HTTPS以及TSL(.crt、.key、.pem区别以及crypto/tls包介绍)
// https://blog.csdn.net/wangshubo1989/article/details/77508738
func main() {
	flag.Parse()
	fmt.Println(*n)
	fmt.Println("config path",*configPath)
	fmt.Println("port", *port)
}
