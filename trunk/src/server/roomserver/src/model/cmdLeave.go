package model

import (
	"errors"
	"github.com/json-iterator/go"
	Xlog "github.com/cihub/seelog"
)

type LeaveCmd struct {
	Cmd       string `json:"cmd"`
	AppID     string `json:"appId"`
	RoomID    string `json:"roomId"`
	UID       string `json:"uid"`
	Uname    string `json:"uname"`
	UserType int   	`json:"userType"`
	Time      string `json:"time"`
}

type LeaveRelayCmd struct {
	Cmd       string `json:"cmd"`
	UID       string `json:"uid"`
	Uname    string `json:"uname"`
	UserType int   	`json:"userType"`
}

type LeaveRespCmd struct {
	Cmd       string `json:"cmd"`
	Result   int    `json:"result"`
	Desc     string `json:"desc"`
}

/**
功能：将序列化的数据转换成结构体数据
data 序列化的数据
 */
func LeaveUnMarshal(data [] byte, leave *LeaveCmd) ( error){
	err :=jsoniter.Unmarshal(data, leave)
	if err != nil {
		return  err
	}

	if leave.AppID == "" {
		return errors.New("AppID is null")
	}

	if leave.UID == "" {
		return errors.New("UID is null")
	}

	if leave.Time == "" {
		return errors.New("Time is null")
	}

	return  err
}

// 打印结构体信息
func LeavePrint(leave *LeaveCmd) {
	Xlog.Info("LeaveCmd Cmd:", leave.Cmd, ", AppID:",leave.AppID, ", RoomID:", leave.RoomID, ", UID:", leave.UID)
}

// 根据leave请求生成转发
func LeaveRelayGenerate(leave *LeaveCmd) (LeaveRelayCmd){
	var relayLeave LeaveRelayCmd
	relayLeave.Cmd = CMD_RELAY_LEAVE
	relayLeave.UID = leave.UID
	relayLeave.Uname = leave.Uname
	relayLeave.UserType = leave.UserType
	return relayLeave
}

// 序列化转发leave
func LeaveRelayMarshal(relayLeave *LeaveRelayCmd) ([] byte, error){
	json,err := jsoniter.Marshal(relayLeave)
	return json, err
}

func LeaveRelayGenerateAndMarshal(leave *LeaveCmd) ([] byte, error){
	relayLeave := LeaveRelayGenerate(leave)
	return LeaveRelayMarshal(&relayLeave)
}

// 根据leave请求生成resp响应
func LeaveRespGenerate(leave *LeaveCmd, ret int, desc string) (LeaveRespCmd){
	var respLeave LeaveRespCmd
	respLeave.Cmd = CMD_RESP_LEAVE
	respLeave.Result = ret
	respLeave.Desc =desc
	return respLeave
}

// 序列化响应leave
func LeaveRespMarshal(respLeave *LeaveRespCmd) ([] byte, error){
	json,err := jsoniter.Marshal(respLeave)
	return json, err
}

func LeaveRespGenerateAndpMarshal(leave *LeaveCmd, ret int, desc string) ([] byte, error){
	respLeave := LeaveRespGenerate(leave, ret, desc)
	return LeaveRespMarshal(&respLeave)
}

// 生成错误原因，并做json序列化
func LeaveRespErrorAndMarshal(ret int, desc string) ([] byte, error){
	var respLeave LeaveRespCmd
	respLeave.Cmd = CMD_RESP_LEAVE
	respLeave.Result = ret
	respLeave.Desc = desc
	json,err := jsoniter.Marshal(respLeave)
	return json, err
}