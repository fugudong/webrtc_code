package model

import (
	"errors"
	"github.com/json-iterator/go"
	Xlog "github.com/cihub/seelog"
)

type CandidateMsg  struct {
	Candidate   string `json:"candidate"`
	CandidateID string `json:"candidateId"`
	Label       int    `json:"label"`
}
type CandidateCmd struct {
	Cmd       string `json:"cmd"`
	AppID     string `json:"appId"`
	RoomID    string `json:"roomId"`
	UID       string `json:"uid"`
	RemoteUID string `json:"remoteUid"`
	Time      string `json:"time"`
	Msg interface{}  `json:"msg"`
}

type CandidateRelayCmd struct {
	Cmd       string `json:"cmd"`
	UID       string `json:"uid"`
	Uname     string `json:"uname"`
	RemoteUID string `json:"remoteUid"`
	Msg interface{}  `json:"msg"`
}

type CandidateRespCmd struct {
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
func CandidateUnMarshal(data [] byte, candidate *CandidateCmd) ( error){
	err :=jsoniter.Unmarshal(data, candidate)

	if err != nil {
		return  err
	}

	if candidate.RemoteUID == "" {
		return errors.New("RemoteUID is null")
	}

	if candidate.UID == "" {
		return errors.New("UID is null")
	}
	var msg CandidateMsg
	if candidate.Msg == msg{
		return  errors.New("ice candidate msg is null")
	}

	return  err
}

// 打印结构体信息
func CandidatePrint(candidate *CandidateCmd) {
	Xlog.Info("CandidateCmd Cmd:", candidate.Cmd, ", AppID:",candidate.AppID,
		", RoomID:", candidate.RoomID, ", UID:", candidate.UID,
		", RemoteUID:", candidate.RemoteUID, ", IceMsg.Candidate:", candidate.Msg)
}

func CandidateRelayGenerate(candidate *CandidateCmd) (CandidateRelayCmd){
	var relayCandidate CandidateRelayCmd
	relayCandidate.Cmd = CMD_RELAY_CANDIDATE
	relayCandidate.UID = candidate.UID
	relayCandidate.RemoteUID = candidate.RemoteUID
	relayCandidate.Msg = candidate.Msg

	return relayCandidate
}

func CandidateRelayMarshal(relayCandidate *CandidateRelayCmd) ([] byte, error){
	json,err := jsoniter.Marshal(relayCandidate)
	return json, err
}

func CandidateRelayGenerateAndMarshal(candidate *CandidateCmd) ([] byte, error){
	relayCandidate := CandidateRelayGenerate(candidate)
	return CandidateRelayMarshal(&relayCandidate)
}

func CandidateRespGenerate(candidate *CandidateCmd, result int, desc string) (CandidateRespCmd){
	var respCandidate CandidateRespCmd
	respCandidate.Cmd = CMD_RESP_CANDIDATE
	respCandidate.UID = candidate.UID
	respCandidate.RemoteUID = candidate.RemoteUID
	respCandidate.Result = result
	respCandidate.Desc = desc

	return respCandidate
}

func CandidateRespMarshal(respCandidate *CandidateRespCmd) ([] byte, error){
	json,err := jsoniter.Marshal(respCandidate)
	return json, err
}

func CandidateRespGenerateAndMarshal(candidate *CandidateCmd, result int, desc string) ([] byte, error){
	respCandidate := CandidateRespGenerate(candidate, result, desc)
	return CandidateRespMarshal(&respCandidate)
}