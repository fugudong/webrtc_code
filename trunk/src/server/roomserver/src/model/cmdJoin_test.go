package model

import "testing"

// 测试Join反序列化，完整数据输入
func Test_JoinUnMarshal_1(t *testing.T) {
	// 输入
	joinMsg := `{
    "cmd": "join",
    "appId": "0123456789",
    "connectType": 0,
    "token": "01234567890123456789012345678901",
	"audio":false,
	"roomId": "124565",
	"roomName": "test1",
	"uid": "124565",
    "uname": "124565",
	"time": "121322332",
    "ip": "118.7.123.2",
    "sdkInfo": "web_rtc_1.0",
	"osName": "windows chrome"
}`
	//t2 := `{"type":"b", id:22222}`

	_byte := []byte(joinMsg)
	var joinStruct JoinCmd
	err :=JoinUnMarshal(_byte, &joinStruct)


	if (err != nil) {
		t.Error("JoinUnMarshal error: ", err)
	} else {
		JoinPrint(&joinStruct)
	}
}

// 测试Join反序列化，缺失部分成分
func Test_JoinUnMarshal_2(t *testing.T) {
	// 输入
	// "token": "01234567890123456789012345678901", 不输入测试结果
	// "connectType": 0,
	joinMsg := `{
    "cmd": "join",
    "appId": "0123456789",
    
	"audio":false,
	"roomId": "124565",
	"roomName": "test1",
	"uid": "124565",
    "uname": "124565",
	"time": "121322332",
    "ip": "118.7.123.2",
    "sdkInfo": "web_rtc_1.0",
	"osName": "windows chrome"
}
`
	//t2 := `{"type":"b", id:22222}`

	_byte := []byte(joinMsg)
	var joinStruct JoinCmd
	err :=JoinUnMarshal(_byte, &joinStruct)

	if (err != nil) {
		t.Error("JoinUnMarshal error: ", err)
	} else {
		JoinPrint(&joinStruct)
	}
	if joinStruct.Token != "" {
		t.Error("token parse  error: ", err)
	}
	if joinStruct.Audio != false {
		t.Error("audio parse  error: ", err)
	}
}

// 测试Join反序列化，增加以及缺失部分成分
func Test_JoinUnMarshal_3(t *testing.T) {
	// 输入
	// 缺失"token": "01234567890123456789012345678901", 不输入测试结果
	// 缺失"connectType": 0,
	// 增加结构体没有的成分  "direct":1,
	joinMsg := `{
    "cmd": "join",
    "appId": "0123456789",
    "direct":1,
	"audio":false,
	"roomId": "124565",
	"roomName": "test1",
	"uid": "124565",
    "uname": "124565",
	"time": "121322332",
    "ip": "118.7.123.2",
    "sdkInfo": "web_rtc_1.0",
	"osName": "windows chrome"
}
`
	//t2 := `{"type":"b", id:22222}`

	_byte := []byte(joinMsg)
	var joinStruct JoinCmd
	err :=JoinUnMarshal(_byte, &joinStruct)

	if (err != nil) {
		t.Error("JoinUnMarshal error: ", err)
	} else {
		JoinPrint(&joinStruct)
	}
	if joinStruct.Token != "" {
		t.Error("token parse  error: ", err)
	}
	if joinStruct.Audio != false {
		t.Error("audio parse  error: ", err)
	}
}
// 测试序列化
func Test_JoinRespMarshal(t *testing.T) {
	// 输入
	joinMsg := `{
    "cmd": "join",
    "appId": "0123456789",
    "connectType": 0,
    "token": "01234567890123456789012345678901",
	"audio":false,
	"roomId": "124565",
	"roomName": "test1",
	"uid": "124565",
    "uname": "124565",
	"time": "121322332",
    "ip": "118.7.123.2",
    "sdkInfo": "web_rtc_1.0",
	"osName": "windows chrome"
}
`
	//t2 := `{"type":"b", id:22222}`

	_byte := []byte(joinMsg)
	var joinStruct JoinCmd
	err :=JoinUnMarshal(_byte, &joinStruct)


	if (err != nil) {
		t.Error("JoinUnMarshal error: ", err)
	} else {
		JoinPrint(&joinStruct)
	}
	userList := [] RoomUser {
		{
			AppID: "13232323",
			UID:"xiafd",
			Uname:"zhangli",
		},
		{
			AppID: "11111",
			UID:"xxxx",
			Uname:"liaoqingfu",
		},
	}
	respJoin := JoinRespGenerate(&joinStruct, userList,0, "ok")
	if respJoin.Cmd != CMD_RESP_JOIN {
		t.Error("respJoin cmd:", respJoin.Cmd, " is not equal to ", CMD_RESP_JOIN)
	}

	respJoinByte, err := JoinRespMarshal(&respJoin)
	if err != nil {
		t.Error("JoinResponseMarshal error")
	}
	t.Logf("respJoin:%s\n", string(respJoinByte))
}
