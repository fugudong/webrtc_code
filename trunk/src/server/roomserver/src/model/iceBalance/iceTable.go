package iceBalance

// 不考虑IP归属地，只记录每个ICE的统计信息，不同客户端的IP到底选择哪个ICE是由 geoIp模块进行处理
import (
	Xlog "github.com/cihub/seelog"
	"sync"
)

type IceTable struct {
	lock  	sync.Mutex		// s锁
	IceStats   map[string]*IceStat		// 房间
}
var IceTab *IceTable	// 将包含的iceserver发送给客户端

func NewIceTable() *IceTable {
	return &IceTable{IceStats: make(map[string]*IceStat)}
}

// 添加一个ICE地址
func (iceTable *IceTable) CreateIce(iceConf IceConf) int  {
	iceTable.lock.Lock()
	defer iceTable.lock.Unlock()

	if _, ok := iceTable.IceStats[iceConf.Ip]; ok {
		Xlog.Warnf("%s have already been created", iceConf.Ip)
		return ICE_HAVE_REGISTERED 		// 如果找到则直接返回
	}
	
	Xlog.Warnf("CreateIce -> ip:%s", iceConf.Ip)
	iceStat := NewIceStat(iceConf, iceTable)
	iceTable.IceStats[iceConf.Ip] = iceStat
	return ICE_OK
}

// 删除一个ICE地址
func (iceTable *IceTable) DeleteIce(ip string) {
	iceTable.lock.Lock()
	defer iceTable.lock.Unlock()
	Xlog.Infof("DeleteIce -> ip:%s", ip)

	if iceStat, ok := iceTable.IceStats[ip]; ok {
		iceStat.deregister()
		delete(iceTable.IceStats, ip)
	} else {
		Xlog.Warnf("DeleteIce -> can't find IceStats of the ip:%s", ip)
	}
}

func (iceTable *IceTable) TimeoutDeleteIce(ip string) {
	iceTable.lock.Lock()
	defer iceTable.lock.Unlock()
	Xlog.Infof("TimeoutDeleteIce -> ip:%s", ip)

	if iceStat, ok := iceTable.IceStats[ip]; ok {
		iceStat.deregister()
		delete(iceTable.IceStats, ip)
	} else {
		Xlog.Warnf("DeleteIce -> can't find IceStats of the ip:%s", ip)
	}
}

func (iceTable *IceTable) FindIce(country string, needBandwidth int) (*IceStat, int){
	iceTable.lock.Lock()
	defer iceTable.lock.Unlock()
	findIce := false
	var iceStat *IceStat
	// 先查找适合区域的ICE
	Xlog.Infof("FindIce -> try find the assigned country:%s server", country)
	for _, iceStat = range iceTable.IceStats {
		if iceStat.Country == country  {
			// 查看负载情况，如果负载已到最大
			if iceStat.IsMaxload(needBandwidth) {
				continue		// 如果已经到了最大负载则查找另一个ICE
			}
			findIce = true
			break
		}
	}
	if findIce {
		Xlog.Infof("FindIce -> found the %s server:%s", iceStat.Country, iceStat.IceConf.Ip)
		return  iceStat, ICE_OK
	}
	findIce = false
	// 查找非同一国家的服务器
	Xlog.Infof("FindIce -> try find the other country server")
	for _, iceStat = range iceTable.IceStats {
		if iceStat.Country != country {
			// 查看负载情况，如果负载已到最大
			if iceStat.IsMaxload(needBandwidth) {
				continue		// 如果已经到了最大负载则查找另一个ICE
			}
			findIce = true
			break
		}
	}
	if findIce {
		Xlog.Infof("FindIce -> found the %s server:%s", iceStat.Country, iceStat.IceConf.Ip)
		return iceStat, ICE_OK
	} else {
		Xlog.Infof("FindIce -> can't find the lightly loaded server")
		return nil, ICE_FULL_LOAD
	}
}

func (iceTable *IceTable) InsertIceClient(ip string, uid string) int {
	iceTable.lock.Lock()
	defer iceTable.lock.Unlock()

	if iceStat, ok := iceTable.IceStats[ip]; ok {
		Xlog.Infof("InsertIceClient -> ip:%s, uid:%s", ip, uid)
		iceStat.InsertIceClient(uid)		// 多人通话时将不正确，fix me
		return ICE_OK
	} else {
		Xlog.Warnf("InsertIceClient -> can't find IceStats of the ip:%s", ip)
		return ICE_NO_FOUND
	}
}

func (iceTable *IceTable) InsertTurnClient(ip string, uid string) int {
	iceTable.lock.Lock()
	defer iceTable.lock.Unlock()

	if iceStat, ok := iceTable.IceStats[ip]; ok {
		Xlog.Infof("InsertTurnClient -> ip:%s, uid:%s", ip, uid)
		iceStat.InsertTurnClient(uid)
		return ICE_OK
	} else {
		Xlog.Warnf("InsertTurnClient -> can't find IceStats of the ip:%s", ip)
		return ICE_NO_FOUND
	}
}

func (iceTable *IceTable) InsertStunClient(ip string, uid string) int {
	iceTable.lock.Lock()
	defer iceTable.lock.Unlock()

	if iceStat, ok := iceTable.IceStats[ip]; ok {
		Xlog.Infof("InsertStunClient -> ip:%s, uid:%s", ip, uid)
		iceStat.InsertStunClient(uid)
		return ICE_OK
	} else {
		Xlog.Warnf("InsertStunClient -> can't find IceStats of the ip:%s", ip)
		return ICE_NO_FOUND
	}
}

func (iceTable *IceTable) DeleteClient(ip string, uid string) {
	iceTable.lock.Lock()
	defer iceTable.lock.Unlock()

	if iceStat, ok := iceTable.IceStats[ip]; ok {
		Xlog.Infof("DeleteClient -> ip:%s, uid:%s", ip, uid)
		iceStat.DeleteClient(uid)
	} else {
		//Xlog.Warnf("DeleteClient -> can't find IceStats of the ip:%s", ip)
	}
}

func (iceTable *IceTable) GetIceClientNumber(ip string, uid string) int {
	iceTable.lock.Lock()
	defer iceTable.lock.Unlock()

	if iceStat, ok := iceTable.IceStats[ip]; ok {
		return iceStat.GetIceClientNumber()

	} else {
		return	-1	// 不存在则退出
	}
}
var TxRxCount = 0
// 更新ICE 当前带宽情况
func (iceTable *IceTable) UpdateIceRxTxRate(ip string, rxRate int, txRate int) int {
	iceTable.lock.Lock()
	defer iceTable.lock.Unlock()

	if iceStat, ok := iceTable.IceStats[ip]; ok {
		TxRxCount++
		if TxRxCount > 20 {
			TxRxCount = 0
			Xlog.Infof("UpdateIceRxTxRate -> ip[%s]:%s,rx:%dkbps, tx:%dkbps", iceStat.Country, ip, rxRate, txRate)
		} else {
			Xlog.Debugf("UpdateIceRxTxRate -> ip[%s]:%s,rx:%dkbps, tx:%dkbps", iceStat.Country, ip, rxRate, txRate)
		}
		iceStat.UpdateIceRxTxRate(rxRate, txRate)
		return ICE_OK// 找到已存在的client
	} else {
		Xlog.Warnf("UpdateIceRxTxRate -> can't find IceStats of the ip:%s", ip)
		return ICE_NO_FOUND
	}
}