package model

import (
	"errors"
	"github.com/json-iterator/go"
	Xlog "github.com/cihub/seelog"
)
type PeerConnectedCmd struct {
	Cmd       	string `json:"cmd"`
	AppID     	string `json:"appId"`
	ConnectType	string `json:"connectType"`
	RoomID    	string `json:"roomId"`
	UID       	string `json:"uid"`
	Uname    	string `json:"uname"`
	RemoteUID 	string `json:"remoteUid"`
	Time      	string `json:"time"`
}

type PeerConnecteRespCmd struct {
	Cmd       	string `json:"cmd"`
	UID       	string `json:"uid"`
	Uname     	string `json:"uname"`
	RemoteUID 	string `json:"remoteUid"`
	Result   	int    `json:"result"`
	Desc     	string `json:"desc"`
}
/**
功能：将序列化的数据转换成结构体数据
data 序列化的数据
 */
func PeerConnecteUnMarshal(data [] byte, peerConnected *PeerConnectedCmd) ( error){
	err :=jsoniter.Unmarshal(data, peerConnected)

	if err != nil {
		return  err
	}

	if peerConnected.RemoteUID == "" {
		return errors.New("RemoteUID is null")
	}

	if peerConnected.UID == "" {
		return errors.New("UID is null")
	}

	return  err
}

// 打印结构体信息
func PeerConnectePrint(peerConnected *PeerConnectedCmd) {
	Xlog.Info("PeerConnectedCmd Cmd:", peerConnected.Cmd, ", AppID:",peerConnected.AppID,
		", RoomID:", peerConnected.RoomID, ", UID:", peerConnected.UID,
		", RemoteUID:", peerConnected.RemoteUID)
}

func PeerConnecteRespGenerate(peerConnected *PeerConnectedCmd, result int, desc string) (PeerConnecteRespCmd){
	var respPeerConnecte PeerConnecteRespCmd
	respPeerConnecte.Cmd = CMD_RESP_PEER_CONNECTED
	respPeerConnecte.UID = peerConnected.UID
	respPeerConnecte.RemoteUID = peerConnected.RemoteUID
	respPeerConnecte.Result = result
	respPeerConnecte.Desc = desc

	return respPeerConnecte
}

func PeerConnecteRespMarshal(respPeerConnecte *PeerConnecteRespCmd) ([] byte, error){
	json,err := jsoniter.Marshal(respPeerConnecte)
	return json, err
}

func PeerConnecteRespGenerateAndMarshal(peerConnected *PeerConnectedCmd, result int, desc string) ([] byte, error){
	respPeerConnecte := PeerConnecteRespGenerate(peerConnected, result, desc)
	return PeerConnecteRespMarshal(&respPeerConnecte)
}