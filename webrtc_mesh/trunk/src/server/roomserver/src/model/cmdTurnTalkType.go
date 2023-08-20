package model

import (
	"errors"
	"github.com/json-iterator/go"
	Xlog "github.com/cihub/seelog"
)

/**
turnTalkType和ReTurnTalkType结构体一样，只是命令不同, 但为了结构体各封装函数风格统一，所以还是分开设置
 */

type TurnTalkTypeCmd struct {
	Cmd       string `json:"cmd"`
	AppID     string `json:"appId"`
	RoomID    string `json:"roomId"`
	UID       string `json:"uid"`
	Uname       string `json:"uname"`
	Index	  int 	  `json:"index"`	// 设备类型，0摄像头，1麦克风，2共享屏幕，3系统声音
	Enable    bool	  `json:"enable"`  // false关闭，true开启
	Time      string `json:"time"`
}

type TurnTalkTypeRelayCmd struct {
	Cmd       string `json:"cmd"`
	UID       string `json:"uid"`
	Uname     string `json:"uname"`
	Index	  int 	  `json:"index"`	// 设备类型，0摄像头，1麦克风，2共享屏幕，3系统声音
	Enable    bool	  `json:"enable"`  // false关闭，true开启
}

type TurnTalkTypeRespCmd struct {
	Cmd       string `json:"cmd"`
	UID       string `json:"uid"`
	Result   int    `json:"result"`
	Desc     string `json:"desc"`
}

/**
功能：将序列化的数据转换成结构体数据
data 序列化的数据
 */
func TurnTalkTypeUnMarshal(data [] byte, turnTalkType *TurnTalkTypeCmd) ( error){
	err :=jsoniter.Unmarshal(data, turnTalkType)

	if err != nil {
		return  err
	}

	if turnTalkType.UID == "" {
		return errors.New("UID is null")
	}

	return  err
}

// 打印结构体信息
func TurnTalkTypePrint(turnTalkType *TurnTalkTypeCmd) {
	Xlog.Info("TurnTalkTypeCmd Cmd:", turnTalkType.Cmd, ", AppID:",turnTalkType.AppID,
		", RoomID:", turnTalkType.RoomID, ", UID:", turnTalkType.UID, "Index:", turnTalkType.Index, ", Enable:", turnTalkType.Enable)
}

func TurnTalkTypeRelayGenerate(turnTalkType *TurnTalkTypeCmd) (TurnTalkTypeRelayCmd){
	var relayTurnTalkType TurnTalkTypeRelayCmd
	relayTurnTalkType.Cmd = CMD_RELAY_TURN_TALK_TYPE
	relayTurnTalkType.UID = turnTalkType.UID
	relayTurnTalkType.Uname = turnTalkType.Uname
	relayTurnTalkType.Index = turnTalkType.Index
	relayTurnTalkType.Enable = turnTalkType.Enable

	return relayTurnTalkType
}

func TurnTalkTypeRelayMarshal(relayTurnTalkType *TurnTalkTypeRelayCmd) ([] byte, error){
	json,err := jsoniter.Marshal(relayTurnTalkType)
	return json, err
}

func TurnTalkTypeRelayGenerateAndMarshal(turnTalkType *TurnTalkTypeCmd) ([] byte, error){
	relayTurnTalkType := TurnTalkTypeRelayGenerate(turnTalkType)
	return TurnTalkTypeRelayMarshal(&relayTurnTalkType)
}


func TurnTalkTypeRespGenerate(turnTalkType *TurnTalkTypeCmd, result int, desc string) (TurnTalkTypeRespCmd){
	var respTurnTalkType TurnTalkTypeRespCmd
	respTurnTalkType.Cmd = CMD_RESP_TURN_TALK_TYPE
	respTurnTalkType.UID = turnTalkType.UID
	respTurnTalkType.Result = result
	respTurnTalkType.Desc = desc

	return respTurnTalkType
}

func TurnTalkTypeRespMarshal(respTurnTalkType *TurnTalkTypeRespCmd) ([] byte, error){
	json,err := jsoniter.Marshal(respTurnTalkType)
	return json, err
}

func TurnTalkTypeRespGenerateAndMarshal(turnTalkType *TurnTalkTypeCmd, result int, desc string) ([] byte, error){
	respTurnTalkType := TurnTalkTypeRespGenerate(turnTalkType, result, desc)
	return TurnTalkTypeRespMarshal(&respTurnTalkType)
}