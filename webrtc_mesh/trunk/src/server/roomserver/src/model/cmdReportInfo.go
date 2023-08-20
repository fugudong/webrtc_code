package model

import (
	"github.com/json-iterator/go"
	"errors"
)
type ReportInfoCmd struct {
	Cmd       string `json:"cmd"`
	AppID     string `json:"appId"`
	RoomID    string `json:"roomId"`
	UID       string `json:"uid"`
	Uname       string `json:"uname"`
	Result	int    `json:"result"`		// 返回值，0为正常，非0则为异常
	Desc 	string `json:"desc"`	// 详细的结果信息
	Data1	string `json:"data1"`	// 可变参数1，比如uid，客户端收到消息时根据自己命令的意义进行解析
	Data2 	string `json:"data2"`	// 可变参数2
	Time      string `json:"time"`
}

func ReportInfoUnMarshal(data [] byte, reportInfo *ReportInfoCmd) ( error){
	err :=jsoniter.Unmarshal(data, reportInfo)

	if err != nil {
		return  err
	}

	if reportInfo.RoomID == "" {
		return errors.New("RoomID is null")
	}

	if reportInfo.UID == "" {
		return errors.New("UID is null")
	}

	return  err
}