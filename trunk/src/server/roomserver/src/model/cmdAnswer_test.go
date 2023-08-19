package model

import "testing"

// 测试Answer反序列化，完整数据输入
func Test_AnswerUnMarshal_1(t *testing.T) {
	// 输入
	answerMsg := `{
    "cmd": "answer",
	"appId": "123332323232",
    "roomId": "3232332",
    "uid": "32323232",
    "uname": "qingfu",
    "remoteUid": "12232323",
	"time": "1456449001000",
    "msg": {
        "sdp": "v=0\\r\\no=- 6871105499568....",
        "type": "answer"
		}
	}`

	_byte := []byte(answerMsg)
	var answerStruct AnswerCmd
	err :=AnswerUnMarshal(_byte, &answerStruct)

	if (err != nil) {
		t.Error("AnswerUnMarshal error: ", err)
	} else {
		AnswerPrint(&answerStruct)
	}
}

// 测试Answer反序列化，缺失部分成分
func Test_AnswerUnMarshal_2(t *testing.T) {
	// 输入
	// "token": "01234567890123456789012345678901", 不输入测试结果
	// "connectType": 0,
	answerMsg := `{
    "cmd": "answer",
	"appId": "123332323232",
    "roomId": "3232332",
    "uname": "qingfu",
    "remoteUid": "12232323",
	"time": "1456449001000",
    "msg": {
        "sdp": "v=0\\r\\no=- 6871105499568....",
        "type": "answer"
		}
	}`

	_byte := []byte(answerMsg)
	var answerStruct AnswerCmd
	err := AnswerUnMarshal(_byte, &answerStruct)

	if (err == nil) {
		t.Error("AnswerUnMarshal error: ", err)
	} else {
		AnswerPrint(&answerStruct)
	}
}

// 测试Answer反序列化，增加部分不需要的成员
func Test_AnswerUnMarshal_3(t *testing.T) {
	// 输入 ,
	// 增加结构体没有的成分  "direct":1,
	answerMsg := `{
    "cmd": "answer",
	"appId": "123332323232",
    "roomId": "3232332",
    "uid": "32323232",
	"direct":1,
    "uname": "qingfu",
    "remoteUid": "12232323",
	"time": "1456449001000",
    "msg": {
        "sdp": "v=0\\r\\no=- 6871105499568....",
        "type": "answer"
		}
	}`
	//t2 := `{"type":"b", id:22222}`

	_byte := []byte(answerMsg)
	var answerStruct AnswerCmd
	err := AnswerUnMarshal(_byte, &answerStruct)

	if (err != nil) {
		t.Error("AnswerUnMarshal error: ", err)
	} else {
		AnswerPrint(&answerStruct)
	}
}
// 测试序列化
func Test_AnswerRespMarshal(t *testing.T) {
	// 输入
	answerMsg := `{
    "cmd": "answer",
	"appId": "123332323232",
    "roomId": "3232332",
    "uid": "32323232",
    "uname": "qingfu",
    "remoteUid": "12232323",
	"time": "1456449001000",
    "msg": {
        "sdp": "v=0\\r\\no=- 6871105499568....",
        "type": "answer"
		}
	}`

	_byte := []byte(answerMsg)
	var answerStruct AnswerCmd
	err :=AnswerUnMarshal(_byte, &answerStruct)

	if (err != nil) {
		t.Error("AnswerUnMarshal error: ", err)
	} else {
		AnswerPrint(&answerStruct)
	}

	respAnswer := AnswerRelayGenerate(&answerStruct)
	if respAnswer.Cmd != CMD_RELAY_ANSWER {
		t.Error("respAnswer cmd:", respAnswer.Cmd, " is not equal to ", CMD_RELAY_ANSWER)
	}

	respAnswerByte, err := AnswerRespMarshal(&respAnswer)
	if err != nil {
		t.Error("AnswerResponseMarshal error")
	}
	t.Logf("respAnswer:%s\n", string(respAnswerByte))
}
