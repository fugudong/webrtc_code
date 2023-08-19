package model

import (
	"errors"
	"github.com/json-iterator/go"
	Xlog "github.com/cihub/seelog"
)

type KeepLiveCmd struct {
	Cmd       string `json:"cmd"`
	AppID     string `json:"appId"`
	RoomID    string `json:"roomId"`
	UID       string `json:"uid"`
	Time      string `json:"time"`
}

type KeepLiveRespCmd struct {
	Cmd       string `json:"cmd"`
	Result   int    `json:"result"`
	Desc     string `json:"desc"`
}

/**
功能：将序列化的数据转换成结构体数据
data 序列化的数据
 */
func KeepLiveUnMarshal(data [] byte, keepLive *KeepLiveCmd) ( error){
	err :=jsoniter.Unmarshal(data, keepLive)
	if err != nil {
		return  err
	}

	if keepLive.AppID == "" {
		return errors.New("AppID is null")
	}

	if keepLive.UID == "" {
		return errors.New("UID is null")
	}

	if keepLive.Time == "" {
		return errors.New("Time is null")
	}

	return  err
}

// 打印结构体信息
func KeepLivePrint(keepLive *KeepLiveCmd) {
	Xlog.Info("KeepLiveCmd Cmd:", keepLive.Cmd, ", AppID:",keepLive.AppID, ", RoomID:", keepLive.RoomID, ", UID:", keepLive.UID)
}


// 根据keepLive请求生成resp响应
func KeepLiveRespGenerate(keepLive *KeepLiveCmd, result int, desc string) (KeepLiveRespCmd){
	var respKeepLive KeepLiveRespCmd
	respKeepLive.Cmd = CMD_RESP_KEEP_LIVE
	respKeepLive.Result = result
	respKeepLive.Desc = desc
	return respKeepLive
}

// 序列化响应keepLive
func KeepLiveRespMarshal(respKeepLive *KeepLiveRespCmd) ([] byte, error) {
	json, err := jsoniter.Marshal(respKeepLive)
	return json, err
}

func KeepLiveRespGenerateAndpMarshal(keepLive *KeepLiveCmd, ret int, desc string) ([] byte, error){
	respKeepLive := KeepLiveRespGenerate(keepLive, ret, desc)
	return KeepLiveRespMarshal(&respKeepLive)
}