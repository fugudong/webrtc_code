package logic

import (
	"roomserver/src/model"
	"roomserver/src/model/iceBalance"
	"fmt"
	Xlog "github.com/cihub/seelog"
	"io"
	"strconv"
	"sync"
	"time"
)
// 当房间在指定的时间不存在client时关闭房间，设置为10秒

//房间容纳人数，在config.ini进行配置
var MaxRoomCapacity = MAX_ROOM_CAPACITY
type Room struct {
	Parent *RoomTable		// 所属房间组
	Id string				// 房间ID
	Name string				// 名字
	Clients map[string]*Client	// 用户列表
	registerTimeout time.Duration	// 保活，超时则释放房间
	lock  	sync.Mutex		// 锁
}

func NewRoom(p *RoomTable, id string, name string, to time.Duration) *Room {
	return &Room{Parent: p, Id: id, Clients: make(map[string]*Client), registerTimeout: to}
}

// 在房间查找用户，如果用户不存在则新建一个用户，如果房间已满则返回错误信息
func (rm *Room) client(join *model.JoinCmd, rwc io.ReadWriteCloser, collider *Collider) (*Client, int) {
	clientID := join.UID
	if c, ok := rm.Clients[clientID]; ok {
		return c, model.ROOM_ERROR_SUCCESS	// 找到已存在的client
	}
	if len(rm.Clients) >= MaxRoomCapacity{
		Xlog.Warnf("Room %s is full, max number is %d, can't add the client %s",
			rm.Id,
			MaxRoomCapacity,
			clientID)
		return nil, model.ROOM_ERROR_FULL	// 房间已经满了，不能再加入新的client
	}

	rm.Clients[clientID] = NewClient(join, rm, collider)
	rm.Clients[clientID].register(rwc)

	Xlog.Debugf("Added client %s to room %s", clientID, rm.Id)

	return rm.Clients[clientID], model.ROOM_ERROR_SUCCESS
}

// remove closes the client connection and removes the client specified by the |clientID|.
func (rm *Room) remove(clientID string) {
	if c, ok := rm.Clients[clientID]; ok {
		c.deregister()	// 此时用户将关闭websocket连接
		c = c	// 仅是为了不报错，这里先不自己关闭连接
		delete(rm.Clients, clientID)
	}
}

// empty returns true if there is no client in the room.
func (rm *Room) ClientsNumber() int {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	return len(rm.Clients)
}

func (rm *Room) wsCount() int {
	count := 0
	for _, c := range rm.Clients {
		if c.registered() {
			count += 1
		}
	}
	return count
}
// 锁room，还是锁client，如果锁的是client，同房间其他client到来，在处理其他clients的时候出现竞态
// 1. 查找client，如果不存在则创建
func (rm *Room) Join(join *model.JoinCmd, rwc io.ReadWriteCloser, collider *Collider) int {
	rm.lock.Lock()				// 锁住房间
	defer func() {
		defer rm.lock.Unlock()		// 解锁房间
		Xlog.Info("Room Join leave")
	}()
	Xlog.Info("Room Join into")
	var msgByte [] byte
	var err error

	// 1. 查找client是否已经存在
	c, ret := rm.client(join, rwc, collider)		// 查找到对应的client
	if ret != model.ROOM_ERROR_SUCCESS {
		errDesc :=  model.ROOM_ERROR_MSG[ret]
		if ret == model.ROOM_ERROR_FULL {
			errDesc = fmt.Sprintf("room is full, it up to max number:%d", MaxRoomCapacity)
		}
		msgByte,err = model.JoinRespGenerateAndMarshal(join, nil, ret, errDesc)
		Send(rwc, msgByte)
		if (err != nil)  {
			Xlog.Error("JoinRespGenerateAndMarshal failed")
		}
		return ret
	}
	if err = c.register(rwc); err != nil {
		Xlog.Errorf("Join error:%s", err)
		msgByte,err = model.JoinRespGenerateAndMarshal(join, nil,1, err.Error())
		Send(rwc, msgByte)
		if (err != nil) {
			Xlog.Error("JoinRespGenerateAndMarshal failed")
		}
		return model.ROOM_ERROR_WEBSOCKET_BROKEN
	}
	// 生成session
	//model.JoinGenerateSession(join)
	Xlog.Infof("Client %s join room %s, room's people is %d", c.UID, rm.Id, len(rm.Clients))
	findIce := iceBalance.ICE_OK
	// 获取join时间
	curTime, err := strconv.ParseInt(join.Time, 10, 64)
	if err != nil {
		Xlog.Errorf("Parse join.time:%s failed: %s",
			join.Time,
			model.CMD_JOIN, err.Error())
		return handleClientMessageFormatFail(model.CMD_JOIN, rwc)
	}
	tm := time.Unix(curTime, 0)
	c.JoinTime = tm		// 保存起始时间
	c.ServJoinTime = time.Now()

	var iceIp string
	var rtcConfig model.RtcConfiguration
	// 2. 通知房间里的其他参与者有新人加入
	if len(rm.Clients) > 1 {	// 当不只是自己在房间
		// 通知其他人有新人加入
		users := make([]model.RoomUser,len(rm.Clients) - 1)
		i := 0
		for _, other := range rm.Clients {
			if c.UID != other.UID  {	// 通知房间的其他人有新人加入
				// 分析两者的ip，返回rtcConfig
				findIce, iceIp, rtcConfig= SmartSelectIce(c, other)
				if findIce != iceBalance.ICE_OK {
					// 报错
					msgByte, err = model.JoinRespGenerateAndMarshal(join, nil,
						model.ROOM_ERROR_ICE_FULL_LOADING,
						model.ROOM_ERROR_MSG[model.ROOM_ERROR_ICE_FULL_LOADING])
					Send(c.rwc, msgByte)
					if(err != nil) {
						Xlog.Error("JoinRespGenerateAndMarshal failed")
					}
					return model.ROOM_ERROR_ICE_FULL_LOADING
				}
				// 建立相互之间的通话关系
				EstablishClientMapRelations(c, other, iceIp, &rtcConfig)
				// ice添加人数
				iceBalance.IceTab.InsertIceClient(iceIp, other.UID)
				iceBalance.IceTab.InsertIceClient(iceIp, c.UID)
				// 通知已在房间的客户端有新人加入
				msgByte,err  := model.JoinNewPeerGenerateAndMarshal(join, rtcConfig)
				if (err != nil) {
					Xlog.Error("JoinNewPeerGenerateAndMarshal failed")
				}
				users[i].AppID = other.AppID
				users[i].UID = other.UID
				users[i].Uname = other.Uname
				users[i].UserType = other.UserType
				users[i].TalkType = other.TalkType
				i++
				Xlog.Infof("Notify %s that %s join %s", other.Uname, c.Uname, rm.Id)
				Send(other.rwc, msgByte)
			}
		}
		// 通知自己目前房间其他加入者的信息，但不包括自己
		msgByte, err = model.JoinRespGenerateAndMarshal(join, users, 0, "ok")
		Send(c.rwc, msgByte)
		if(err != nil) {
			Xlog.Error("JoinRespGenerateAndMarshal failed")
		}
	} else {	// 只有自己在房间
		msgByte, err = model.JoinRespGenerateAndMarshal(join, nil, 0, "ok")
		Send(c.rwc, msgByte)
		if(err != nil) {
			Xlog.Error("JoinRespGenerateAndMarshal failed")
		}
	}
	err = joinOperateDb(join, c)
	if(err != nil) {
		Xlog.Error("joinOperateDb failed")
		// 邮件通知报错？
	}
	return model.ROOM_ERROR_SUCCESS
}

/**
leave leave命令
rwc 对应的websocket
status 离开类型，0：正常关闭，1为超时关闭
 */
func (rm *Room) Leave(leave *model.LeaveCmd, rwc io.ReadWriteCloser, timeout bool) int {
	rm.lock.Lock()
	defer func() {
		Xlog.Info("Room Leave leave")
		rm.lock.Unlock()		// 解锁房间
	}()
	Xlog.Info("Room Leave into")
	var msgByte [] byte
	var err error
	c, ok := rm.Clients[leave.UID]			// 查找对应的client
	if !ok {
		// 如果不存在在返回不存在的结果
		Xlog.Warnf("In room:%s can't find the uid:%s", rm.Id, leave.UID)
		msgByte,err = model.LeaveRespGenerateAndpMarshal(leave, model.ROOM_ERROR_NOT_FIND_UID,
			"In room "+ rm.Id + " can't find the client " + leave.UID)
		Send(rwc, msgByte)
		if (err != nil) {
			Xlog.Error("LeaveRespGenerateAndpMarshal failed")
		}
		return model.ROOM_ERROR_NOT_FIND_UID
	}

	c.IsTimeoutLeave = timeout
	// 通知房间的其他client有人退出
	if len(rm.Clients) > 1 {
		msgByte,err = model.LeaveRelayGenerateAndMarshal(leave)
		for _, other := range rm.Clients {
			if c.UID != other.UID  {
				// 记录finish time
				rc, ok := c.RemoteClients[other.UID]
				if ok {
					curTime := time.Now()
					offsetTime := curTime.Sub(c.ServJoinTime)
					rc.FinishTime = c.JoinTime.Add(offsetTime)
					iceBalance.IceTab.DeleteClient(rc.IceIP, rc.UID)	// 将其从ice删除
					// 写数据库
					err = LeaveSubOperateDb(leave, c, rc)
					// 解除关系
					delete(c.RemoteClients, other.UID)
				}
				rc, ok = other.RemoteClients[c.UID]
				if ok {
					curTime := time.Now()
					offsetTime := curTime.Sub(other.ServJoinTime)
					rc.FinishTime = other.JoinTime.Add(offsetTime)
					iceBalance.IceTab.DeleteClient(rc.IceIP, rc.UID)	// 将其从ice删除
					// 写数据库
					err = LeaveSubOperateDb(leave, other, rc)
					// 解除关系
					delete(other.RemoteClients, c.UID)
				}
				Send(other.rwc, msgByte)
			}
		}
	}

	// 响应client，此时可以自己关闭RTCPeerConnection了
	msgByte,err = model.LeaveRespGenerateAndpMarshal(leave, 0, "ok")
	if (err == nil) {
		Send(c.rwc, msgByte)
	}

	// 写数据库
	err = LeaveOperateDb(leave, c)
//room rm r
//client c
	rm.remove(leave.UID)	// 删除client并释放client的资源
	//delete(rm.Clients, leave.UID)				// 删除client
	Xlog.Infof("Client %s leave room %s, remaining number is %d", leave.UID, rm.Id, len(rm.Clients))
	return model.ROOM_ERROR_SUCCESS
}

/**
leave leave命令
rwc 对应的websocket
status 离开类型，0：正常关闭，1为超时关闭
 */
func (rm *Room) JoinLeave(leave *model.LeaveCmd, rwc io.ReadWriteCloser) int {
	rm.lock.Lock()
	defer func() {
		defer rm.lock.Unlock()		// 解锁房间
		Xlog.Info("Room JoinLeave leave")
	}()
	Xlog.Info("Room JoinLeave into")
	c, ok := rm.Clients[leave.UID]			// 查找对应的client
	if !ok {
		return model.ROOM_ERROR_NOT_FIND_UID
	}

	// 通知房间的其他client有人退出
	if len(rm.Clients) > 1 {
		msgByte,err := model.LeaveRelayGenerateAndMarshal(leave)
		if err != nil {
			Xlog.Warnf("LeaveRelayGenerateAndMarshal failed:%s", err.Error())
		}
		for _, other := range rm.Clients {
			if c.UID != other.UID  {
				Send(other.rwc, msgByte)
			}
		}
	}
	// 如果代码执行到这里则说明出现了rejoin（重连的情况），可以记录到数据库中
	rm.remove(leave.UID)	// 删除client并释放client的资源
	Xlog.Debugf("Removed exist client %s from room %s", leave.UID, rm.Id)
	return model.ROOM_ERROR_SUCCESS
}

/**
1. 检测uid的client是否在房间内，如果不存在则报错
2. 检测房间中的RemoteUID是否存在
3. 正常流程走完则不用再发送响应信息
 */
func (rm *Room) Offer(offer *model.OfferCmd, rwc io.ReadWriteCloser) int {
	rm.lock.Lock()
	defer func() {
		defer rm.lock.Unlock()		// 解锁房间
		Xlog.Info("Room Offer leave")
	}()
	Xlog.Infof("Room Offer into")
	var msgByte [] byte
	var err error

	//1. 检测uid的client是否在房间内，如果不存在则报错
	c, ok := rm.Clients[offer.UID]			// 查找对应的client
	if !ok {
		// 如果不存在在返回不存在的结果
		Xlog.Warn("In room "+ rm.Id + " can't find the client " + offer.UID)
		msgByte, err = model.OfferRespGenerateAndMarshal(offer, 1, "In room "+ rm.Id + " can't find the client " + offer.UID, "")
		Send(rwc, msgByte)
		if (err != nil) {
			Xlog.Error("OfferRespGenerateAndMarshal failed")
		}
		return model.ROOM_ERROR_NOT_FIND_UID
	}

	//2. 检测房间中的RemoteUID是否存在,如果存在则发送offer，如果不存在resp client并给出报错信息
	sendOffer := false
	// 通知房间的其他client
	for _, other := range rm.Clients {
		if offer.RemoteUID == other.UID   {
			rc, ok := other.RemoteClients[c.UID]			// 查找对应的client
			if !ok {
				msgByte, err = model.OfferRespGenerateAndMarshal(offer, 1, "In room "+ rm.Id + " can't find the RemoteClients " + offer.RemoteUID, "")
				sendOffer = false
			} else {
				msgByte,err = model.OfferRelayGenerateAndMarshal(offer, rc.RtcConfig, rc.SubSessionId)
				sendOffer = true
			}

			if (err == nil) {
				Send(other.rwc, msgByte)
			} else {
				Xlog.Error("OfferRelayGenerateAndMarshal failed")
			}
			break
		}
	}
	if !sendOffer {
		msgByte, err = model.OfferRespGenerateAndMarshal(offer, model.ROOM_ERROR_NOT_FIND_REMOTEID,
			"In room "+ rm.Id + " can't find the remote client " + offer.RemoteUID, "")
		if (err == nil) {
			Send(c.rwc, msgByte)
		}else {
			Xlog.Error("OfferRespGenerateAndMarshal failed")
		}
		return model.ROOM_ERROR_NOT_FIND_REMOTEID
	} else {
		rc, ok := c.RemoteClients[offer.RemoteUID]			// 查找对应的client
		if !ok {
			msgByte, err = model.OfferRespGenerateAndMarshal(offer, 1, "In room "+ rm.Id + " can't find the RemoteClients " + offer.UID, "")
		} else {
			msgByte, err = model.OfferRespGenerateAndMarshal(
				offer,
				model.ROOM_ERROR_SUCCESS,
				"ok",
				rc.SubSessionId)
			curTime := time.Now()
			offsetTime := curTime.Sub(c.ServJoinTime)
			rc.BeginTime = c.JoinTime.Add(offsetTime)
			rc.ConnectTime = rc.BeginTime
			Xlog.Infof("uid:%s begin try talk to %s at time:%s",
				c.UID,
				rc.UID,
				rc.BeginTime.Format("2006/1/2 15:04:05") )
		}
		if (err == nil) {
			Send(c.rwc, msgByte)
		}else {
			Xlog.Error("OfferRespGenerateAndMarshal failed")
		}
		return model.ROOM_ERROR_SUCCESS
	}
}

/**
1. 检测uid的client是否在房间内，如果不存在则报错
2. 检测房间中的RemoteUID是否存在
3. 正常流程走完则不用再发送响应信息
 */
func (rm *Room) Answer(answer *model.AnswerCmd, rwc io.ReadWriteCloser) int {
	rm.lock.Lock()
	defer func() {
		defer rm.lock.Unlock()		// 解锁房间
		Xlog.Info("Room Answer leave")
	}()
	Xlog.Infof("Room Answer into")
	var msgByte [] byte
	var err error

	//1. 检测uid的client是否在房间内，如果不存在则报错
	c, ok := rm.Clients[answer.UID]			// 查找对应的client
	if !ok {
		// 如果不存在在返回不存在的结果
		Xlog.Warn("In room "+ rm.Id + " can't find the client " + answer.UID)
		msgByte, err = model.AnswerRespGenerateAndMarshal(answer, 1, "In room "+ rm.Id + " can't find the client " + answer.UID)
		Send(rwc, msgByte)
		if (err != nil) {
			Xlog.Error("AnswerRespGenerateAndMarshal failed")
		}
		return model.ROOM_ERROR_NOT_FIND_UID
	}

	//2. 检测房间中的RemoteUID是否存在,如果存在则发送offer，如果不存在resp client并给出报错信息
	sendAnswer := false
	// 通知房间的其他client
	for _, other := range rm.Clients {
		if answer.RemoteUID == other.UID   {
			msgByte,err :=model.AnswerRelayGenerateAndMarshal(answer)
			if (err == nil) {
				Send(other.rwc, msgByte)
				sendAnswer = true
			} else {
				Xlog.Error("AnswerRelayGenerateAndMarshal failed")
			}
			break
		}
	}
	if !sendAnswer {

		msgByte, err = model.AnswerRespGenerateAndMarshal(answer, 2, "In room "+ rm.Id + " can't find the remote client " + answer.RemoteUID)
		if (err == nil) {
			Send(c.rwc, msgByte)
		}else {
			Xlog.Error("AnswerRespGenerateAndMarshal failed")
		}
		return model.ROOM_ERROR_NOT_FIND_REMOTEID
	} else {
		rc, ok := c.RemoteClients[answer.RemoteUID]			// 查找对应的client
		if ok {
			curTime := time.Now()
			offsetTime := curTime.Sub(c.ServJoinTime)
			rc.BeginTime = c.JoinTime.Add(offsetTime)
			rc.ConnectTime = rc.BeginTime
			Xlog.Infof("uid:%s begin try talk to %s at time:%s",
				c.UID,
				rc.UID,
				rc.BeginTime.Format("2006/1/2 15:04:05") )
		}
	}
	//3. 正常流程走完则不用再发送响应信息
	return model.ROOM_ERROR_SUCCESS
}

/**
1. 检测uid的client是否在房间内，如果不存在则报错
2. 检测房间中的RemoteUID是否存在
3. 正常流程走完则不用再发送响应信息
 */
func (rm *Room) Candidate(candidate *model.CandidateCmd, rwc io.ReadWriteCloser) int {
	rm.lock.Lock()
	defer func() {
		defer rm.lock.Unlock()		// 解锁房间
		Xlog.Debug("Room Candidate leave")
	}()
	Xlog.Debug("Room Candidate into")
	var msgByte [] byte
	var err error

	//1. 检测uid的client是否在房间内，如果不存在则报错
	c, ok := rm.Clients[candidate.UID]			// 查找对应的client
	if !ok {
		// 如果不存在在返回不存在的结果
		Xlog.Warn("In room "+ rm.Id + " can't find the client " + candidate.UID)
		msgByte, err = model.CandidateRespGenerateAndMarshal(candidate, 1, "In room "+ rm.Id + " can't find the client " + candidate.UID)
		if (err == nil) {
			Send(rwc, msgByte)
		} else {
			Xlog.Error("CandidateRespGenerateAndMarshal failed")
		}
		return model.ROOM_ERROR_NOT_FIND_UID
	}

	//2. 检测房间中的RemoteUID是否存在,如果存在则发送candidate，如果不存在resp client并给出报错信息
	sendCandidate := false
	// 通知房间的其他client
	for _, other := range rm.Clients {
		if candidate.RemoteUID == other.UID   {
			msgByte,err :=model.CandidateRelayGenerateAndMarshal(candidate)
			if (err == nil) {
				Send(other.rwc, msgByte)
				sendCandidate = true
			} else {
				Xlog.Error("CandidateRelayGenerateAndMarshal failed")
			}
			break
		}
	}
	if !sendCandidate {
		msgByte, err = model.CandidateRespGenerateAndMarshal(candidate, model.ROOM_ERROR_NOT_FIND_REMOTEID,
			"In room "+ rm.Id + " can't find the remote client " + candidate.RemoteUID)
		if (err == nil) {
			Send(c.rwc, msgByte)
		}else {
			Xlog.Error("CandidateRespGenerateAndMarshal failed")
		}
		return model.ROOM_ERROR_NOT_FIND_REMOTEID
	}
	//3. 正常流程走完则不用再发送响应信息
	return model.ROOM_ERROR_SUCCESS
}


/**
1. 检测uid的client是否在房间内，如果不存在则报错
2. 检测房间中的RemoteUID是否存在
3. 正常流程走完则不用再发送响应信息
 */
func (rm *Room) TurnTalkType(turnTalkType *model.TurnTalkTypeCmd, rwc io.ReadWriteCloser) int {
	rm.lock.Lock()
	defer func() {
		defer rm.lock.Unlock()		// 解锁房间
		Xlog.Info("Room TurnTalkType leave")
	}()
	Xlog.Info("Room TurnTalkType into")

	var msgByte [] byte
	var err error
	c, ok := rm.Clients[turnTalkType.UID]			// 查找对应的client
	if !ok {
		// 如果不存在在返回不存在的结果
		Xlog.Warnf("In room:%s can find clientid:%s", rm.Id, turnTalkType.UID)
		msgByte,err = model.TurnTalkTypeRespGenerateAndMarshal(turnTalkType, model.ROOM_ERROR_NOT_FIND_UID,
			"In room "+ rm.Id + " can't find the client " + turnTalkType.UID)
		if (err == nil) {
			Send(rwc, msgByte)
		}
		return model.ROOM_ERROR_NOT_FIND_UID
	}
	switch turnTalkType.Index {
	case 0:	// 处理摄像头
		if turnTalkType.Enable == true	 {
			// 开启摄像头
			c.NewDevType |= 0x1	// bit 0置位
		} else {
			c.NewDevType &= 0xfe	// bit 0清零
		}
		break
	case 1:	// 麦克风
		if turnTalkType.Enable == true	 {
			// 开启摄像头
			c.NewDevType |= 0x2	// bit 0置位
		} else {
			c.NewDevType &= 0xfd	// bit 0清零
		}
		break
	}
	// 通知房间的其他client
	if len(rm.Clients) > 1 {
		msgByte,err = model.TurnTalkTypeRelayGenerateAndMarshal(turnTalkType)
		for _, other := range rm.Clients {
			if c.UID != other.UID  {
				Send(other.rwc, msgByte)
			}
		}
	}

	return model.ROOM_ERROR_SUCCESS
}

func (rm *Room) KeepLive(keepLive *model.KeepLiveCmd, rwc io.ReadWriteCloser) int {
	//Xlog.Info("Room KeepLive into")
	rm.lock.Lock()
	defer rm.lock.Unlock()

	var msgByte [] byte
	var err error
	c, ok := rm.Clients[keepLive.UID]			// 查找对应的client
	if !ok {
		// 如果不存在在返回不存在的结果
		Xlog.Warnf("In room:%s can't find the uid:%s", rm.Id, keepLive.UID)
		msgByte,err = model.KeepLiveRespGenerateAndpMarshal(keepLive, model.ROOM_ERROR_NOT_FIND_UID,
			"In room "+ rm.Id + " can't find the client " + keepLive.UID)
		Send(rwc, msgByte)
		if (err != nil) {
			Xlog.Error("LeaveRespGenerateAndpMarshal failed")
		}
		return model.ROOM_ERROR_NOT_FIND_UID
	}

	// 更新定时器，
	c.resetTimer()
	// 响应client
	msgByte,err = model.KeepLiveRespGenerateAndpMarshal(keepLive, model.ROOM_ERROR_SUCCESS, "ok")
	if (err == nil) {
		Send(c.rwc, msgByte)
	}

	//Xlog.Debugf("KeepLive client %s to room %s", keepLive.UID, rm.Id)
	return model.ROOM_ERROR_SUCCESS
}

func notifyRemoteClientDevType(peerConnected *model.PeerConnectedCmd, c * Client, remote * Client)  {
	var turnTalkType model.TurnTalkTypeCmd
	turnTalkType.Cmd = model.CMD_TURN_TALK_TYPE
	turnTalkType.AppID = peerConnected.AppID
	turnTalkType.RoomID = peerConnected.RoomID
	turnTalkType.UID = peerConnected.UID
	turnTalkType.Uname = peerConnected.Uname
	turnTalkType.Time = peerConnected.Time
	Xlog.Infof("NewDevType:0x%x, OrigDevType:0x%x", c.NewDevType, c.OrigDevType)
	// 判断video是否关闭
	if (c.NewDevType & 0x1) != (c.OrigDevType & 0x1) {
		// 状态发生变化
		if (c.NewDevType & 0x1) == 0 {
			turnTalkType.Index = 0
			turnTalkType.Enable = false
			msgByte, err := model.TurnTalkTypeRelayGenerateAndMarshal(&turnTalkType)
			if err != nil {
				Xlog.Errorf("TurnTalkTypeRelayGenerateAndMarshal failed: %s",
					err.Error())
			} else {
				Send(remote.rwc, msgByte)
			}
		}
	}
	// 判断audio是否关闭
	if (c.NewDevType & 0x2) != (c.OrigDevType & 0x2) {
		// 状态发生变化
		if (c.NewDevType & 0x2) == 0 {
			turnTalkType.Index = 1
			turnTalkType.Enable = false
			msgByte, err := model.TurnTalkTypeRelayGenerateAndMarshal(&turnTalkType)
			if err != nil {
				Xlog.Errorf("TurnTalkTypeRelayGenerateAndMarshal failed: %s",
					err.Error())
			} else {
				Send(remote.rwc, msgByte)
			}
		}
	}
}

func (rm *Room) PeerConnected(peerConnected *model.PeerConnectedCmd, rwc io.ReadWriteCloser) int {
	rm.lock.Lock()
	defer func() {
		defer rm.lock.Unlock()		// 解锁房间
		Xlog.Info("Room PeerConnected leave")
	}()
	Xlog.Info("Room PeerConnected into")

	var msgByte [] byte
	var err error

	//1. 检测uid的client是否在房间内，如果不存在则报错
	c, ok := rm.Clients[peerConnected.UID]			// 查找对应的client
	if !ok {
		// 如果不存在在返回不存在的结果
		Xlog.Warn("In room "+ rm.Id + " can't find the client " + peerConnected.UID)
		msgByte, err = model.PeerConnecteRespGenerateAndMarshal(peerConnected, 1,
			"In room "+ rm.Id + " can't find the client " + peerConnected.UID)
		if (err == nil) {
			Send(rwc, msgByte)
		} else {
			Xlog.Error("PeerConnecteRespGenerateAndMarshal failed")
		}
		return model.ROOM_ERROR_NOT_FIND_UID
	}

	// 当前的DevType和刚join时不同时则通知remote id
	Xlog.Infof("NewDevType:0x%x, OrigDevType:0x%x", c.NewDevType, c.OrigDevType)
	if c.NewDevType != c.OrigDevType {
		remote, ok := rm.Clients[peerConnected.RemoteUID]			// 查找对应的client
		if ok {
			notifyRemoteClientDevType(peerConnected, c, remote)
		}
	}
	c.ConnectType = peerConnected.ConnectType
	rc, ok := c.RemoteClients[peerConnected.RemoteUID]			// 查找对应的client
	if ok {
		if peerConnected.ConnectType == "TURN" {
			iceBalance.IceTab.InsertTurnClient(rc.IceIP, c.UID)
			rc.ConnectType = "TURN"
		} else if  peerConnected.ConnectType == "STUN"{
			iceBalance.IceTab.InsertStunClient(rc.IceIP, c.UID)
			rc.ConnectType = "STUN"
		} else {
			iceBalance.IceTab.InsertTurnClient(rc.IceIP, c.UID)
			rc.ConnectType = "Unknown"	// 苹果的Safari不能检测出连接类型
		}
		rc.IsPeerConnected = true
		curTime, err := strconv.ParseInt(peerConnected.Time, 10, 64)
		if err != nil {
			Xlog.Errorf("Parse peerConnected.time:%s failed: %s",
				peerConnected.Time,
				model.CMD_PEER_CONNECTED, err.Error())
			return handleClientMessageFormatFail(model.CMD_PEER_CONNECTED, rwc)
		}
		tm := time.Unix(curTime, 0)
		rc.ConnectTime = tm		// 更新
		Xlog.Infof("uid:%s connect to %s at time:%s",
			c.UID,
			rc.UID,
			rc.ConnectTime.Format("2006/1/2 15:04:05") )
	} else  {
		Xlog.Errorf("can't find the RemoteClients of uid %", peerConnected.RemoteUID)
	}

	return model.ROOM_ERROR_SUCCESS
}

func (rm *Room) ReportStats(reportStats *model.ReportStatsCmd, rwc io.ReadWriteCloser) int {
	rm.lock.Lock()
	defer func() {
		defer rm.lock.Unlock()		// 解锁房间
		Xlog.Debug("Room ReportStats leave")
	}()

	var err error
	Xlog.Debug("Room ReportStats into")

	//1. 检测uid的client是否在房间内，如果不存在则报错
	c, ok := rm.Clients[reportStats.UID]			// 查找对应的client
	if !ok {
		// 如果不存在在返回不存在的结果
		Xlog.Warn("In room "+ rm.Id + " can't find the client " + reportStats.UID)
		return model.ROOM_ERROR_NOT_FIND_UID
	}

	rc, ok := c.RemoteClients[reportStats.RemoteUID]			// 查找对应的client
	if !ok {
		Xlog.Warnf("reportStats can't find RemoteUID:%s, UID:%s->RemoteUID:%s, SubSessionIds:%s",
			reportStats.RemoteUID,
			reportStats.UID,
			reportStats.RemoteUID,
			c.SubSessionIds)
		return model.ROOM_ERROR_SUCCESS
	} else {
		// 3.写数据库 记录客户端统计信息
		err = reportStatsOperateDb(reportStats, rc.SubSessionId)
		if err != nil {
			Xlog.Errorf("reportStatsOperateDb failed:%s", err.Error())
		}
	}


	return model.ROOM_ERROR_SUCCESS
}