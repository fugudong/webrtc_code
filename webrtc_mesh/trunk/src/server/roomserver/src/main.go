package main

import (
	"roomserver/src/connect"
	"roomserver/src/db"
	"roomserver/src/logic"
	"roomserver/src/model"
	"roomserver/src/model/iceBalance"
	"roomserver/src/util"
	"errors"
	"time"
	"flag"
	"github.com/Unknwon/goconfig"
	Xlog "github.com/cihub/seelog"
	"log"
	"strconv"
)

func SetupLogger(logConfig string)  int{
	logger, err := Xlog.LoggerFromConfigAsFile(logConfig)
	if err != nil {
		log.Fatal("LoggerFromConfigAsFile failed:%s\n", err)
		return -1
	}

	Xlog.ReplaceLogger(logger)

	log.Printf("SetupLogger ok\n",)
	return 0
}

/**
配置文件参数
 */
type RoomServerConfig struct {
	// mysql数据库
	mysqlUser string			// 用户名
	mysqlPassword string		// 密码
	mysqlDatabaseName string	// 库名
	mysqlUrl string 			// 数据库地址

	// HTTP配置
	https int		// 是否为https, 0：http， 1：https
	httpListenIp string//https或https的监听IP
	httpListenPort	int	// http或https的监听端口
	httpsKeyPath string 	// server.crt", "server.key
	httpsCrtPath string		//

	// websocket配置
	wss int		// 是否为wss, 0：ws， 1：wss
	wsListenIp string	//ws或wss的监听IP
	wsListenPort	int			// ws或wss的监听端口
	wssKeyPath string 	// server.crt", "server.key
	wssCrtPath string		//

	iceConfs [] iceBalance.IceConf

	// client心跳包超时
	clientKeepLiveTimeout int
	// ice server心跳包超时
	iceKeepLiveTimeout int
	// 客户端报告统计结果的间隔
	clientReportStatsInterval int
	// 房间最大容纳人数
	maxRoomCapacity int
}

func readIceConfiguration(config *RoomServerConfig, cfg *goconfig.ConfigFile) error{
	iceType := cfg.MustValueArray("ice", "iceType", ",")
	bundlePolicy := cfg.MustValueArray("ice", "bundlePolicy", ",")
	rtcpMuxPolicy := cfg.MustValueArray("ice", "rtcpMuxPolicy",",")
	iceTransportPolicy := cfg.MustValueArray("ice", "iceTransportPolicy", ",")
	ip := cfg.MustValueArray("ice", "ip", ",")
	turnPort := cfg.MustValueArray("ice", "turnPort", ",")
	stunPort := cfg.MustValueArray("ice", "stunPort", ",")
	maxBandwidth := cfg.MustValueArray("ice", "maxBandwidth", ",")
	username := cfg.MustValueArray("ice", "username", ",")
	credential := cfg.MustValueArray("ice", "credential", ",")
	length := len(iceType)
	if length == 0 {
		return errors.New("length is 0")
	}
	Xlog.Infof("ICE have %d pair", length)
	if length != len(bundlePolicy) {
		return errors.New("bundlePolicy may be falied")
	}
	if length != len(rtcpMuxPolicy) {
		return errors.New("rtcpMuxPolicy may be falied")
	}
	if length != len(iceTransportPolicy) {
		return errors.New("iceTransportPolicy may be falied")
	}
	if length != len(ip) {
		return errors.New("ip may be falied")
	}
	if length != len(turnPort) {
		return errors.New("turnPort may be falied")
	}
	if length != len(stunPort) {
		return errors.New("stunPort may be falied")
	}
	if length != len(maxBandwidth) {
		return errors.New("maxBandwidth may be falied")
	}
	if length != len(username) {
		return errors.New("username may be falied")
	}
	if length != len(credential) {
		return errors.New("credential may be falied")
	}
	for i:=0; i< len(bundlePolicy); i++ {
		var iceConf iceBalance.IceConf
		iceConf.IceType = iceType[i]
		iceConf.BundlePolicy = bundlePolicy[i]
		iceConf.RtcpMuxPolicy = rtcpMuxPolicy[i]
		iceConf.IceTransportPolicy = iceTransportPolicy[i]
		iceConf.Ip = ip[i]
		valueInt, err := strconv.Atoi(turnPort[i])
		if err != nil {
			Xlog.Criticalf("无法转换%s成int类型：%s", turnPort[i], err)
			return  errors.New("Atoi turnPort[i] to int  failed")
		}
		iceConf.TurnPort = valueInt
		valueInt, err = strconv.Atoi(stunPort[i])
		if err != nil {
			Xlog.Criticalf("无法转换%s成int类型：%s", stunPort[i], err)
			return  errors.New("Atoi stunPort[i] to int  failed")
		}
		iceConf.StunPort = valueInt
		valueInt, err = strconv.Atoi(maxBandwidth[i])
		if err != nil {
			Xlog.Criticalf("无法转换%s成int类型：%s", maxBandwidth[i], err)
			return  errors.New("Atoi maxBandwidth[i] to int  failed")
		}
		iceConf.MaxBandwidth = valueInt
		iceConf.UserName = username[i]
		iceConf.Credential = credential[i]
		config.iceConfs = append(config.iceConfs, iceConf)
	}

	return nil
}
//1. 读取配置文件
func readConfiguration(filePath string)  (RoomServerConfig, error){
	var config RoomServerConfig

	Xlog.Info("LoadConfigFile")
	cfg, err := goconfig.LoadConfigFile(filePath)
	if err != nil {
		Xlog.Critical("Read config %s failed: %s", filePath, err.Error())
		return  config, errors.New("LoadConfigFile failed")
	}

	var valueString string
	var valueInt int

	Xlog.Info("default")
	// 读取 clientKeepLiveTimeout
	valueString, err = cfg.GetValue("default", "clientKeepLiveTimeout")
	if err != nil {
		Xlog.Criticalf("Can't get the key value(%s)：%s", "clientKeepLiveTimeout", err)
		return config, errors.New("GetValue clientKeepLiveTimeout failed")
	}
	valueInt, err = strconv.Atoi(valueString)
	if err != nil {
		Xlog.Criticalf("无法转换%s成int类型：%s", valueString, err)
		return config, errors.New("Atoi clientKeepLiveTimeout to int  failed")
	}
	config.clientKeepLiveTimeout = valueInt

	// 读取 iceKeepLiveTimeout
	valueString, err = cfg.GetValue("default", "iceKeepLiveTimeout")
	if err != nil {
		Xlog.Criticalf("Can't get the key value(%s)：%s", "iceKeepLiveTimeout", err)
		return config, errors.New("GetValue iceKeepLiveTimeout failed")
	}
	valueInt, err = strconv.Atoi(valueString)
	if err != nil {
		Xlog.Criticalf("无法转换%s成int类型：%s", valueString, err)
		return config, errors.New("Atoi iceKeepLiveTimeout to int  failed")
	}
	config.iceKeepLiveTimeout = valueInt

	// 读取clientReportStatsInterval
	valueString, err = cfg.GetValue("default", "clientReportStatsInterval")
	if err != nil {
		Xlog.Criticalf("Can't get the key value(%s)：%s", "clientReportStatsInterval", err)
		return config, errors.New("GetValue clientReportStatsInterval failed")
	}
	valueInt, err = strconv.Atoi(valueString)
	if err != nil {
		Xlog.Criticalf("无法转换%s成int类型：%s", valueString, err)
		return config, errors.New("Atoi clientReportStatsInterval to int  failed")
	}
	config.clientReportStatsInterval = valueInt

	// 读取maxRoomCapacity
	valueString, err = cfg.GetValue("default", "maxRoomCapacity")
	if err != nil {
		Xlog.Criticalf("Can't get the key value(%s)：%s", "maxRoomCapacity", err)
		return config, errors.New("GetValue maxRoomCapacity failed")
	}
	valueInt, err = strconv.Atoi(valueString)
	if err != nil {
		Xlog.Criticalf("无法转换%s成int类型：%s", valueString, err)
		return config, errors.New("Atoi maxRoomCapacity to int  failed")
	}
	config.maxRoomCapacity = valueInt

	Xlog.Info("mysql")
	// 读取mysql配置
	// 读取用户名
	valueString, err = cfg.GetValue("mysql", "mysqlUser")
	if err != nil {
		Xlog.Criticalf("Can't get the key value(%s), failed:%s", "mysqlUser", err)
		return  config, errors.New("GetValue mysqlUser failed")
	}
	config.mysqlUser = valueString
	// 读取密码
	valueString, err = cfg.GetValue("mysql", "mysqlPassword")
	if err != nil {
		Xlog.Criticalf("Can't get the key value(%s), failed:%s", "mysqlPassword", err)
		return  config, errors.New("GetValue mysqlPassword failed")
	}
	config.mysqlPassword = valueString
	// 读取db名
	valueString, err = cfg.GetValue("mysql", "mysqlDatabaseName")
	if err != nil {
		Xlog.Criticalf("Can't get the key value(%s), failed:%s", "mysqlDatabaseName", err)
		return  config, errors.New("GetValue mysqlDatabaseName failed")
	}
	config.mysqlDatabaseName = valueString
	// 读取db所在地址
	valueString, err = cfg.GetValue("mysql", "mysqlUrl")
	if err != nil {
		Xlog.Criticalf("Can't get the key value(%s), failed:%s", "mysqlUrl", err)
		return  config, errors.New("GetValue mysqlUrl failed")
	}
	config.mysqlUrl = valueString

	Xlog.Info("websocket")
	// 读取websocket配置
	// 读取 wss
	valueString, err = cfg.GetValue("websocket", "wss")
	if err != nil {
		Xlog.Criticalf("Can't get the key value(%s)：%s", "wss", err)
		return config, errors.New("GetValue wss failed")
	}
	valueInt, err = strconv.Atoi(valueString)
	if err != nil {
		Xlog.Criticalf("无法转换%s成int类型：%s", valueString, err)
		return config, errors.New("Atoi wss to int failed")
	}
	config.wss = valueInt
	// 读取wsListenIp
	valueString, err = cfg.GetValue("websocket", "wsListenIp")
	if err != nil {
		Xlog.Criticalf("Can't get the key value(%s)：%s", "wsListenIp", err)
		return  config, errors.New("GetValue wsListenIp failed")
	}
	config.wsListenIp = valueString
	// 读取 wsListenPort
	valueString, err = cfg.GetValue("websocket", "wsListenPort")
	if err != nil {
		Xlog.Criticalf("Can't get the key value(%s)：%s", "wsListenPort", err)
		return config, errors.New("GetValue wsListenPort failed")
	}
	valueInt, err = strconv.Atoi(valueString)
	if err != nil {
		Xlog.Criticalf("无法转换%s成int类型：%s", valueString, err)
		return config, errors.New("Atoi wsListenPort to int  failed")
	}
	config.wsListenPort = valueInt
	// 读取wsKeyPath
	valueString, err = cfg.GetValue("websocket", "wssKeyPath")
	if err != nil {
		Xlog.Criticalf("Can't get the key value(%s)：%s", "wssKeyPath", err)
		return  config, errors.New("GetValue wssKeyPath failed")
	}
	config.wssKeyPath = valueString
	// 读取wsKeyPath
	valueString, err = cfg.GetValue("websocket", "wssCrtPath")
	if err != nil {
		Xlog.Criticalf("Can't get the key value(%s)：%s", "wssCrtPath", err)
		return  config, errors.New("GetValue wssCrtPath failed")
	}
	config.wssCrtPath = valueString

	Xlog.Info("ice")
	err = readIceConfiguration(&config, cfg)
	if err != nil {
		Xlog.Error("readIceConfiguration failed:", err.Error())
	}
	return  config, err
}


var configPath = flag.String("config","config.ini", "input config file path, ex. --config /etc/rtc/roomserver/config.ini")
//1. 解析参数输入，配置文件路径
//2. 读取配置文件

func main() {
	defer Xlog.Flush()

	log.Println("Start RTC Room server ..........")

	if SetupLogger("seelog_config.xml") != 0 {
		log.Fatal("init seelog failed, please check it... ")
		return
	}
	
	flag.Parse()
	Xlog.Infof("configPath = %s", *configPath)
	log.Println("readConfiguration ..........")
	config, err := readConfiguration(*configPath)
	if err != nil {
		Xlog.Criticalf("readConfiguration failed: %s, ", err.Error())
		Xlog.Flush()
		return
	}

	logic.ClientKeepLiveTimeout =  time.Duration(config.clientKeepLiveTimeout)
	iceBalance.IceKeepLiveTimeOut = time.Duration(config.iceKeepLiveTimeout)
	model.ClientReportStatsInterval = config.clientReportStatsInterval
	logic.MaxRoomCapacity = config.maxRoomCapacity

	Xlog.Infof("util.GeoIpMapNew")
	log.Println("util.GeoIpMapNew ..........")
	util.GeoIpTable = util.GeoIpMapNew("GeoLite2-City.mmdb")
	if util.GeoIpTable == nil {
		Xlog.Criticalf("GeoIpMapNew failed")
		return
	}

	log.Println("iceBalance.NewIceTable() ..........")
	Xlog.Infof("iceBalance.NewIceTable()")
	iceBalance.IceTab = iceBalance.NewIceTable()
	for key, iceConf := range config.iceConfs {
		key = key	// 仅是防止报错
		iceBalance.IceTab.CreateIce(iceConf)
	}


	Xlog.Infof("websocket wss:%d, listenIp:%s, listPort:%d, key:%s, crt:%s",
		config.wss, config.wsListenIp,config.wsListenPort, config.wssKeyPath, config.wssCrtPath)
	// 连接数据库
	err = db.ConnectDB(config.mysqlUser, config.mysqlPassword, config.mysqlDatabaseName, config.mysqlUrl)
	if err != nil {
		Xlog.Criticalf("Connect database failed：", err.Error())
		return
	}

	c := connect.NewCollider("0voice")
	c.Run(config.wsListenPort, false)
}
