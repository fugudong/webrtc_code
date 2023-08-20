package main

import (
	"./db"
	"./util"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/Unknwon/goconfig"
	log "github.com/cihub/seelog"
	"html/template"
	"time"
)

type TalkReportConfig struct {
	// mysql数据库
	mysqlConfig db.MysqlConfig
	mailConfig	util.MailInfo
}

var talkReportconfig *TalkReportConfig

type DailyTalkReport struct {
	Title 	string
	NowDate string	// 当前日期
	Note string	// 备注
	TalkReport [] *db.TalkSessionReprt
}


//定时创建数据库

func loopReportTalkStats(done chan int) {
	defer func() { done <- 1 } ()	// 退出函数时发送channel通知阻塞的主线程退出程序
	for {
		now := time.Now()//获取当前时间，放到now里面，要给next用  
		next := now.Add(time.Hour * 24) //通过now偏移24小时
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location()) //获取下一个凌晨的日期
		// Format("2006-01-02 15:04:05")
		theDate :=  fmt.Sprintf("%04d%02d%02d",  now.Year(), now.Month(), now.Day())
		log.Infof("Now date:%s, Next date:%s， dif_time:%v, the_date:%s",
			now.Format("2006-01-02 15:04:05"),
			next.Format("2006-01-02 15:04:05"),
			next.Sub(now),
			theDate)
		t := time.NewTimer(next.Sub(now))//计算当前时间到凌晨的时间间隔，设置一个定时器
		//t := time.NewTimer(5 * time.Second)//计算当前时间到凌晨的时间间隔，设置一个定时器

		<-t.C		// 等待时间到达

		// 连接数据库
		err := db.ConnectDB(&talkReportconfig.mysqlConfig)
		if err != nil {
			log.Errorf("-->ConnectDB:%s", err.Error())
			return
		}
		defer db.DisconnectDB()
		// 读取数据库

		count, err := db.TalkSessionGetCountByDate(theDate)
		if err != nil {
			log.Errorf("-->TalkSessionGetCountByDate:%s", err.Error())
			return
		}
		log.Infof("TalkSessionGetCountByDate count:%d", count)
		talkList, err := db.TalkSessionsList(theDate, 0, count)
		if err != nil {
			log.Errorf("-->TalkSessionsList:%s", err.Error())
		}
		var dailyReport DailyTalkReport
		dailyReport.TalkReport = talkList
		log.Infof("dailyReport.TalkReport = %d", len(dailyReport.TalkReport))

		dailyReport.Title = fmt.Sprintf("音视频通话统计日报 %04d-%02d-%02d", now.Year(), now.Month(), now.Day())
		dailyReport.NowDate = time.Now().Format("2006-01-02 15:04:05")
		tmpl := template.Must(template.ParseFiles("index.html"))
		b := new(bytes.Buffer)
		tmpl.Execute(b,  struct{ DailyTalkReport DailyTalkReport }{dailyReport} )
		result := string(b.String())
		sendMail := &talkReportconfig.mailConfig
		sendMail.Subject = dailyReport.Title
		err = sendMail.SendMail(result)
		if err != nil {
			log.Errorf("SendMail failed:%s", err.Error())
		}
		break
	}
}

func SetupLogger(logConfig string)  int{
	logger, err := log.LoggerFromConfigAsFile(logConfig)
	if err != nil {
		log.Critical("LoggerFromConfigAsFile failed:%s\n", err)
		return -1
	}

	log.ReplaceLogger(logger)

	log.Infof("SetupLogger ok\n",)
	return 0
}

//1. 读取配置文件
func readConfiguration(filePath string)  (*TalkReportConfig, error){
	var config TalkReportConfig

	log.Info("LoadConfigFile")
	cfg, err := goconfig.LoadConfigFile(filePath)
	if err != nil {
		log.Critical("Read config %s failed: %s", filePath, err.Error())
		return  &config, errors.New("LoadConfigFile failed")
	}

	var valueString string

	log.Info("mysql")
	// 读取mysql配置
	// 读取用户名
	valueString, err = cfg.GetValue("mysql", "userName")
	if err != nil {
		log.Criticalf("Can't get the key value(%s), failed:%s", "userName", err)
		return  &config, errors.New("GetValue userName failed")
	}
	config.mysqlConfig.UserName = valueString
	// 读取密码
	valueString, err = cfg.GetValue("mysql", "password")
	if err != nil {
		log.Criticalf("Can't get the key value(%s), failed:%s", "password", err)
		return  &config, errors.New("GetValue password failed")
	}
	config.mysqlConfig.Password = valueString
	// 读取db名
	valueString, err = cfg.GetValue("mysql", "databaseName")
	if err != nil {
		log.Criticalf("Can't get the key value(%s), failed:%s", "databaseName", err)
		return  &config, errors.New("GetValue databaseName failed")
	}
	config.mysqlConfig.DatabaseName = valueString
	// 读取db所在地址
	valueString, err = cfg.GetValue("mysql", "url")
	if err != nil {
		log.Criticalf("Can't get the key value(%s), failed:%s", "url", err)
		return  &config, errors.New("GetValue url failed")
	}
	config.mysqlConfig.Url = valueString

	log.Info("mail")
	// 读取mail配置
	// 读取 服务器外网ip地址
	valueString, err = cfg.GetValue("mail", "ipAddr")
	if err != nil {
		log.Criticalf("Can't get the key value(%s)：%s", "ipAddr", err)
		return &config, errors.New("GetValue ipAddr failed")
	}

	config.mailConfig.SendIp = valueString
	// 读取邮件收件人
	valueString, err = cfg.GetValue("mail", "receivers")
	if err != nil {
		log.Criticalf("Can't get the key value(%s)：%s", "receivers", err)
		return  &config, errors.New("GetValue receivers failed")
	}
	config.mailConfig.Receivers = valueString
	// 读取 sender发送者
	valueString, err = cfg.GetValue("mail", "sender")
	if err != nil {
		log.Criticalf("Can't get the key value(%s)：%s", "sender", err)
		return &config, errors.New("GetValue sender failed")
	}
	config.mailConfig.Sender = valueString
	// 读取 发送人登录用户名
	valueString, err = cfg.GetValue("mail", "user")
	if err != nil {
		log.Criticalf("Can't get the key value(%s)：%s", "user", err)
		return &config, errors.New("GetValue user failed")
	}
	config.mailConfig.User = valueString
	// 读取邮件发件人邮箱密码
	valueString, err = cfg.GetValue("mail", "passwd")
	if err != nil {
		log.Criticalf("Can't get the key value(%s)：%s", "passwd", err)
		return  &config, errors.New("GetValue passwd failed")
	}
	config.mailConfig.Passwd = valueString

	// 读取邮件发件人邮箱密码
	valueString, err = cfg.GetValue("mail", "smtpHostAndPort")
	if err != nil {
		log.Criticalf("Can't get the key value(%s)：%s", "smtpHostAndPort", err)
		return  &config, errors.New("GetValue smtpHostAndPort failed")
	}
	config.mailConfig.SmtpHostAndPort = valueString

	// 读取邮件发件人邮箱密码
	valueString, err = cfg.GetValue("mail", "subject")
	if err != nil {
		log.Criticalf("Can't get the key value(%s)：%s", "subject", err)
		return  &config, errors.New("GetValue subject failed")
	}
	config.mailConfig.Subject = valueString

	return  &config, err
}


const DEFAULT_DATE = "20190527"
var configPath = flag.String("config","conf.ini", "input config file path, ex. --config /etc/rtc/roomserver/config.conf")
func main()  {
	defer log.Flush()

	fmt.Println("Start Talk Report server ..........")

	if SetupLogger("seelog_config.xml") != 0 {
		log.Critical("init seelog failed, please check it... ")
		return
	}

	config, err := readConfiguration(*configPath)
	if err != nil {
		log.Criticalf("readConfiguration failed: %s, ", err.Error())
		fmt.Printf("readConfiguration failed: %s, ", err.Error())
		log.Flush()
		return
	}
	talkReportconfig = config

	util.GeoIpTable = util.GeoIpMapNew("GeoLite2-City.mmdb")
	if util.GeoIpTable == nil {
		log.Errorf("GeoIpMapNew failed")
		return
	}
	c := make(chan int)

	go loopReportTalkStats(c)
	log.Infof("wait Talk Report eixt")
	<-c	// 使用channel阻塞主线程

	log.Error("Talk Report eixt")	//正常运行中不应该退出程序，如果退出了则是出错了，需要处理
}