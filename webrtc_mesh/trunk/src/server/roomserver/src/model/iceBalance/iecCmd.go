package iceBalance

import "github.com/json-iterator/go"

// ICE注册到中心节点
const CMD_ICE_REGISTER = "iceRegister"
// 取消注册
const CMD_ICE_DEREGISTER = "iceDeregister"
// 客户端报告带宽负载
const CMD_ICE_REPORT_RX_TX_RATE = "iceReportRxTxRate"

type IceRegisterCmd struct {
	Cmd 	string	`json:"cmd"`// 命令
	Ip 		string	`json:"ip"`
	IceType string	`json:"iceType"`
	BundlePolicy string	`json:"bundlePolicy"`
	RtcpMuxPolicy      string `json:"rtcpMuxPolicy"` // "negotiate" | "require"
	IceTransportPolicy string `json:"iceTransportPolicy"`	// relay, all
	TurnPort 	int	`json:"turnPort"`
	StunPort 	int	`json:"stunPort"`
	MaxBandwidth int `json:"maxBandwidth"`		// 最大带宽
	Username 	string	`json:"username"`
	Credential 	string	`json:"credential"`
}

type IceDeregisterCmd struct {
	Cmd string	`json:"cmd"`// 命令
	Ip  string	`json:"ip"`// ICE服务器外网ip
}

type IceReportRxTxRateCmd struct {
	Cmd string	`json:"cmd"`// 命令
	Ip  string	`json:"ip"`// ICE服务器外网ip
	RxRate int	`json:"rxRate"`// 下行实时带宽负载
	TxRate int  `json:"txRate"`// 上行实时带宽负载
}



func IceRegisterUnMarshal(data [] byte, iceReg *IceRegisterCmd) ( error){
	err :=jsoniter.Unmarshal(data, iceReg)

	if err != nil {
		return  err
	}

	return  err
}

func IceDeregisterUnMarshal(data [] byte, iceDereg *IceDeregisterCmd) ( error){
	err :=jsoniter.Unmarshal(data, iceDereg)

	if err != nil {
		return  err
	}

	return  err
}

func  IceRegister2IceConf(iceReg *IceRegisterCmd) IceConf {
	var iceConf IceConf
	iceConf.Ip 			= iceReg.Ip
	iceConf.IceType 	= iceReg.IceType
	iceConf.BundlePolicy	= iceReg.BundlePolicy
	iceConf.IceTransportPolicy	= iceReg.IceTransportPolicy
	iceConf.RtcpMuxPolicy	= iceReg.RtcpMuxPolicy
	iceConf.TurnPort	= iceReg.TurnPort
	iceConf.StunPort	= iceReg.StunPort
	iceConf.UserName	= iceReg.Username
	iceConf.Credential	= iceReg.Credential
	iceConf.MaxBandwidth = iceReg.MaxBandwidth

	return iceConf
}

func IceReportRxTxRateUnMarshal(data [] byte,  iceReport *IceReportRxTxRateCmd) ( error){
	err :=jsoniter.Unmarshal(data, iceReport)

	if err != nil {
		return  err
	}

	return  err
}