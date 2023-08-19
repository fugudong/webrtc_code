package db

const TALK_OK 					= 0	// 正常通话
const TALK_ICE_FULL_LOADING 	= 1	// join时遇到ICE server full loading

// 连接类型
const DB_CONNECT_TYPE_INIT = "Unknown"
const DB_CONNECT_TYPE_TURN = "TURN"
const DB_CONNECT_TYPE_STUN = "STUN"
const DB_COUNTRY_INIT = "Unknown"

//0:正常通话和正常退出; 1:正常通话超时退出;2:异常通话正常退出;3异常通话异常退出
const TALK_SESSION_STATUS_INIT = -1
const TALK_SESSION_STATUS_NORMAL = 0
const TALK_SESSION_STATUS_TIMEOUT_LEAVE = 1
const TALK_SESSION_STATUS_FAILD = 2
const TALK_SESSION_STATUS_FAILD_AND_TIMEOUT_LEAVE = 3