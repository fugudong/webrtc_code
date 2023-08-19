package db

import (
	"../../src/util"
	"database/sql"
	"fmt"
	log "github.com/cihub/seelog"
	"strings"
)




// 指定日期的通话数量
func TalkSessionGetCountByDate(date string) (int,  error)  {
	query := fmt.Sprintf("select session_id from rtc_talk_sessions where  DATE_FORMAT(create_time,'%%Y%%m%%d') = '%s'", date)
	log.Infof("query:%s", query)
	//var session_id string
	stmtOut, err := dbConn.Prepare(query)
	if err != nil {
		log.Info("Prepare err:%s\n",err.Error());
		return 0, err
	}
	defer stmtOut.Close()
	rows, err := stmtOut.Query() // Query
	if err != nil {
		log.Info("Query err:%s\n",err.Error());
		return 0, err
	}
	count := 0
	for rows.Next() {		// 通过next的方式进行计数
		count += 1
	}

	return count, err
}

func printTalkReport(report *TalkSessionReprt)  {
	log.Infof("uname:%s, type:%s, bt:%s, d:%s, s:%s, c:%s, uip:%s, iceip:%s, osb:%s",
		report.UserName,
		report.TalkType,
		report.BeginTime,
		report.Duration,
		report.Status,
		report.ConnectType,
		report.UserCountryIp,
		report.IceCountryIp,
		report.Browser,
		)
}

func printTalkSession(ses *TalkSession)  {
	log.Infof("id:%d s:%s, sub:%s, appid:%s, u:%s un:%s, tt:%d, uip:%s, t:%t, %s %s %s, d:%d, %s %s %s",
		ses.Id, ses.SessionID, ses.SubSessionIds, ses.AppID, ses.UserId,
		ses.UserName, ses.TalkType, ses.UserIp, ses.Timeout, ses.CreateTime.Format("2006-01-02 15:04:05"),
		ses.JoinTime.Format("2006-01-02 15:04:05"), ses.LeaveTime.Format("2006-01-02 15:04:05"), ses.Duration, ses.OsName, ses.Browser,
		ses.SdkInfo)
}

func printTalkSubSession(subSes *TalkSubSession)  {
	log.Infof("sid:%s, c:%s, s:%d, rid:%s, rname:%s, rip:%s, iceip:%s, bt:%s, ct:%s, ft:%s, d:%d",
		subSes.SubSessionId,
		subSes.ConnectType,
		subSes.Status,
		subSes.RemoteUserId,
		subSes.RemoteUserName,
		subSes.RemoteUserIp,
		subSes.IceIp,
		subSes.BeginTime.Format("2006-01-02 15:04:05"),
		subSes.ConnectTime.Format("2006-01-02 15:04:05"),
		subSes.FinishTime.Format("2006-01-02 15:04:05"),
		subSes.Duration,
	)
}

func TalkSessionsList(date string, start int, count int) ([] *TalkSessionReprt, error)  {

	var talkReports []*TalkSessionReprt
	query := fmt.Sprintf("select * from rtc_talk_sessions where  DATE_FORMAT(create_time,'%%Y%%m%%d') = '%s' limit %d,%d", date, start, count)
	log.Infof("query:%s", query)
	//var session_id string
	stmtOut, err := dbConn.Prepare(query)
	if err != nil {
		log.Errorf("Prepare err:%s\n",err.Error());
		return talkReports, err
	}
	defer stmtOut.Close()
	
	rows, err := stmtOut.Query() // Query
	if err != nil {
		log.Errorf("Query err:%s\n",err.Error());
		return talkReports, err
	}

	for rows.Next() {		// 通过next的方式进行计数
		var  talk TalkSession

		if err := rows.Scan(&talk.Id, &talk.SessionID, &talk.SubSessionIds, &talk.AppID, &talk.UserId,
			&talk.UserName, &talk.TalkType, &talk.UserIp, &talk.Timeout, &talk.CreateTime,
			&talk.JoinTime, &talk.LeaveTime, &talk.Duration, &talk.OsName, &talk.Browser,
			&talk.SdkInfo); err != nil {
			return talkReports, err
		}
		printTalkSession(&talk)
		subSessionIds := strings.Split(talk.SubSessionIds, ",")
		subNum := len(subSessionIds)
		log.Infof("SubSessionIds:%s, len = %d", talk.SubSessionIds, subNum)
		if subNum == 0 || talk.SubSessionIds == "" {
			var talkReport TalkSessionReprt
			talkReport.AppId = talk.AppID
			talkReport.UserName = talk.UserName
			talkReport.Status = STATUS_NO_TALK
			talkReport.TalkType = getTalkTypeString(talk.TalkType)
			talkReport.ConnectType = CONNECT_TYPE_UNKNOWN
			talkReport.UserCountryIp = "USER " + util.GeoIpTable.GeoIpGetCountry(talk.UserIp) + ":" + talk.UserIp
			talkReport.BeginTime = talk.JoinTime.Format("01-02 15:04:05") //"2006-01-02 15:04:05"
			h,m,s := util.ResolveTime(talk.Duration)
			duration := fmt.Sprintf("%02d:%02d:%02d", h, m, s)
			talkReport.Duration =  duration
			talkReport.Browser = talk.OsName + ":" + talk.Browser
			printTalkReport(&talkReport)
			talkReports = append(talkReports, &talkReport)			// 加入数组
		} else {
			var stmtOutSubsession *sql.Stmt
			for i:=0; i < subNum; i++ {
				var talkReport TalkSessionReprt
				// 读取子会话
				subSessionId := subSessionIds[i]
				query := fmt.Sprintf("select * from rtc_talk_subsessions where  sub_session_id = '%s'", subSessionId)
				stmtOutSubsession, err := dbConn.Prepare(query)
				if err != nil {
					log.Errorf("Prepare rtc_talk_subsessions faild:%s\n",err.Error());
					continue
				}
				subRows, err := stmtOutSubsession.Query() // Query
				if err != nil {
					log.Errorf("Query rtc_talk_subsessions failed:%s\n",err.Error());
					continue
				}

				var subTalk TalkSubSession
				for subRows.Next() {
					if err := subRows.Scan(&subTalk.Id, &subTalk.SubSessionId, &subTalk.Status, &subTalk.RemoteUserId, &subTalk.RemoteUserName,
						&subTalk.RemoteUserIp, &subTalk.IceIp, &subTalk.ConnectType, &subTalk.BeginTime, &subTalk.ConnectTime,
						&subTalk.FinishTime, &subTalk.Duration, &subTalk.Cost); err != nil {
							log.Errorf("Scan rtc_talk_subsessions failed:%s", err.Error())
						continue
					}
				}
				printTalkSubSession(&subTalk)
				talkReport.AppId = talk.AppID
				talkReport.UserName = talk.UserName
				talkReport.Status = getStatusString(subTalk.Status)
				talkReport.ConnectType = getConnectTypeString(subTalk.ConnectType)
				//log.Infof("talk.UserIp:%s", talk.UserIp)
				talkReport.UserCountryIp = "USER "+ util.GeoIpTable.GeoIpGetCountry(talk.UserIp) + ":" + talk.UserIp
				if subTalk.IceIp != "" {
					talkReport.IceCountryIp =  "ICE  "+ util.GeoIpTable.GeoIpGetCountry(subTalk.IceIp) + ":" + subTalk.IceIp
				}
				h,m,s := util.ResolveTime(subTalk.Duration)
				duration := fmt.Sprintf("%02d:%02d:%02d", h, m, s)
				talkReport.Duration = duration
				// 读取质量统计信息
				queryStats := fmt.Sprintf("select aud_recv_lost_rate, aud_recv_bitrate, vid_recv_lost_rate, vid_recv_bitrate from rtc_talk_quality_stats_infos where  sub_session_id = '%s'", subSessionId)
				stmtOutStats, err := dbConn.Prepare(queryStats)
				if err != nil {
					log.Errorf("Prepare err:%s\n",err.Error());
					continue
				}
				subRows, err = stmtOutStats.Query() // Query
				if err != nil {
					log.Errorf("Query err:%s\n",err.Error());
					continue
				}
				var audAvgBitRate, audMinBitRate, audMaxBitRate, audSumBitRate int
				var vidAvgBitRate, vidMinBitRate, vidMaxBitRate, vidSumBitRate int
				var audAvgLostRate, audMinLostRate, audMaxLostRate, audSumLostRate float32
				var vidAvgLostRate, vidMinLostRate, vidMaxLostRate, vidSumLostRate float32
				audMinBitRate = 1000
				vidMinBitRate = 2000
				audMinLostRate = 1.0
				vidMinLostRate = 1.0
				count := 0
				for subRows.Next() {
					var talkStat TalkQualityStatsInfo
					if err := subRows.Scan(&talkStat.AudRecvLostRate, &talkStat.AudRecvBitRate,  &talkStat.VidRecvLostRate,  &talkStat.VidRecvBitRate); err != nil {
						log.Errorf("Scan rtc_talk_quality_stats_infos failed:%s", err.Error())
						continue
					}
					count += 1
					// 码率 音频
					if talkStat.AudRecvBitRate < audMinBitRate { audMinBitRate = talkStat.AudRecvBitRate}
					if talkStat.AudRecvBitRate > audMaxBitRate { audMaxBitRate = talkStat.AudRecvBitRate}
					audSumBitRate += talkStat.AudRecvBitRate
					// 码率 视频
					if talkStat.VidRecvBitRate < vidMinBitRate { vidMinBitRate = talkStat.VidRecvBitRate}
					if talkStat.VidRecvBitRate > vidMaxBitRate { vidMaxBitRate = talkStat.VidRecvBitRate}
					vidSumBitRate += talkStat.VidRecvBitRate

					// 丢包率 音频
					if talkStat.AudRecvLostRate < audMinLostRate { audMinLostRate = talkStat.AudRecvLostRate}
					if talkStat.AudRecvLostRate > audMaxLostRate { audMaxLostRate = talkStat.AudRecvLostRate}
					audSumLostRate += talkStat.AudRecvLostRate
					// 丢包率 视频
					if talkStat.VidRecvLostRate < vidMinLostRate { vidMinLostRate = talkStat.VidRecvLostRate}
					if talkStat.VidRecvLostRate > vidMaxLostRate { vidMaxLostRate = talkStat.VidRecvLostRate}
					vidSumLostRate += talkStat.VidRecvLostRate
				}
				if count == 0 {
					talkReport.AudAvgMaxMinBitRate = STATS_NO_DATA
					talkReport.AudAvgMaxMinLostRate = STATS_NO_DATA
					talkReport.VidAvgMaxMinBitRate = STATS_NO_DATA
					talkReport.VidAvgMaxMinLostRate = STATS_NO_DATA
				} else {
					audAvgBitRate = audSumBitRate / count
					audAvgLostRate = audSumLostRate / float32(count)
					vidAvgBitRate = vidSumBitRate / count
					vidAvgLostRate = vidSumLostRate / float32(count)
					if talk.TalkType == TALK_TYPE_AUDIO_ONLY {
						talkReport.AudAvgMaxMinBitRate = fmt.Sprintf("A:%d,%d,%d",audAvgBitRate, audMinBitRate, audMaxBitRate)
						talkReport.AudAvgMaxMinLostRate= fmt.Sprintf("A:%0.2f,%0.2f,%0.2f",audAvgLostRate, audMinLostRate, audMaxLostRate)
						talkReport.VidAvgMaxMinBitRate = STATS_NO_DATA
						talkReport.VidAvgMaxMinLostRate = STATS_NO_DATA
					} else if talk.TalkType == TALK_TYPE_AUDIO_VIDEO{
						talkReport.AudAvgMaxMinBitRate = fmt.Sprintf("A:%d,%d,%d", audAvgBitRate, audMinBitRate, audMaxBitRate)
						talkReport.VidAvgMaxMinBitRate = fmt.Sprintf("V:%d,%d:%d", vidAvgBitRate, vidMinBitRate, vidMaxBitRate)
						talkReport.AudAvgMaxMinLostRate= fmt.Sprintf("A:%0.2f,%0.2f,%0.2f",
							audAvgLostRate, audMinLostRate, audMaxLostRate,)
						talkReport.VidAvgMaxMinLostRate= fmt.Sprintf("V:%0.2f,%0.2f,%0.2f",
							vidAvgLostRate, vidMinLostRate, vidMaxLostRate)
					}
				}
				talkReport.TalkType = getTalkTypeString(talk.TalkType)
				talkReport.BeginTime = subTalk.BeginTime.Format("01-02 15:04:05")
				talkReport.Browser = talk.OsName + ":" + talk.Browser
				printTalkReport(&talkReport)
				talkReports = append(talkReports, &talkReport)			// 加入数组
			}
			if stmtOutSubsession != nil {
				stmtOutSubsession.Close()
			}
		}
	}

	return talkReports, err
}