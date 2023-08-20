package iceBalance

import (
	"sync"
	"roomserver/src/util"
	Xlog "github.com/cihub/seelog"
	"time"
)
var SingleAudioVideoBandwidth = 200	// 单向300kb
var SingleAudioBandwidth = 40	// 单向300kb

var IceKeepLiveTimeOut time.Duration = 30

type IceConf struct {
	Ip           		string // 对应的ice ip地址
	IceType				string
	BundlePolicy 		string 	// "balanced" | "max-compat" | "max-bundle";
	IceTransportPolicy 	string 	// relay, all
	RtcpMuxPolicy      	string  // "negotiate" | "require"
	TurnPort     int
	StunPort     int
	UserName     string
	Credential   string
	MaxBandwidth int    // 最大带宽，单位kb
}

type IceStat struct {
	IceConf 	IceConf		// 配置信息
	Country 	string		// 国家
	Subdivision string  	// 省或州
	City		string		// 城市
	//UsedBandwidth  		int  // 已用带宽，单位kb，当已用带宽达到MaxBandwidth时则说明服务器带宽已满，再连入用户
	RxRate		int			//  实时下行网络负载 单位kb
	TxRate		int			// 实时上行网络负载
	IceClients map[string] string		// 总共人数，key和value都是uid				使用map的方式存储，计算的时候也使用map的len计算
	TurnClients map[string] string		// 使用中继的人数
	StunClients map[string] string		//
	timer  *time.Timer	//
	lock  	sync.Mutex		// 锁
}

func NewIceStat(iceConf IceConf, iceTable *IceTable) *IceStat {
	// 调用geoip获取 Country  Subdivision City
	Xlog.Infof("NewIceStat -> iceConf.Ip:%s", iceConf.Ip)
	country := util.GeoIpTable.GeoIpGetCountry(iceConf.Ip)
	//subdivision := util.GeoIpTable.GeoIpGetSubdivisions(iceConf.Ip)
	//city := util.GeoIpTable.GeoIpGetCity(iceConf.Ip)
	subdivision :=""
	city		:= ""
	Xlog.Infof("iceConf country:%s, subdivision:%s, city:%s, timeout:%d", country, subdivision, city, IceKeepLiveTimeOut)
	iceStat:=IceStat{IceConf:iceConf,
		Country:country,
		Subdivision:subdivision,
		City:city,
		IceClients: make(map[string]string),
		TurnClients:make(map[string]string),
		StunClients:make(map[string]string)}
	iceStat.timer = time.AfterFunc(IceKeepLiveTimeOut * time.Second, func() {
		iceTable.TimeoutDeleteIce(iceStat.IceConf.Ip)
	})

	return &iceStat
}
// 是否为最大负载
func (iceStat *IceStat) IsMaxload(needBandwidth int) bool {
	Xlog.Infof("ICE Ip:%s, Bandwidth max:%dkbps, tx:%dkbps, rx:%dkbps, need:%dkpbs, Clients ice:%d, turn:%d, stun:%d",
		iceStat.IceConf.Ip,
		iceStat.IceConf.MaxBandwidth,
		iceStat.TxRate,
		iceStat.RxRate,
		needBandwidth,
		len(iceStat.IceClients),
		len(iceStat.TurnClients),
		len(iceStat.StunClients))
	if iceStat.TxRate + needBandwidth < iceStat.IceConf.MaxBandwidth {
		return false
	} else {
		return true
	}
}

// 增加一个用户id
func (iceStat *IceStat) InsertIceClient(uid string)  {
	iceStat.lock.Lock()
	defer iceStat.lock.Unlock()
	iceStat.IceClients[uid] = uid
}

// 增加一个TURN用户id
func (iceStat *IceStat) InsertTurnClient(uid string)  {
	iceStat.lock.Lock()
	defer iceStat.lock.Unlock()
	iceStat.TurnClients[uid] = uid
}

// 增加一个STUN用户id
func (iceStat *IceStat) InsertStunClient(uid string)  {
	iceStat.lock.Lock()
	defer iceStat.lock.Unlock()
	iceStat.StunClients[uid] = uid
}
// 删除一个用户，同时操作三个ICE、TURN、STUN 3个map
func (iceStat *IceStat) DeleteClient(uid string)  {
	iceStat.lock.Lock()
	defer iceStat.lock.Unlock()
	delete(iceStat.IceClients, uid)
	delete(iceStat.TurnClients, uid)
	delete(iceStat.StunClients, uid)
}

// 总共有多少个用户在使用该ICE服务器
func (iceStat *IceStat) GetIceClientNumber() int {
	iceStat.lock.Lock()
	defer iceStat.lock.Unlock()
	return len(iceStat.IceClients)
}

// 该ICE服务器总共有多少个用户采用了TURN中继的方式
func (iceStat *IceStat) GetTurnClientNumber() int {
	iceStat.lock.Lock()
	defer iceStat.lock.Unlock()
	return len(iceStat.TurnClients)
}

// 该ICE服务器总共有多少个用户采用了STUN P2P的方式
func (iceStat *IceStat) GetStunClientNumber() int {
	iceStat.lock.Lock()
	defer iceStat.lock.Unlock()
	return len(iceStat.StunClients)
}

func (iceStat *IceStat) UpdateIceRxTxRate(rxRate int, txRate int)  {
	iceStat.lock.Lock()
	defer iceStat.lock.Unlock()
	iceStat.RxRate = rxRate
	iceStat.TxRate = txRate
	iceStat.resetTimer()
}

func (iceStat *IceStat) resetTimer()  {
	Xlog.Debugf("Ice ip:%s resetTimer  = %d", iceStat.IceConf.Ip, IceKeepLiveTimeOut)
	ret := iceStat.timer.Reset(IceKeepLiveTimeOut * time.Second)
	if !ret {
		Xlog.Errorf("Ice ip:%s resetTimer failed", iceStat.IceConf.Ip)
	}
}

func (iceStat *IceStat) deregister() {
	if iceStat.timer != nil {
		iceStat.timer.Stop()
	}
}