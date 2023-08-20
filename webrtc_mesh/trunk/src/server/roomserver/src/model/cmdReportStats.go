package model

import (
	"errors"
	"github.com/json-iterator/go"
)

type ReportStatsCmd struct {
	Cmd          	string `json:"cmd"`
	AppID 		  	string `json:"appId"`
	RoomID       	string `json:"roomId"`
	UID          	string `json:"uid"`
	RemoteUID   	string `json:"remoteUid"`
	Audio struct {
		Send struct {
			BitRate         int `json:"bitRate"`
			PacketsLostRate string `json:"packetsLostRate"`
		} `json:"send"`
		Recv struct {
			BitRate         int `json:"bitRate"`
			PacketsLostRate string `json:"packetsLostRate"`
		} `json:"recv"`
	} `json:"audio"`

	Video        struct {
		Recv struct {
			BitRate           int `json:"bitRate"`
			CodecName         string `json:"codecName"`
			FrameRateRecv     int `json:"frameRateRecv"`
			Height            int `json:"height"`
			PacketsLostRate   string `json:"packetsLostRate"`
			Width             int `json:"width"`
		} `json:"recv"`
		Send struct {
			BitRate         int    `json:"bitRate"`
			CodecName       string `json:"codecName"`
			FrameRateSent   int `json:"frameRateSent"`
			Height          int `json:"height"`
			PacketsLostRate string    `json:"packetsLostRate"`
			Width           int `json:"width"`
		} `json:"send"`
	} `json:"video"`
	Time         string `json:"time"`
}

func ReportStatsUnMarshal(data [] byte, reportStats *ReportStatsCmd) ( error){
	err :=jsoniter.Unmarshal(data, reportStats)

	if err != nil {
		return  err
	}

	if reportStats.RoomID == "" {
		return errors.New("RoomID is null")
	}

	if reportStats.UID == "" {
		return errors.New("UID is null")
	}

	if reportStats.RemoteUID == "" {
		return errors.New("RemoteUID is null")
	}

	return  err
}
