package model

import (
	"errors"
	Xlog "github.com/cihub/seelog"
	"github.com/json-iterator/go"
)

/**
offer和ReOffer结构体一样，只是命令不同, 但为了结构体各封装函数风格统一，所以还是分开设置
 */
type OfferCmd struct {
	Cmd       string `json:"cmd"`
	AppID     string `json:"appId"`
	RoomID    string `json:"roomId"`
	UID       string `json:"uid"`
	Uname     string `json:"uname"`
	RemoteUID string `json:"remoteUid"`
	IsIceReset bool  `json:"isIceReset"`
	Time      string `json:"time"`
	Msg interface{} `json:"msg"`
}

type OfferRelayCmd struct {
	Cmd       string `json:"cmd"`
	UID       string `json:"uid"`
	Uname     string `json:"uname"`
	RemoteUID string `json:"remoteUid"`
	IsIceReset bool  `json:"isIceReset"`
	SubSessionId  string `json:"subSessionId"`
	Msg interface{} `json:"msg"`
	RtcConfig RtcConfiguration `json:"rtcConfig"`
}

type OfferRespCmd struct {
	Cmd       string `json:"cmd"`
	RemoteUID string `json:"remoteUid"`
	Result   int    `json:"result"`
	Desc     string `json:"desc"`
	SubSessionId  string `json:"subSessionId"`
}

/**
功能：将序列化的数据转换成结构体数据
data 序列化的数据
 */
func OfferUnMarshal(data [] byte, offer *OfferCmd) ( error){
	err :=jsoniter.Unmarshal(data, offer)

	if err != nil {
		return  err
	}

	if offer.RemoteUID == "" {
		return errors.New("RemoteUID is null")
	}

	if offer.UID == "" {
		return errors.New("UID is null")
	}
		var msg SdpMsg
	if offer.Msg == msg{
		return  errors.New("sdp msg is null")
	}

	return  err
}

// 打印结构体信息
func OfferPrint(offer *OfferCmd) {
	Xlog.Info("OfferCmd Cmd:", offer.Cmd, ", AppID:",offer.AppID,
		", RoomID:", offer.RoomID, ", UID:", offer.UID, ", Uname:", offer.Uname,
		", RemoteUID:", offer.RemoteUID, ", SdpMsg.sdp:", offer.Msg)
}

func OfferRelayGenerate(offer *OfferCmd, rtcConfig RtcConfiguration, subSessionId string) (OfferRelayCmd){
	var relayOffer OfferRelayCmd
	relayOffer.Cmd = CMD_RELAY_OFFER
	relayOffer.UID = offer.UID
	relayOffer.Uname = offer.Uname
	relayOffer.RemoteUID = offer.RemoteUID
	relayOffer.SubSessionId = subSessionId
	relayOffer.Msg = offer.Msg
	relayOffer.RtcConfig = rtcConfig

	return relayOffer
}

func OfferRelayMarshal(relayOffer *OfferRelayCmd) ([] byte, error){
	json,err := jsoniter.Marshal(relayOffer)
	return json, err
}

func OfferRelayGenerateAndMarshal(offer *OfferCmd, rtcConfig RtcConfiguration, subSessionId string) ([] byte, error){
	relayOffer := OfferRelayGenerate(offer, rtcConfig, subSessionId)
	return OfferRelayMarshal(&relayOffer)
}


func OfferRespGenerate(offer *OfferCmd, result int, desc string, subSessionId string) (OfferRespCmd){
	var respOffer OfferRespCmd
	respOffer.Cmd = CMD_RESP_OFFER
	respOffer.RemoteUID = offer.RemoteUID
	respOffer.Result = result
	respOffer.Desc = desc
	respOffer.SubSessionId = subSessionId

	return respOffer
}

func OfferRespMarshal(respOffer *OfferRespCmd) ([] byte, error){
	json,err := jsoniter.Marshal(respOffer)
	return json, err
}

func OfferRespGenerateAndMarshal(offer *OfferCmd, result int, desc string, subSessionId string) ([] byte, error){
	respOffer := OfferRespGenerate(offer, result, desc, subSessionId)
	return OfferRespMarshal(&respOffer)
}

func OfferRespErrorAndMarshal(result int, desc string,) ([] byte, error){
	var respOffer OfferRespCmd
	respOffer.Cmd = CMD_RESP_OFFER
	respOffer.Result = result
	respOffer.Desc = desc

	json,err := jsoniter.Marshal(respOffer)
	return json, err
}

