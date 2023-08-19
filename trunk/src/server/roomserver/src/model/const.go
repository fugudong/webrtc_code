package model

type ROOM_ int32

// 与交互命令相关的常量
// 连接类型

// 用户类型 userType
const USER_TYPE_NORMAL = 0		// 普通用户
const USER_TYPE_OBSERVER = 1	// 观察者模式
// 通话类型talkType
const TALK_TYPE_AUDIO_ONLY = 0      //无视频有音频,
const TALK_TYPE_AUDIO_VIDEO = 1     //有视频有音频,
const TALK_TYPE_VIDEO_ONLY = 2      //有视频无音频,
const TALK_TYPE_NO_AUDIO_VIDEO = 3   //无视频无音频


const (
	ROOM_ERROR_SUCCESS  	= 0		// 为0
	ROOM_ERROR_FULL  	= 1		// 房间已经满了
	ROOM_ERROR_NOT_FIND_UID = 2		// 不能找到指定的用户ID
	ROOM_ERROR_PARSE_FAILED = 3		// 解析json错误
	ROOM_ERROR_WEBSOCKET_BROKEN 	= 4		// websocket出现异常
	ROOM_ERROR_NOT_FIND_RID		= 5	 	// 没有找到指定的房间
	ROOM_ERROR_NOT_FIND_REMOTEID 	= 6		// 远程ID
	ROOM_ERROR_WEBSOCKET_FAILED 	= 7		// websocket出错
	ROOM_ERROR_ICE_FULL_LOADING 	= 8		// ice server负载已满
	ROOM_ERROR_ICE_INVALID 		= 9		// ice server失效
	ROOM_ERROR_NO_MICROPHONE_DEV 	= 10
)
// 和ROOM_ERROR_XXX 常量对应
var ROOM_ERROR_MSG = []string{
	"successful",	// 0
	"room is full, it up to max number.",	// 1
	"can't find the designated uid, it may leave the room halfway", // 2
	"server parse the message faild, please check the format",	// 3
	"websocket may be broken",	// 4
	"can't find the designated room id",	// 5
	"can't find the remote user id",	//6
	"can't connect to room server ",	// 7
	"ice server bandwidth is full loading", //8
	"ice server invalid, please report to the 0voice technology", // 9
	"no microphone device, you can't use talk feature"}