package model

import "testing"

// 测试Candidate反序列化，完整数据输入
func Test_CandidateUnMarshal_1(t *testing.T) {
	// 输入
	candidateMsg := `{
    "cmd": "candidate",
	"appId": "123332323232",
    "roomId": "3232332",
    "uid": "32323232",
    "remoteUid": "12232323",
	"time": "1456449001000",
    "msg": {
        "candidateId": "fsdfsdf",
   	 	"label": 1,
    	"candidate":"xxsfsdfsdfsd"
		}
	}`

	_byte := []byte(candidateMsg)
	var candidateStruct CandidateCmd
	err :=CandidateUnMarshal(_byte, &candidateStruct)

	if (err != nil) {
		t.Error("CandidateUnMarshal error: ", err)
	} else {
		CandidatePrint(&candidateStruct)
	}
}

// 测试Candidate反序列化，缺失部分成分
func Test_CandidateUnMarshal_2(t *testing.T) {
	// 输入
	// "token": "01234567890123456789012345678901", 不输入测试结果
	// "connectType": 0,
	candidateMsg := `{
    "cmd": "candidate",
	"appId": "123332323232",
    "roomId": "3232332",
    "remoteUid": "12232323",
	"time": "1456449001000",
    "msg": {
        "sdp": "v=0\\r\\no=- 6871105499568....",
        "type": "candidate"
		}
	}`

	_byte := []byte(candidateMsg)
	var candidateStruct CandidateCmd
	err := CandidateUnMarshal(_byte, &candidateStruct)

	if (err == nil) {
		t.Error("CandidateUnMarshal error: ", err)
	} else {
		CandidatePrint(&candidateStruct)
	}
}

// 测试Candidate反序列化，增加部分不需要的成员
func Test_CandidateUnMarshal_3(t *testing.T) {
	// 输入 ,
	// 增加结构体没有的成分  "direct":1,
	candidateMsg := `{
    "cmd": "candidate",
	"appId": "123332323232",
    "roomId": "3232332",
    "uid": "32323232",
	"direct":1,
    "remoteUid": "12232323",
	"time": "1456449001000",
    "msg": {
        "candidateId": "fsdfsdf",
    	"label": 1,
    	"candidate":"xxsfsdfsdfsd"
		}
	}`
	//t2 := `{"type":"b", id:22222}`

	_byte := []byte(candidateMsg)
	var candidateStruct CandidateCmd
	err :=CandidateUnMarshal(_byte, &candidateStruct)

	if (err != nil) {
		t.Error("CandidateUnMarshal error: ", err)
	} else {
		CandidatePrint(&candidateStruct)
	}
}
// 测试序列化
func Test_CandidateRelayMarshal(t *testing.T) {
	// 输入
	candidateMsg := `{
    "cmd": "candidate",
	"appId": "123332323232",
    "roomId": "3232332",
    "uid": "32323232",
    "remoteUid": "12232323",
	"time": "1456449001000",
    "msg": {
         "candidateId": "fsdfsdf",
    	"label": 1,
    	"candidate":"xxsfsdfsdfsd"
		}
	}`

	_byte := []byte(candidateMsg)
	var candidateStruct CandidateCmd
	err :=CandidateUnMarshal(_byte, &candidateStruct)

	if (err != nil) {
		t.Error("CandidateUnMarshal error: ", err)
	} else {
		CandidatePrint(&candidateStruct)
	}

	respCandidate := CandidateRelayGenerate(&candidateStruct)
	if respCandidate.Cmd != CMD_RELAY_CANDIDATE {
		t.Error("respCandidate cmd:", respCandidate.Cmd, " is not equal to ", CMD_RELAY_CANDIDATE)
	}

	respCandidateByte, err := CandidateRespMarshal(&respCandidate)
	if err != nil {
		t.Error("CandidateResponseMarshal error")
	}
	t.Logf("respCandidate:%s\n", string(respCandidateByte))
}
