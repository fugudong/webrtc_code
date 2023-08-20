package model

import "testing"

// 测试KeepLive反序列化，完整数据输入
func Test_KeepLiveUnMarshal_1(t *testing.T) {
	// 输入
	keepLiveMsg := `{
    "cmd": "keepLive",
	"appId": "123332323232",
    "roomId": "3232332",
    "uid": "32323232",
	"time": "1456449001000"
	}`

	_byte := []byte(keepLiveMsg)
	var keepLiveStruct KeepLiveCmd
	err := KeepLiveUnMarshal(_byte, &keepLiveStruct)

	if (err != nil) {
		t.Error("KeepLiveUnMarshal error: ", err)
	} else {
		KeepLivePrint(&keepLiveStruct)
	}
}

// 测试KeepLive反序列化，缺失部分成分
func Test_KeepLiveUnMarshal_2(t *testing.T) {
	// 输入
	// "appId": "123332323232",, 不输入测试结果
	keepLiveMsg := `{
    "cmd": "keepLive",
    "roomId": "3232332",
    "uid": "32323232",
	"time": "1456449001000"
	}`

	_byte := []byte(keepLiveMsg)
	var keepLiveStruct KeepLiveCmd
	err := KeepLiveUnMarshal(_byte, &keepLiveStruct)

	if (err == nil) {
		t.Error("KeepLiveUnMarshal error: ", err)
	} else {
		KeepLivePrint(&keepLiveStruct)
	}
}

// 测试KeepLive反序列化，增加部分不需要的成员
func Test_KeepLiveUnMarshal_3(t *testing.T) {
	// 输入 ,
	// 增加结构体没有的成分  "direct":1,
	keepLiveMsg := `{
    "cmd": "keepLive",
	"appId": "123332323232",
    "roomId": "3232332",
	"direct":1,
    "uid": "32323232",
	"time": "1456449001000"
	}`

	_byte := []byte(keepLiveMsg)
	var keepLiveStruct KeepLiveCmd
	err := KeepLiveUnMarshal(_byte, &keepLiveStruct)

	if (err != nil) {
		t.Error("KeepLiveUnMarshal error: ", err)
	} else {
		KeepLivePrint(&keepLiveStruct)
	}
}

// 测试序列化
func Test_KeepLiveRespMarshal(t *testing.T) {
	// 输入
	keepLiveMsg := `{
    "cmd": "keepLive",
	"appId": "123332323232",
    "roomId": "3232332",
    "uid": "32323232",
	"time": "1456449001000"
	}`

	_byte := []byte(keepLiveMsg)
	var keepLiveStruct KeepLiveCmd
	err :=KeepLiveUnMarshal(_byte, &keepLiveStruct)

	if (err != nil) {
		t.Error("KeepLiveUnMarshal error: ", err)
	} else {
		KeepLivePrint(&keepLiveStruct)
	}

	// 测试relay
	respKeepLive := KeepLiveRespGenerate(&keepLiveStruct)
	if respKeepLive.Cmd != CMD_RESP_KEEP_LIVE {
		t.Error("respKeepLive cmd:", respKeepLive.Cmd, " is not equal to ", CMD_RESP_KEEP_LIVE)
	}

	respKeepLiveByte, err := KeepLiveRespMarshal(&respKeepLive)
	if err != nil {
		t.Error("KeepLiveRelayMarshal error")
	}

	t.Logf("respKeepLive:%s\n", string(respKeepLiveByte))
}

