package logic

import (
	"roomserver/src/model"
	"roomserver/src/model/iceBalance"
	"roomserver/src/util"
	"fmt"
	Xlog "github.com/cihub/seelog"
)

func getNeedBandwidth(talkType int) int {
	switch talkType {
	case  model.TALK_TYPE_AUDIO_ONLY:
		return iceBalance.SingleAudioBandwidth
	case model.TALK_TYPE_AUDIO_VIDEO:
		return iceBalance.SingleAudioVideoBandwidth
	case model.TALK_TYPE_VIDEO_ONLY:
		return iceBalance.SingleAudioVideoBandwidth
	case model.TALK_TYPE_NO_AUDIO_VIDEO:
		return 0
	}
	return 0
}
//// 规则，一 地区， 二 负载情况
func  SmartSelectIce(src *Client, dst *Client) (int, string, model.RtcConfiguration)  {
	var rtcConfig model.RtcConfiguration
	var iceStat * iceBalance.IceStat
	var ret int
	srcCountry := util.GeoIpTable.GeoIpGetCountry(src.ClientIp)
	dstCountry := util.GeoIpTable.GeoIpGetCountry(dst.ClientIp)
	Xlog.Infof("src ip:%s, country:%s; dst ip:%s, country:%s", src.ClientIp, srcCountry, dst.ClientIp, dstCountry)
	needBandwidth := 0
	needBandwidth += getNeedBandwidth(src.TalkType)
	needBandwidth += getNeedBandwidth(dst.TalkType)
	// 先简单选择
	if (srcCountry == ICE_USA) ||  (dstCountry == ICE_USA) {
		// 只要一个客户端在美国就优选选择美国的ICE server
		iceStat, ret = iceBalance.IceTab.FindIce(ICE_USA, needBandwidth)
		if ret == iceBalance.ICE_OK {
			Xlog.Infof("find the US ICE:%s", iceStat.IceConf.Ip)
		} else {
			Xlog.Warn("can't find the US ICE")
		}
	} else if (srcCountry == ICE_CN) &&  (srcCountry == ICE_CN) {
		iceStat, ret = iceBalance.IceTab.FindIce(ICE_CN, needBandwidth)
	} else {
		// 没有美国用户时则选择香港中继节点
		iceStat, ret = iceBalance.IceTab.FindIce(ICE_HK, needBandwidth)
	}
	if iceStat == nil {
		return  ret,"",rtcConfig
	}

	// 填充rtcConfig
	rtcConfig.BundlePolicy = iceStat.IceConf.BundlePolicy
	// TURN
	turn1 := fmt.Sprintf("turn:%s:%d?transport=udp", iceStat.IceConf.Ip,
		iceStat.IceConf.TurnPort)
	turn2 := fmt.Sprintf("turn:%s:%d?transport=tcp", iceStat.IceConf.Ip,
		iceStat.IceConf.TurnPort)
	iceTurnServer := model.IceServer{
		Credential:iceStat.IceConf.Credential,
		Username:iceStat.IceConf.UserName,
	}
	iceTurnServer.Urls = append(iceTurnServer.Urls, turn1)
	iceTurnServer.Urls = append(iceTurnServer.Urls, turn2)
	// STUN
	stun := fmt.Sprintf("stun:%s:%d", iceStat.IceConf.Ip,
		iceStat.IceConf.StunPort)
	iceStunServer := model.IceServer{}
	iceStunServer.Urls = append(iceStunServer.Urls, stun)
	// 加入TURN 和 STUN server
	rtcConfig.IceServers = append(rtcConfig.IceServers, iceTurnServer)
	rtcConfig.IceServers = append(rtcConfig.IceServers, iceStunServer)

	rtcConfig.IceTransportPolicy = iceStat.IceConf.IceTransportPolicy
	rtcConfig.RtcpMuxPolicy = iceStat.IceConf.RtcpMuxPolicy
	// 填充client
	iceIP := iceStat.IceConf.Ip

	return 	iceBalance.ICE_OK, iceIP, rtcConfig
}