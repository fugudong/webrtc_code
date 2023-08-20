package model

import "github.com/json-iterator/go"



// 客户端和服务器的交互命令
// 带resp前缀则为服务器回复客户端
// 带relay则为服务器转发客户端的消息给相应的客户端

const CMD_GENERAL_MSG_RESP = "generalMsgResp"

const CMD_JOIN = "join"
const CMD_RESP_JOIN = "respJoin"

const CMD_LEAVE 		= "leave"
const CMD_RESP_LEAVE 	= "respLeave"
const CMD_RELAY_LEAVE 	= "relayLeave" // relay为前缀，前缀后的字段表明是客户端->服务器端的命令

const CMD_OFFER = "offer"
const CMD_RESP_OFFER = "respOffer"
const CMD_RELAY_OFFER = "relayOffer"

const CMD_ANSWER = "answer"
const CMD_RESP_ANSWER = "respAnswer"
const CMD_RELAY_ANSWER = "relayAnswer"

const CMD_CANDIDATE = "candidate"
const CMD_RESP_CANDIDATE = "respCandidate"
const CMD_RELAY_CANDIDATE =  "relayCandidate"

const CMD_KEEP_LIVE = "keepLive"
const CMD_RESP_KEEP_LIVE = "respKeepLive"

const CMD_TURN_TALK_TYPE = "turnTalkType"
const CMD_RESP_TURN_TALK_TYPE = "respTurnTalkType"
const CMD_RELAY_TURN_TALK_TYPE = "relayTurnTalkType"

const CMD_PEER_CONNECTED = "peerConnected" // 通话成功
const CMD_RESP_PEER_CONNECTED = "respPeerConnected"

// 报告状态信息
const CMD_REPORT_INFO	= "reportInfo"
const CMD_REPORT_STATS	= "reportStats"

// 服务器端主动发送消息给客户端
const CMD_NOTIFY_NEW_SESSION = "notifyNewSession" // 客户端需要保存收到的session，在向服务器发送请求时需要使用sessionID
const CMD_NOTIFY_NEW_PEER = "notifyNewPeer"       // 通知其他客户端有新人加入


type SdpMsg struct {
	Sdp  string `json:"sdp"`
	Type string `json:"type"`
}

//  服务器回应客户端的字段，该结构体为通用结构体，一般在服务器在解析客户端的message异常时使用
type GeneralRespMsg struct {
	Cmd    	string `json:"cmd"`		// 回应哪一条命令
	Result	int    `json:"result"`		// 返回值，0为正常，非0则为异常
	Desc 	string `json:"desc"`	// 详细的结果信息
	Data1	string `json:"data1"`	// 可变参数1，比如uid，客户端收到消息时根据自己命令的意义进行解析
	Data2 	string `json:"data2"`	// 可变参数2
}
//
/**
* 生成通用的回应消息，并做json序列化
* cmd 回复的命令类型，比如CmdJoin，CmdOffer等等
* ret 返回码，0为正常，其他值则为出错值
* desc 文字描述ret的意义
* data1 不同的命令意义不同
* data2 不同的命令意义不同
 */
func GeneralRespMsgMarshal(cmd string, ret int, desc string,
	data1 string, data2 string) ([] byte, error){
	var respMsg GeneralRespMsg
	respMsg.Cmd = CMD_GENERAL_MSG_RESP
	respMsg.Result = ret
	respMsg.Desc = desc
	respMsg.Data1 = data1
	respMsg.Data2 = data2
	respBytes,err := jsoniter.Marshal(respMsg)
	return respBytes, err
}

