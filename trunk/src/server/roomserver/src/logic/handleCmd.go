package logic

import (
	"roomserver/src/model"
	"roomserver/src/model/iceBalance"
	"strconv"
	"io"
	"time"
	"golang.org/x/net/websocket"
	Xlog "github.com/cihub/seelog"
)




// Collider 碰撞机
type Collider struct {
	roomTable *RoomTable		// 房间table 字典
	WsCount int
}

func NewCollider() *Collider {
	return &Collider{
		roomTable: NewRoomTable(time.Second* REGISTER_TIMEOUT_SEC),
		WsCount: 0,
	}
}

// 当解析client's message出错时回应客户端相应的提示消息
func handleClientMessageFormatFail(cmd string, rwc io.ReadWriteCloser) int {
	respMsg, err2  := model.GeneralRespMsgMarshal(cmd, model.ROOM_ERROR_PARSE_FAILED,
		model.ROOM_ERROR_MSG[model.ROOM_ERROR_PARSE_FAILED], "", "")
	Send(rwc, respMsg)	// 错误
	if err2 != nil {
		Xlog.Warn("GeneralRespMsgMarshal ROOM_ERROR_PARSE_FAILED failed")
	}

	return model.ROOM_ERROR_PARSE_FAILED
}

// 当房间不存在时则回应客户端相应的提示信息
func handleFindRoomFail(cmd string, rwc io.ReadWriteCloser) int {
	respMsg, err2  := model.GeneralRespMsgMarshal(cmd, model.ROOM_ERROR_NOT_FIND_RID,
		model.ROOM_ERROR_MSG[model.ROOM_ERROR_NOT_FIND_RID], "", "")
	Send(rwc, respMsg)	// 错误
	if err2 != nil {
		Xlog.Warn("GeneralRespMsgMarshal ROOM_ERROR_NOT_FIND_RID failed")
	}

	return model.ROOM_ERROR_NOT_FIND_RID
}

// 处理加入房间事件
func (c *Collider)HandleJoin(message []byte, rwc *websocket.Conn) int {
	Xlog.Debug("HandleJoin into")
	//1. 验证message是否可以正常解析json
	var join model.JoinCmd
	err := model.JoinUnMarshal(message, &join)
	if err != nil {
		Xlog.Errorf("Cmd:%s -> JoinUnMarshal failed: %s", model.CMD_JOIN, err.Error())
		return handleClientMessageFormatFail(model.CMD_JOIN, rwc)
	}

	join.IP = rwc.Request().Header.Get("Remote_addr")		// 获取IP
	//join.IP = "43.250.201.20"
	// 2. 先清除原有的信息
	leave := model.JoinGenerateLeave(&join)
	r := c.roomTable.GetRoomLocked(leave.RoomID)	// 查找对应的房间是否存在
	if r != nil {	// 只有找到房间的情况下才去查找client是否存在
		//ret := r.JoinLeave(&leave, rwc)
		//if ret != model.ROOM_ERROR_SUCCESS {
			// 打印错误原因，不影响下一步操作
		//}
		Xlog.Warnf("%s join room %s again", join.UID, join.RoomID)
	}

	//3.查找房间是否存在，如果不存在则创建一个
	r = c.roomTable.room(join.RoomID,join.RoomName)
	//4.处理加入
	ret := r.Join(&join, rwc, c)

	return ret
}

// 处理离开房间事件
func (c *Collider)HandleLeave(message []byte, rwc io.ReadWriteCloser) (int) {
	Xlog.Debug("HandleLeave into")
	//1. 验证json格式
	var leave model.LeaveCmd
	err := model.LeaveUnMarshal(message, &leave)
	if err != nil {
		Xlog.Errorf("Cmd:%s -> LeaveUnMarshal failed", model.CMD_LEAVE)
		return handleClientMessageFormatFail(model.CMD_LEAVE, rwc)
	}
	r := c.roomTable.GetRoomLocked(leave.RoomID)
	if r != nil {
		// 离开房间
		ret := r.Leave(&leave, rwc, false)
		if ret == model.ROOM_ERROR_SUCCESS {
			// 如果房间没人存在则删除房间
			c.roomTable.removeRoom(leave.RoomID)
		}
		return ret
	} else {
		Xlog.Errorf("Cmd:%s -> Client:%s can't find the room:%s", leave.Cmd, leave.UID, leave.RoomID)
		return handleFindRoomFail(model.CMD_LEAVE, rwc)
	}
}

// 处理offer信令交互
func (c *Collider)HandleOffer(message []byte, rwc io.ReadWriteCloser) (int) {
	Xlog.Debug("HandleOffer into")

	// 1. 验证客户是否已经在房间和token是否已经过时
	var offer model.OfferCmd
	err := model.OfferUnMarshal(message, &offer)
	if err != nil {
		Xlog.Errorf("Cmd:%s -> OfferUnMarshal failed: %s", model.CMD_OFFER, err.Error())
		return handleClientMessageFormatFail(model.CMD_OFFER, rwc)
	}
	// 查找房间是否存在
	r := c.roomTable.GetRoomLocked(offer.RoomID)
	if r != nil {
		return r.Offer(&offer, rwc)
	} else {
		Xlog.Errorf("Cmd:%s -> Client:%s can't find the room:%s", offer.Cmd, offer.UID, offer.RoomID)
		return handleFindRoomFail(model.CMD_OFFER, rwc)
	}
}

// 处理answer信令交互
func (c *Collider)HandleAnswer(message []byte, rwc io.ReadWriteCloser) (int) {
	Xlog.Debug("HandleAnswer into")

	// 1. 验证客户是否已经在房间和token是否已经过时
	var answer model.AnswerCmd
	err := model.AnswerUnMarshal(message, &answer)
	if err != nil {
		Xlog.Errorf("Cmd:%s -> AnswerUnMarshal failed: %s", model.CMD_ANSWER, err.Error())
		return handleClientMessageFormatFail(model.CMD_ANSWER, rwc)
	}
	// 查找房间是否存在
	r := c.roomTable.GetRoomLocked(answer.RoomID)
	if r != nil {
		return  r.Answer(&answer, rwc)
	} else {
		Xlog.Errorf("Cmd:%s -> Client:%s can't find the room:%s", answer.Cmd, answer.UID, answer.RoomID)
		return handleFindRoomFail(model.CMD_ANSWER, rwc)
	}
}

// 处理candidate信令交互
func (c *Collider)HandleCandidate(message []byte, rwc io.ReadWriteCloser) (int) {
	Xlog.Debug("HandleCandidate into")

	// 1. 验证客户是否已经在房间和token是否已经过时
	var candidate model.CandidateCmd
	err := model.CandidateUnMarshal(message, &candidate)
	if err != nil {
		Xlog.Errorf("Cmd:%s -> CandidateUnMarshal failed: %s", model.CMD_CANDIDATE, err.Error())
		return handleClientMessageFormatFail(model.CMD_CANDIDATE, rwc)
	}
	// 查找房间是否存在
	r := c.roomTable.GetRoomLocked(candidate.RoomID)
	if r != nil {
		return r.Candidate(&candidate, rwc)
	} else {
		Xlog.Errorf("Cmd:%s -> Client:%s can't find the room:%s", candidate.Cmd, candidate.UID, candidate.RoomID)
		return handleFindRoomFail(model.CMD_CANDIDATE, rwc)
	}
}

// 报告xx已经和xx处于连接状态，怎么样才能认为是已经成功连接？检测到到正常发送则认为是连接成功？
func (c *Collider)HandlePeerConnected(message []byte, rwc io.ReadWriteCloser) (int) {
	Xlog.Debug("HandlePeerConnected into")

	// 1. 验证客户是否已经在房间和token是否已经过时
	var peerConnected model.PeerConnectedCmd
	err := model.PeerConnecteUnMarshal(message, &peerConnected)
	if err != nil {
		Xlog.Errorf("Cmd:%s -> PeerConnecteUnMarshal failed: %s", model.CMD_PEER_CONNECTED, err.Error())
		return handleClientMessageFormatFail(model.CMD_PEER_CONNECTED, rwc)
	}
	// 查找房间是否存在
	r := c.roomTable.GetRoomLocked(peerConnected.RoomID)
	if r != nil {
		return r.PeerConnected(&peerConnected, rwc)
	} else {
		Xlog.Errorf("Cmd:%s -> Client:%s can't find the room:%s",
			peerConnected.Cmd, peerConnected.UID, peerConnected.RoomID)
		return handleFindRoomFail(model.CMD_PEER_CONNECTED, rwc)
	}
}

// 心跳包
func (c *Collider)HandleKeepLive(message []byte, rwc io.ReadWriteCloser) (int) {
	//Xlog.Debug("HandleKeepLive into")
	var keppLive model.KeepLiveCmd
	err := model.KeepLiveUnMarshal(message, &keppLive)
	if err != nil {
		Xlog.Errorf("Cmd:%s -> KeepLiveUnMarshal failed: %s", model.CMD_KEEP_LIVE, err.Error())
		return handleClientMessageFormatFail(model.CMD_KEEP_LIVE, rwc)
	}
	// 查找房间是否存在
	r := c.roomTable.GetRoomLocked(keppLive.RoomID)
	if r != nil {
		return r.KeepLive(&keppLive, rwc)
	} else {
		Xlog.Errorf("Cmd:%s -> Client:%s can't find the room:%s", keppLive.Cmd, keppLive.UID, keppLive.RoomID)
		return handleFindRoomFail(model.CMD_KEEP_LIVE, rwc)
	}
}

func (c *Collider)HandleTurnTalkType(message []byte, rwc io.ReadWriteCloser) (int) {
	Xlog.Debug("HandleTurnTalkType into")
	// 1. 验证客户是否已经在房间和token是否已经过时
	var turnTalkType model.TurnTalkTypeCmd
	err := model.TurnTalkTypeUnMarshal(message, &turnTalkType)
	if err != nil {
		Xlog.Errorf("Cmd:%s -> TurnTalkTypeUnMarshal failed: %s", model.CMD_TURN_TALK_TYPE, err.Error())
		return handleClientMessageFormatFail(model.CMD_TURN_TALK_TYPE, rwc)
	}

	// 查找房间是否存在
	r := c.roomTable.GetRoomLocked(turnTalkType.RoomID)
	if r != nil {
		return  r.TurnTalkType(&turnTalkType, rwc)
	} else {
		Xlog.Errorf("Cmd:%s -> Client:%s can't find the room:%s", turnTalkType.Cmd, turnTalkType.UID, turnTalkType.RoomID)
		return handleFindRoomFail(model.CMD_TURN_TALK_TYPE, rwc)
	}
}

// 处理超时离开事件
func (c *Collider)HandleTimeoutLeave(client *Client , room *Room, rwc io.ReadWriteCloser) (int) {
	Xlog.Info("HandleTimeoutLeave into")
	//1. 验证json格式
	var leave model.LeaveCmd

	// 封装leave
	leave.Cmd = model.CMD_LEAVE
	leave.AppID = client.AppID
	leave.RoomID = room.Id // 房间号
	leave.UID = client.UID
	leave.Uname = client.Uname
	leave.UserType = client.UserType
	curTime := time.Now()			// 当前时间
	offsetTime := curTime.Sub(client.ServJoinTime)	// 获取经过的时间
	leaveTime := client.JoinTime.Add(offsetTime)	// 换算成结束时间
	leave.Time =  strconv.FormatInt(leaveTime.Unix(), 10)

	r := c.roomTable.GetRoomLocked(leave.RoomID)
	if r != nil {
		// 离开房间
		ret := r.Leave(&leave, rwc, true)
		if ret == model.ROOM_ERROR_SUCCESS {
			// 如果房间没人存在则删除房间
			c.roomTable.removeRoom(leave.RoomID)
		}
		return ret
	} else {
		Xlog.Errorf("Cmd:%s -> Client:%s can't find the room:%s", leave.Cmd, leave.UID, leave.RoomID)
		return handleFindRoomFail(model.CMD_LEAVE, rwc)
	}
}

func (c *Collider)HandleReportInfo(message []byte, rwc io.ReadWriteCloser) (int) {
	Xlog.Info("HandleReportInfo into")

	// 1. 验证客户是否已经在房间和token是否已经过时
	var reportInfoCmd model.ReportInfoCmd
	err := model.ReportInfoUnMarshal(message, &reportInfoCmd)
	if err != nil {
		Xlog.Errorf("Cmd:%s -> ReportInfoUnMarshal failed: %s", model.CMD_REPORT_INFO, err.Error())
		return handleClientMessageFormatFail(model.CMD_REPORT_INFO, rwc)
	}
	if reportInfoCmd.Result != model.ROOM_ERROR_SUCCESS {
		Xlog.Errorf("Room:%s, uid:%s, uname:%s report failed: reselut:%d, desc:%s, data1:%s, data2:%s",
			reportInfoCmd.RoomID, reportInfoCmd.UID, reportInfoCmd.Uname,
			reportInfoCmd.Result, reportInfoCmd.Desc, reportInfoCmd.Data1, reportInfoCmd.Data2)
	}
	return model.ROOM_ERROR_SUCCESS
}

func (c *Collider)HandleReportStats(message []byte, rwc io.ReadWriteCloser) (int) {
	Xlog.Debug("HandleReportInfo into")

	// 1. 验证客户是否已经在房间和token是否已经过时
	var reportStats model.ReportStatsCmd
	err := model.ReportStatsUnMarshal(message, &reportStats)
	if err != nil {
		Xlog.Errorf("Cmd:%s -> ReportInfoUnMarshal failed: %s", model.CMD_REPORT_STATS, err.Error())
		return handleClientMessageFormatFail(model.CMD_REPORT_STATS, rwc)
	}
	// 查找房间是否存在
	r := c.roomTable.GetRoomLocked(reportStats.RoomID)
	if r != nil {
		return  r.ReportStats(&reportStats, rwc)
	} else {
		Xlog.Errorf("Cmd:%s -> Client:%s can't find the room:%s", reportStats.Cmd, reportStats.UID, reportStats.RoomID)
		return handleFindRoomFail(model.CMD_REPORT_STATS, rwc)
	}
}

// ice server cmd handle
func (c *Collider)HandleIceRegister(message []byte, ip string) (int) {
	Xlog.Info("HandleIceRegister into")

	var iceReg  iceBalance.IceRegisterCmd
	err := iceBalance.IceRegisterUnMarshal(message, &iceReg)
	if err != nil {
		Xlog.Errorf("Cmd:%s -> IceRegisterUnMarshal failed: %s", iceBalance.CMD_ICE_REGISTER, err.Error())
		return iceBalance.ICE_JSON_PARSE_FALIED
	}
	iceConf := iceBalance.IceRegister2IceConf(&iceReg)
	Xlog.Infof("ip -> act:%s, conf:%s", ip, iceConf.Ip)

	return iceBalance.IceTab.CreateIce(iceConf)
}

func (c *Collider)HandleIceReportIceRxTxRate(message []byte) (int) {
	Xlog.Debug("HandleIceReportIceRxTxRate into")

	var iceReport  iceBalance.IceReportRxTxRateCmd
	err := iceBalance.IceReportRxTxRateUnMarshal(message, &iceReport)
	if err != nil {
		Xlog.Errorf("Cmd:%s -> IceReportRxTxRateUnMarshal failed: %s", iceBalance.CMD_ICE_REGISTER, err.Error())
		return iceBalance.ICE_JSON_PARSE_FALIED
	}

	return iceBalance.IceTab.UpdateIceRxTxRate(iceReport.Ip, iceReport.RxRate, iceReport.TxRate)
}

func (c *Collider)HandleIceDeregister(message []byte) (int) {
	Xlog.Info("HandleIceDeregister into")

	var iceDereg  iceBalance.IceDeregisterCmd
	err := iceBalance.IceDeregisterUnMarshal(message, &iceDereg)
	if err != nil {
		Xlog.Errorf("Cmd:%s -> HandleIceDeregister failed: %s", iceBalance.CMD_ICE_REGISTER, err.Error())
		return iceBalance.ICE_JSON_PARSE_FALIED
	}

	iceBalance.IceTab.DeleteIce(iceDereg.Ip)
	return iceBalance.ICE_OK
}