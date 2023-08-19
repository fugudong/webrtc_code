package db

import (
	"time"
)
// 数据库名称 rtc_room_server
// 第一期数据库操作目的：（1）跟踪记录；（2）搜集用户信息，所以我们先只在join, leave和connect的时候做处理

// 通话会话表 rtc_talk_sessions
// join时创建，connec成功时更新起始时间，leave时更新结束时间
type TalkSession struct {
	SessionId         	string    //所属会话ID	主键
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

// 统计客户端的通话质量，包含分辨率，帧率，视频码率
// 每个客户端只统计发送端
type TalkQualityStatsInfo struct {
	RemoteUserId 		string 	// 对端client的用户id
	SubSessionId 		string	// 子会话ID，因为一个client join后可能同时和多个client进行对话，1:1关系认为是一个子会话
	Time 		 	time.Time
	AudSendLostRate 	float32 // 音频发送丢包率
	AudSendBitRate 		int	// 音频发送比特率
	AudRecvLostRate 	float32 // 音频接收丢包率
	AudRecvBitRate 		int	// 音频接收比特率

	VidSendLostRate 	float32 // 视频发送丢包率
	VidSendBitRate 		int	// 视频发送比特率
	SendWidth 		int	// 发送分辨率
	SendHeight	 	int	// 发送分辨率
	//VidSendFrameRateInput	int	// 输入帧率
	SendFrameRateSent	int	// 实际发送的帧率

	VidRecvLostRate 	float32 // 视频接收丢包率
	VidRecvBitRate 		int	// 视频接收比特率
	//VidRecvWidth 		int	// 发送分辨率
	//VidRecvHeight	 	int	// 发送分辨率
	//VidRecvFrameRateOutput	int	// 对方发送的帧率
	RecvFrameRecv	int	//  实际收到的帧率
}

// 参考 GO语言使用数据库——MySQL https://blog.csdn.net/TDCQZD/article/details/82667785
// 添加通话会话记录
func AddTalkSession(session *TalkSession) error {
	stmtIns, err := dbConn.Prepare(`INSERT INTO rtc_talk_sessions
		(session_id, app_id, user_id, user_name, talk_type, 
		user_ip, create_time, join_time, leave_time, os_name, 
		browser, sdk_info) 
		VALUES(?, ?, ?, ?, ?,  ?, ?, ?, ?, ?,  ?, ?)`)

	if err != nil {
		return err
	}
	defer stmtIns.Close()
	_, err = stmtIns.Exec(session.SessionId, session.AppID, session.UserId,session.UserName, session.TalkType,
		session.UserIp, session.CreateTime, session.JoinTime, session.LeaveTime, session.OsName,
		session.Browser, session.SdkInfo)

	return  err
}
//
// 更新通话结束时间和总时间，在做数据分析时如果发现 FinishTime < BeginTime则说明属于通话异常
func UpdateTalkSessionFinishTimeAndDuration(session *TalkSession) error {
	// 先读取begin time

	sqlStr := `UPDATE rtc_talk_sessions SET time_out = ?, leave_time = ?, duration = ?, sub_session_ids = ? where session_id = ` + "\"" + session.SessionId  + "\""
	//fmt.Println("UpdateTalkSessionFinishTimeAndDuration: " + sqlStr)
	stmtIns, err := dbConn.Prepare(sqlStr)
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	_, err = stmtIns.Exec(session.Timeout, session.LeaveTime, session.Duration, session.SubSessionIds)
	if err != nil {
		return  err
	}

	return  nil
}

// 添加
func AddTalkSubSession(session *TalkSubSession) error {
	stmtIns, err := dbConn.Prepare(`INSERT INTO rtc_talk_subsessions
		(sub_session_id, status, remote_user_id, remote_user_name, remote_user_ip, 
		ice_ip, connect_type, begin_time, connect_time, finish_time, 
		duration) 
		VALUES(?, ?, ?, ?, ?,  ?, ?, ?, ?, ?, ?)`)

	if err != nil {
		return err
	}
	defer stmtIns.Close()
	_, err = stmtIns.Exec(session.SubSessionId, session.Status, session.RemoteUserId, session.RemoteUserName, session.RemoteUserIp,
		session.IceIp, session.ConnectType, session.BeginTime, session.ConnectTime, session.FinishTime,
		session.Duration)

	return  err
}

// 添加质量统计信息
func AddTalkQualityStatsInfo(stats *TalkQualityStatsInfo)  error {
	// 自动生成uid？
	//id := util.GenerateUserInfoId()
	sqlStr := `INSERT INTO rtc_talk_quality_stats_infos
		(remote_user_id, sub_session_id, create_time, aud_send_lost_rate, aud_send_bitrate, 
		aud_recv_lost_rate, aud_recv_bitrate, vid_send_lost_rate, vid_send_bitrate, send_width, 
		send_height, send_framerate_sent, vid_recv_lost_rate, vid_recv_bitrate, recv_framerate_recv)
		 VALUES(?, ?, ?, ?, ?,  ?, ?, ?, ?,?,  ?, ?, ?, ?, ?)`
	stmtIns, err := dbConn.Prepare(sqlStr)
	if err != nil {
		return   err
	}
	_, err = stmtIns.Exec(stats.RemoteUserId, stats.SubSessionId, stats.Time, stats.AudSendLostRate,
		stats.AudSendBitRate, stats.AudRecvLostRate, stats.AudRecvBitRate, stats.VidSendLostRate, stats.VidSendBitRate,
		stats.SendWidth, stats.SendHeight, stats.SendFrameRateSent, stats.VidRecvLostRate, stats.VidRecvBitRate,
		stats.RecvFrameRecv)
	if err != nil {
		return  err
	}
	defer stmtIns.Close()
	return  nil
}