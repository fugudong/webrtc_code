package model

import (
	"roomserver/src/util"
	"github.com/json-iterator/go"
)

// join格式 后期需要加上Ice server？

type JoinCmd struct {
	Cmd      string `json:"cmd"`
	AppID    string `json:"appId"`
	Token    string `json:"token"`
	RoomID   string `json:"roomId"`
	RoomName string `json:"roomName"`
	UID      string `json:"uid"`
	Uname    string `json:"uname"`
	UserType int   	`json:"userType"`
	TalkType int   	`json:"talkType"`
	Time     string `json:"time"`		// 时间同样使用string，在处理时转换为int64
	IP       string `json:"ip"`
	OsName   string `json:"osName"`
	Browser   string `json:"browser"`
	SdkInfo  string `json:"sdkInfo"`
}

type JoinNewPeerCmd struct {
	Cmd      string `json:"cmd"`
	UID      string `json:"uid"`
	Uname    string `json:"uname"`
	UserType int   	`json:"userType"`
	TalkType int   	`json:"talkType"`	// 既定的通话类型
	//DevType	int 	`json:"devType"`	// 新人进来使用TalkType即可
	RtcConfig RtcConfiguration `json:"rtcConfig"`
}

type RoomUser struct {
	AppID     string `json:"appId"`
	UID       string `json:"uid"`
	Uname     string `json:"uname"`
	UserType int   	`json:"userType"`
	TalkType int   	`json:"talkType"`
}

type JoinRespCmd struct {
	Cmd      string `json:"cmd"`
	Result   int    `json:"result"`
	Desc     string `json:"desc"`
	RoomID   string `json:"roomId"`
	RoomName string `json:"roomName"`
	UID      string `json:"uid"`
	Uname    string `json:"uname"`
	ReportStatsInterval int `json:"reportStatsInterval"`
	UserList [] RoomUser `json:"userList"`
}
var ClientReportStatsInterval = 30000
/**
功能：将序列化的数据转换成结构体数据
data 序列化的数据
 */
func JoinUnMarshal(data [] byte, join *JoinCmd) ( error){
	err :=jsoniter.Unmarshal(data, join)

	// 生成roomid uid
	if join.RoomID == "" {
		join.RoomID = util.GenerateRoomID()
	}
	if join.UID == "" {
		join.UID = util.GenerateUserID()
	}

	return  err
}

func JoinRespGenerate(join *JoinCmd, users [] RoomUser, ret int, desc string) (JoinRespCmd){
	var respJoin JoinRespCmd
	respJoin.Cmd = CMD_RESP_JOIN
	respJoin.Result = ret
	respJoin.ReportStatsInterval = ClientReportStatsInterval
	respJoin.Desc = desc
	respJoin.RoomID = join.RoomID
	respJoin.RoomName = join.RoomName
	respJoin.UID = join.UID
	respJoin.Uname = join.Uname
	respJoin.UserList = users
	return respJoin
}

// 生成错误原因，并做json序列化
func JoinRespErrorAndMarshal(ret int, desc string) ([] byte, error){
	var respJoin JoinRespCmd
	respJoin.Cmd = CMD_RESP_JOIN
	respJoin.Result = ret
	respJoin.Desc = desc
	json,err := jsoniter.Marshal(respJoin)
	return json, err
}

func JoinRespMarshal(respJoin *JoinRespCmd) ([] byte, error){
	json,err := jsoniter.Marshal(respJoin)
	return json, err
}

func JoinRespGenerateAndMarshal(join *JoinCmd, users [] RoomUser, ret int, desc string) ([] byte, error){
	respJoin := JoinRespGenerate(join, users, ret, desc)
	return JoinRespMarshal(&respJoin)
}

func JoinNewPeerGenerate(join *JoinCmd, rtcConfig RtcConfiguration) (JoinNewPeerCmd){
	var newJoin JoinNewPeerCmd
	newJoin.Cmd = CMD_NOTIFY_NEW_PEER
	newJoin.UID = join.UID
	newJoin.Uname = join.Uname
	newJoin.UserType = join.UserType	// 用户类型
	newJoin.TalkType = join.TalkType	// 通话类型
	newJoin.RtcConfig = rtcConfig
	return newJoin
}

func JoinNewPeerMarshal(newJoin *JoinNewPeerCmd) ([] byte, error){
	json,err := jsoniter.Marshal(newJoin)
	return json, err
}

func JoinNewPeerGenerateAndMarshal(join *JoinCmd, rtcConfig RtcConfiguration) ([] byte, error) {
	newJoin := JoinNewPeerGenerate(join, rtcConfig)
	return JoinNewPeerMarshal(&newJoin)
}

func JoinGenerateLeave(join *JoinCmd)  (LeaveCmd) {
	var leave LeaveCmd
	leave.Cmd = CMD_LEAVE
	leave.AppID = join.AppID
	leave.RoomID = join.RoomID
	leave.UID = join.UID
	leave.Time = join.Time

	return  leave
}




