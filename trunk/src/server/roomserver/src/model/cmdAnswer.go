package model

import (
	"errors"
	"github.com/json-iterator/go"
	Xlog "github.com/cihub/seelog"
)

type AnswerCmd struct {
	Cmd       string `json:"cmd"`
	AppID     string `json:"appId"`
	ConnectType      int    `json:"connectType"`
	RoomID    string `json:"roomId"`
	UID       string `json:"uid"`
	Uname     string `json:"uname"`
	RemoteUID string `json:"remoteUid"`
	Time      string `json:"time"`
	//Iceserver IceServer `json:"iceserver"`
	//Msg SdpMsg `json:"msg"`
	Msg interface{} `json:"msg"`
}

type AnswerRelayCmd struct {
	Cmd       string `json:"cmd"`
	ConnectType      int    `json:"connectType"`
	UID       string `json:"uid"`
	Uname     string `json:"uname"`
	RemoteUID string `json:"remoteUid"`
	//Iceserver IceServer `json:"iceserver"`
	//Msg SdpMsg `json:"msg"`
	Msg interface{} `json:"msg"`
}

type AnswerRespCmd struct {
	Cmd       string `json:"cmd"`
	UID       string `json:"uid"`
	Uname     string `json:"uname"`
	RemoteUID string `json:"remoteUid"`
	Result   int    `json:"result"`
	Desc     string `json:"desc"`
}

/**
功能：将序列化的数据转换成结构体数据
data 序列化的数据
 */
func AnswerUnMarshal(data [] byte, answer *AnswerCmd) ( error){
	err :=jsoniter.Unmarshal(data, answer)

	if err != nil {
		return  err
	}

	if answer.RemoteUID == "" {
		return errors.New("RemoteUID is null")
	}

	if answer.UID == "" {
		return errors.New("UID is null")
	}

	var msg SdpMsg
	if answer.Msg == msg{
		return  errors.New("sdp msg is null")
	}

	return  err
}

// 打印结构体信息
func AnswerPrint(answer *AnswerCmd) {
	Xlog.Info("AnswerCmd Cmd:", answer.Cmd, ", AppID:",answer.AppID, ", ConnectType:", answer.ConnectType, "," +
		", RoomID:", answer.RoomID, ", UID:", answer.UID, ", Uname:", answer.Uname,
		", RemoteUID:", answer.RemoteUID, ", SdpMsg.sdp:", answer.Msg)
}

func AnswerRelayGenerate(answer *AnswerCmd) (AnswerRelayCmd){
	var relayAnswer AnswerRelayCmd
	relayAnswer.Cmd = CMD_RELAY_ANSWER
	relayAnswer.ConnectType = answer.ConnectType
	relayAnswer.UID = answer.UID
	relayAnswer.Uname = answer.Uname
	relayAnswer.RemoteUID = answer.RemoteUID
	relayAnswer.Msg = answer.Msg

	return relayAnswer
}

func AnswerRelayMarshal(relayAnswer *AnswerRelayCmd) ([] byte, error){
	json,err := jsoniter.Marshal(relayAnswer)
	return json, err
}

func AnswerRelayGenerateAndMarshal(answer *AnswerCmd) ([] byte, error){
	relayAnswer := AnswerRelayGenerate(answer)
	return  AnswerRelayMarshal(&relayAnswer)
}

func AnswerRespGenerate(answer *AnswerCmd, result int, desc string) (AnswerRespCmd){
	var respAnswer AnswerRespCmd
	respAnswer.Cmd = CMD_RESP_ANSWER
	respAnswer.UID = answer.UID
	respAnswer.Uname = answer.Uname
	respAnswer.RemoteUID = answer.RemoteUID
	respAnswer.Result = result
	respAnswer.Desc = desc

	return respAnswer
}

func AnswerRespMarshal(respAnswer *AnswerRespCmd) ([] byte, error){
	json,err := jsoniter.Marshal(respAnswer)
	return json, err
}

func AnswerRespGenerateAndMarshal(answer *AnswerCmd, result int, desc string) ([] byte, error){
	answerOffer := AnswerRespGenerate(answer, result, desc)
	return  AnswerRespMarshal(&answerOffer)
}