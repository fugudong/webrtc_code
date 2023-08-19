package logic

import (
	"roomserver/src/model"
	"roomserver/src/util"
	Xlog "github.com/cihub/seelog"
	"io"
	"time"
)

var ClientKeepLiveTimeout time.Duration =  30

type RemoteClient struct{
	UID string	// 用户id
	Uname string 	// 用户名
	UserIp		string
	SubSessionId string
	BeginTime 	time.Time	// 在offer或者answer时写入起始时间
	ConnectTime time.Time	// 在connect时写入连接时间
	FinishTime	time.Time	// 在leave时写入结束时间
	//DifferTime  time.Duration	// 记录和join的时差
	ConnectType	string
	IceIP		string		// 分配的ICE地址
	RtcConfig   model.RtcConfiguration
	IsPeerConnected bool	// 是否已经通话成功
}
// 客户信息
type Client struct {
	AppID string
	UID string	// 用户id
	Uname string 	// 用户名
	UserType int  	// 用户类型
	TalkType int   // 通话类型
	ConnectType string
	OrigDevType int
	NewDevType	 int	// 以bit为单位
	// rwc is the interface to access the websocket connection.
	// It is set after the client registers with the server.
	rwc io.ReadWriteCloser
	// timer is used to remove this client if unregistered after a timeout.
	timer          *time.Timer	//
	JoinTime       time.Time
	ServJoinTime   time.Time	// 超时时用来推算实际的结束时间
	SessionId      string
	ClientIp       string		// 对应客户端外网ip地址
	RemoteClients  map[string] *RemoteClient	//对应remote_id，比如SubSessionIds[remote_id]
	SubSessionIds  [] string
	room           *Room
	collider       *Collider
	IsTimeoutLeave bool		// 是否为超时离开
}

func EstablishClientMapRelations(src *Client, dst *Client, iceIp string, rtcConfig *model.RtcConfiguration)  {
	curTime := time.Now()
	// 新建一个RemoteClient
	subsessionId1 := util.GenerateTalkSubsessionID()
	offsetTime := curTime.Sub(src.ServJoinTime)
	rc1 := RemoteClient{
		UID:dst.UID,
		Uname:dst.Uname,
		UserIp:dst.ClientIp,
		SubSessionId:subsessionId1,
		IceIP:iceIp,
		BeginTime:src.JoinTime.Add(offsetTime),
		ConnectTime:src.JoinTime.Add(offsetTime),
		FinishTime:src.JoinTime.Add(offsetTime),
		RtcConfig: (*rtcConfig),
		IsPeerConnected:false,
		ConnectType:"Unkown",
	}
	src.RemoteClients[dst.UID] = &rc1
	src.SubSessionIds = append(src.SubSessionIds, subsessionId1)

	subsessionId2 := util.GenerateTalkSubsessionID()
	rc2 := RemoteClient{
		UID:src.UID,
		Uname:src.Uname,
		UserIp:src.ClientIp,
		SubSessionId:subsessionId2,
		IceIP:iceIp,
		BeginTime:dst.JoinTime.Add(offsetTime),
		ConnectTime:dst.JoinTime.Add(offsetTime),
		FinishTime:dst.JoinTime.Add(offsetTime),
		RtcConfig: (*rtcConfig),
		IsPeerConnected:false,
		ConnectType:"Unkown",
	}
	dst.RemoteClients[src.UID] = &rc2
	dst.SubSessionIds = append(dst.SubSessionIds, subsessionId2)
}

func NewRemoteClient(uid string, uname string, beginTime time.Time,
	finishTime time.Time, iceIp string, rtcConfig model.RtcConfiguration) *RemoteClient  {
	return &RemoteClient{
		UID:uid,
		Uname:uname,
		BeginTime:beginTime,
		FinishTime:finishTime,
		IceIP:iceIp,
		RtcConfig:rtcConfig,
		IsPeerConnected:false,
	}
}
// 创建client
func NewClient(join *model.JoinCmd, room *Room, collider *Collider) *Client {
	devType := 0
	switch  join.TalkType  {
	case model.TALK_TYPE_AUDIO_ONLY:
		devType = 0x2
		break
	case model.TALK_TYPE_AUDIO_VIDEO:
		devType = 0x3
		break
	case model.TALK_TYPE_VIDEO_ONLY:
		devType = 0x1
		break
	case model.TALK_TYPE_NO_AUDIO_VIDEO:
		devType = 0
	}
	c := Client{AppID:join.AppID,
		UID:            join.UID,
		Uname:          join.Uname,
		UserType:       join.UserType,
		TalkType:       join.TalkType,
		ConnectType:    "",
		OrigDevType:    devType,	// 记录join时的devType开启情况
		NewDevType:     devType,		// 记录通话过程中的devType开启情况
		SessionId:      util.GenerateTalkSessionID(),
		ClientIp:       join.IP,		// 客户端外网IP
		room:           room,
		collider:       collider,
		IsTimeoutLeave: false,
		RemoteClients:  make(map[string]*RemoteClient)}

	c.timer = time.AfterFunc(ClientKeepLiveTimeout * time.Second , func() {
		c.collider.HandleTimeoutLeave(&c, room, c.rwc)
	})

	return &c
}

// register binds the ReadWriteCloser to the client if it's not done yet.
// 掉线之后要考虑重连的问题,
func (c *Client) register(rwc io.ReadWriteCloser) error {
	c.rwc = rwc
	return nil
}

// deregister closes the ReadWriteCloser if it exists.
func (c *Client) deregister() {
	Xlog.Debugf("client:%s deregister rwc.Close", c.UID)
	if c.rwc != nil {
		//c.rwc.Close()
		c.rwc = nil
	}
	if c.timer != nil {
		c.timer.Stop()
	}
}

// registered returns true if the client has registered.
func (c *Client) registered() bool {
	return c.rwc != nil
}

func (c *Client) resetTimer()  {
	//Xlog.Infof("client:%s resetTimer", c.UID)
	ret := c.timer.Reset(ClientKeepLiveTimeout * time.Second)
	if !ret {
		Xlog.Warnf("client uid:%s resetTimer failed", c.UID)
	}
}