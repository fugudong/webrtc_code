package logic

import (
	"roomserver/src/db"
	"roomserver/src/model"
	Xlog "github.com/cihub/seelog"
	"strconv"
	"time"
)



func joinOperateDb(join *model.JoinCmd, c *Client) error {
	// 写数据库
	var session db.TalkSession
	session.SessionId = c.SessionId;
	session.AppID = join.AppID
	session.UserId = join.UID
	session.UserName = join.Uname
	session.TalkType = join.TalkType
	session.UserIp = join.IP


	session.CreateTime = time.Now()
	session.JoinTime = c.JoinTime	// 加入和离开则使用当地时间
	session.LeaveTime = c.JoinTime	// 时间戳
	session.Duration = 0
	session.OsName = join.OsName
	session.Browser = join.Browser
	session.SdkInfo = join.SdkInfo


	//Xlog.Infof("join session：join:%s, leave:%s duration = %d",
	//	session.JoinTime.Format("2006-01-02 15:04:05"),
	//	session.LeaveTime.Format("2006-01-02 15:04:05"),
	//	session.Duration)
	// 添加通话记录
	err := db.AddTalkSession(&session)
	if err != nil {
		Xlog.Error("AddTalkSession failed: " + err.Error())
		// 需发邮件通知写数据库失败
		return err
	}

	return nil
}


func LeaveOperateDb(leave *model.LeaveCmd, c *Client) error {
	// 写数据库
	var session db.TalkSession

	session.SessionId = c.SessionId
	session.Timeout =  c.IsTimeoutLeave

	curTime, err := strconv.ParseInt(leave.Time, 10, 64)
	if err != nil {
		Xlog.Errorf("Parse leave.time:%s failed: %s",leave.Time)
	}
	tm := time.Unix(curTime, 0)
	session.LeaveTime = tm		// 时间戳					离开时间
	session.JoinTime = c.JoinTime
	// 起始时间使用在join或者offer或者answer时保存的BeginTime
	session.Duration = int64(session.LeaveTime.Sub(session.JoinTime).Seconds())	// 持续时长

	session.SubSessionIds = "";
	length := len(c.SubSessionIds)				// 子会话
	for i := 0; i < length; i++ {
		session.SubSessionIds += c.SubSessionIds[i]
		if i < length - 1 {
			session.SubSessionIds += ","
		}
	}
	Xlog.Infof("leave session：join:%s, leave:%s duration = %d,SubSessionIds:%s ",
		session.JoinTime.Format("2006-01-02 15:04:05"),
		session.LeaveTime.Format("2006-01-02 15:04:05"),
		session.Duration, session.SubSessionIds)
	// 添加通话记录
	err = db.UpdateTalkSessionFinishTimeAndDuration(&session)
	if err != nil {
		Xlog.Error("UpdateTalkSessionFinishTimeAndDuration failed: " + err.Error())
		// 需发邮件通知写数据库失败
	}
	return err
}

func LeaveSubOperateDb(leave *model.LeaveCmd, c *Client, rc *RemoteClient) error {
	// 写数据库
	var session db.TalkSubSession
	var status int

	session.SubSessionId = rc.SubSessionId
	if rc.IsPeerConnected && !c.IsTimeoutLeave {
		status = db.TALK_SESSION_STATUS_NORMAL	// 通话正常并正常退出
	} else if rc.IsPeerConnected && c.IsTimeoutLeave{
		status = db.TALK_SESSION_STATUS_TIMEOUT_LEAVE// 通话正常超时退出
	} else if !rc.IsPeerConnected && !c.IsTimeoutLeave{
		status = db.TALK_SESSION_STATUS_FAILD// 通话异常正常退出
	} else {
		status = db.TALK_SESSION_STATUS_FAILD_AND_TIMEOUT_LEAVE// 通话异常超时退出
	}
	session.Status = status
	session.RemoteUserId = rc.UID
	session.RemoteUserName = rc.Uname
	session.RemoteUserIp = rc.IceIP
	session.IceIp = rc.IceIP
	session.ConnectType = rc.ConnectType
	session.BeginTime = rc.BeginTime
	session.ConnectTime = rc.ConnectTime
	session.FinishTime = rc.FinishTime

	session.Duration = int64(rc.FinishTime.Sub(rc.ConnectTime).Seconds())
	if session.Duration < 0 {
		Xlog.Error("rc.FinishTime.Sub(rc.ConnectTime) < 0 ")
		session.Duration = 0
	}

	// 添加子通话记录
	err := db.AddTalkSubSession(&session)
	if err != nil {
		Xlog.Error("AddTalkSubSession failed: " + err.Error())
		// 需发邮件通知写数据库失败
	}
	return err
}

func reportStatsOperateDb(reportStats *model.ReportStatsCmd, subSessionId string) error {
	var statsInfo db.TalkQualityStatsInfo
	statsInfo.SubSessionId = subSessionId
	statsInfo.RemoteUserId = reportStats.RemoteUID
	curTime, err := strconv.ParseInt(reportStats.Time, 10, 64)
	if err != nil {
		Xlog.Errorf("reportStats leave.time:%s failed: %s",reportStats.Time)
	}
	tm := time.Unix(curTime, 0)
	statsInfo.Time = tm
	// 音频相关
	packetsLostRate, err := strconv.ParseFloat(reportStats.Audio.Send.PacketsLostRate, 32)
	statsInfo.AudSendLostRate = float32(packetsLostRate)
	statsInfo.AudSendBitRate = reportStats.Audio.Send.BitRate
	packetsLostRate, err = strconv.ParseFloat(reportStats.Audio.Recv.PacketsLostRate, 32)
	statsInfo.AudRecvLostRate = float32(packetsLostRate)
	statsInfo.AudRecvBitRate = reportStats.Audio.Recv.BitRate
	// 视频相关
	packetsLostRate, err = strconv.ParseFloat(reportStats.Video.Send.PacketsLostRate, 32)
	statsInfo.VidSendLostRate = float32(packetsLostRate)
	statsInfo.VidSendBitRate = reportStats.Video.Send.BitRate
	statsInfo.SendWidth = reportStats.Video.Send.Width
	statsInfo.SendHeight = reportStats.Video.Send.Height
	statsInfo.SendFrameRateSent = reportStats.Video.Send.FrameRateSent

	packetsLostRate, err = strconv.ParseFloat(reportStats.Video.Recv.PacketsLostRate, 32)
	statsInfo.VidRecvLostRate = float32(packetsLostRate)
	statsInfo.VidRecvBitRate = reportStats.Video.Recv.BitRate
	statsInfo.RecvFrameRecv = reportStats.Video.Recv.FrameRateRecv

	// 添加通话记录
	err = db.AddTalkQualityStatsInfo(&statsInfo)
	if err != nil {
		Xlog.Error("AddTalkQualityStatsInfo failed: " + err.Error())
		// 需发邮件通知写数据库失败
	}

	return err
}