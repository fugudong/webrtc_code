package db

import (
	"time"
)

type TalkSession struct {
	Id         			int64
	SessionID         	string    //所属会话ID	主键
	SubSessionIds 		string    // 每个1:1通话关系生成一个子会话
	AppID             	string    //应用ID
	UserId            	string    //用户ID
	UserName			string	  // 用户名
	UserIp				string	  // 用户IP
	TalkType          	int       // 通话类型
	Timeout				bool
	CreateTime         	time.Time // 使用系统时间统一创建
	JoinTime         	time.Time // 加入时间
	LeaveTime        	time.Time // 离开时间
	Duration          	int64     // 通话时长，单位秒（属于冗余字段，但方便后台直接查看通话时间）
	OsName             	string
	Browser				string
	SdkInfo				string
}

type TalkSubSession struct {
	Id         			int64
	SubSessionId 		string    // 每个1:1通话关系生成一个子会话
	ConnectType       	string    //  TRUN和STUN, 初始为 Unknown
	Status            	int       // 0:正常通话和正常退出; 1:正常通话超时退出;2:异常通话正常退出;3异常通话异常退出
	RemoteUserId		string	  // 对端ID
	RemoteUserName		string	  // 对端名字
	RemoteUserIp		string	  //对端用户ip
	IceIp            	string    // 默认为""
	BeginTime        	time.Time // 通话起始时间
	ConnectTime			time.Time
	FinishTime        	time.Time // 通话结束时间
	Duration          	int64     // 通话时长，单位秒（属于冗余字段，但方便后台直接查看通话时间）
	Cost              	int32     // 计费-单位分
}

type TalkQualityStatsInfo struct {
	AudRecvLostRate 	float32 // 音频接收丢包率
	AudRecvBitRate 		int	// 音频接收比特率
	VidRecvLostRate 	float32 // 视频接收丢包率
	VidRecvBitRate 		int	// 视频接收比特率
}

//应用ID 用户ID 通话类型 通话结束情况 连接类型
// 用户国家和IP ICE服务器国家和IP
//全部用文字说明，避免使用数字无法理解
type TalkSessionReprt struct {
	AppId 			string
	UserName 			string
	TalkType 		string		// 音频 音视频 视频 无音视频
	Status 			string		// 正常结束，超时退出，未连接
	ConnectType 	string		// 连接类型：P2P，中继，未知
	UserCountryIp 	string		// 用户国家IP
	IceCountryIp 	string		// ICE国家ip
	BeginTime 		string		// 通话起始时间
	Duration  		string		// 通话时长
	Browser			string
	// 以下为接收的数据统计，如果通话时间太短没有收到统计消息则设置为未知
	AudAvgMaxMinBitRate		string	//平均码率 最大码率, 最小码率 当没有统计结果时显示为未知，合并一起avg:127,max:200,min:100,
	VidAvgMaxMinBitRate		string	//平均码率 最大码率, 最小码率 当没有统计结果时显示为未知，合并一起avg:127,max:200,min:100,
	AudAvgMaxMinLostRate  	string	// 平均丢包率,最大丢包率，最小丢包率,合并一起 avg:0.1,max:0.2,min:0
	VidAvgMaxMinLostRate  	string	// 平均丢包率,最大丢包率，最小丢包率,合并一起 avg:0.1,max:0.2,min:0
}
// 通话类型
const TALK_TYPE_STR_AUDIO 			= "A"
const TALK_TYPE_STR_AUDIO_VIDEO 	= "AV"
const TALK_TYPE_STR_VIDEO 			= "V"
const TALK_TYPE_STR_NO_AV 			= "NO"
const TALK_TYPE_STR_UNKNOWN 	    = "NG"
// 是否正常通话
const STATUS_NORMAL 			= "正常"
const STATUS_TIMEOUT_LEAVE 		= "超时退出"
const STATUS_NO_CONNECT 		= "未连接"
const STATUS_NO_TALK			= "未发起通话"
const STATUS_UNKNOWN			= "未知"

const TALK_UNKOWN 				= "未知"

const CONNECT_TYPE_RELAY	= "中继"
const CONNECT_TYPE_P2P		= "P2P"
const CONNECT_TYPE_UNKNOWN = "未知"

const TALK_SESSION_STATUS_INIT = -1
const TALK_SESSION_STATUS_NORMAL = 0
const TALK_SESSION_STATUS_TIMEOUT_LEAVE = 1
const TALK_SESSION_STATUS_FAILD = 2
const TALK_SESSION_STATUS_FAILD_AND_TIMEOUT_LEAVE = 3

// 通话类型talkType
const TALK_TYPE_AUDIO_ONLY = 0      //无视频有音频,
const TALK_TYPE_AUDIO_VIDEO = 1     //有视频有音频,
const TALK_TYPE_VIDEO_ONLY = 2      //有视频无音频,
const TALK_TYPE_NO_AUDIO_VIDEO = 3   //无视频无音频

const STATS_NO_DATA = "--"

func getStatusString(status int) string {
	if status == TALK_SESSION_STATUS_NORMAL {
		return "通话正常"
	} else if status == TALK_SESSION_STATUS_TIMEOUT_LEAVE {
		return "正常通话超时退出"
	} else if status == TALK_SESSION_STATUS_FAILD {
		return "通话异常"
	} else if status == TALK_SESSION_STATUS_FAILD_AND_TIMEOUT_LEAVE {
		return "通话异常超时退出"
	} else {
		return "未知"
	}
}

func getConnectTypeString(connectType string) string {
	if connectType == "TURN" {
		return CONNECT_TYPE_RELAY
	} else if connectType == "STUN" {
		return CONNECT_TYPE_P2P
	} else {
		return CONNECT_TYPE_UNKNOWN
	}
}

func getTalkTypeString(talkType int) string {
	if talkType == TALK_TYPE_AUDIO_ONLY {
		return TALK_TYPE_STR_AUDIO
	} else if talkType == TALK_TYPE_AUDIO_VIDEO {
		return TALK_TYPE_STR_AUDIO_VIDEO
	} else if talkType == TALK_TYPE_VIDEO_ONLY{
		return TALK_TYPE_STR_VIDEO
	} else if talkType == TALK_TYPE_NO_AUDIO_VIDEO{
		return TALK_TYPE_STR_NO_AV
	} else  {
		return  TALK_TYPE_STR_UNKNOWN
	}
}



