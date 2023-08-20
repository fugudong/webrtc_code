package model

import "testing"

// 测试Leave反序列化，完整数据输入
func Test_LeaveUnMarshal_1(t *testing.T) {
	// 输入
	leaveMsg := `{
    "cmd": "leave",
	"appId": "123332323232",
    "roomId": "3232332",
    "uid": "32323232",
	"time": "1456449001000"
	}`

	_byte := []byte(leaveMsg)
	var leaveStruct LeaveCmd
	err := LeaveUnMarshal(_byte, &leaveStruct)

	if (err != nil) {
		t.Error("LeaveUnMarshal error: ", err)
	} else {
		LeavePrint(&leaveStruct)
	}
}

// 测试Leave反序列化，缺失部分成分
func Test_LeaveUnMarshal_2(t *testing.T) {
	// 输入
	// "appId": "123332323232",, 不输入测试结果
	leaveMsg := `{
    "cmd": "leave",
    "roomId": "3232332",
    "uid": "32323232",
	"time": "1456449001000"
	}`

	_byte := []byte(leaveMsg)
	var leaveStruct LeaveCmd
	err := LeaveUnMarshal(_byte, &leaveStruct)

	if (err == nil) {
		t.Error("LeaveUnMarshal error: ", err)
	} else {
		LeavePrint(&leaveStruct)
	}
}

// 测试Leave反序列化，增加部分不需要的成员
func Test_LeaveUnMarshal_3(t *testing.T) {
	// 输入 ,
	// 增加结构体没有的成分  "direct":1,
	leaveMsg := `{
    "cmd": "leave",
	"appId": "123332323232",
    "roomId": "3232332",
	"direct":1,
    "uid": "32323232",
	"time": "1456449001000"
	}`

	_byte := []byte(leaveMsg)
	var leaveStruct LeaveCmd
	err := LeaveUnMarshal(_byte, &leaveStruct)

	if (err != nil) {
		t.Error("LeaveUnMarshal error: ", err)
	} else {
		LeavePrint(&leaveStruct)
	}
}
// 测试序列化
func Test_LeaveRelayMarshal(t *testing.T) {
	// 输入
	leaveMsg := `{
    "cmd": "leave",
	"appId": "123332323232",
    "roomId": "3232332",
    "uid": "32323232",
	"time": "1456449001000"
	}`

	_byte := []byte(leaveMsg)
	var leaveStruct LeaveCmd
	err :=LeaveUnMarshal(_byte, &leaveStruct)

	if (err != nil) {
		t.Error("LeaveUnMarshal error: ", err)
	} else {
		LeavePrint(&leaveStruct)
	}

	// 测试relay
	relayLeave := LeaveRelayGenerate(&leaveStruct)
	if relayLeave.Cmd != CMD_RELAY_LEAVE {
		t.Error("respLeave cmd:", relayLeave.Cmd, " is not equal to ", CMD_RELAY_LEAVE)
	}

	relayLeaveByte, err := LeaveRelayMarshal(&relayLeave)
	if err != nil {
		t.Error("LeaveRelayMarshal error")
	}

	t.Logf("relayLeave:%s\n", string(relayLeaveByte))
}

// 测试序列化
func Test_LeaveRespMarshal(t *testing.T) {
	// 输入
	leaveMsg := `{
    "cmd": "leave",
	"appId": "123332323232",
    "roomId": "3232332",
    "uid": "32323232",
	"time": "1456449001000"
	}`

	_byte := []byte(leaveMsg)
	var leaveStruct LeaveCmd
	err :=LeaveUnMarshal(_byte, &leaveStruct)

	if (err != nil) {
		t.Error("LeaveUnMarshal error: ", err)
	} else {
		LeavePrint(&leaveStruct)
	}

	// 测试relay
	respLeave := LeaveRespGenerate(&leaveStruct)
	if respLeave.Cmd != CMD_RESP_LEAVE {
		t.Error("respLeave cmd:", respLeave.Cmd, " is not equal to ", CMD_RESP_LEAVE)
	}

	respLeaveByte, err := LeaveRespMarshal(&respLeave)
	if err != nil {
		t.Error("LeaveRelayMarshal error")
	}

	t.Logf("respLeave:%s\n", string(respLeaveByte))
}

