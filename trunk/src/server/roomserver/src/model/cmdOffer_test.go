package model

import "testing"

// 测试Offer反序列化，完整数据输入
func Test_OfferUnMarshal_1(t *testing.T) {
	// 输入
	offerMsg := `{
    "cmd": "offer",
	"appId": "123332323232",
    "roomId": "3232332",
    "uid": "32323232",
    "uname": "qingfu",
    "remoteUid": "12232323",
	"time": "1456449001000",
    "msg": {
        "sdp": "v=0\\r\\no=- 6871105499568....",
        "type": "offer"
		}
	}`

	_byte := []byte(offerMsg)
	var offerStruct OfferCmd
	err :=OfferUnMarshal(_byte, &offerStruct)

	if (err != nil) {
		t.Error("OfferUnMarshal error: ", err)
	} else {
		OfferPrint(&offerStruct)
	}
}

// 测试Offer反序列化，缺失部分成分
func Test_OfferUnMarshal_2(t *testing.T) {
	// 输入
	// "token": "01234567890123456789012345678901", 不输入测试结果
	// "connectType": 0,
	offerMsg := `{
    "cmd": "offer",
	"appId": "123332323232",
    "roomId": "3232332",
    "uname": "qingfu",
    "remoteUid": "12232323",
	"time": "1456449001000",
    "msg": {
        "sdp": "v=0\\r\\no=- 6871105499568....",
        "type": "offer"
		}
	}`

	_byte := []byte(offerMsg)
	var offerStruct OfferCmd
	err := OfferUnMarshal(_byte, &offerStruct)

	if (err == nil) {
		t.Error("OfferUnMarshal error: ", err)
	} else {
		OfferPrint(&offerStruct)
	}
}

// 测试Offer反序列化，增加部分不需要的成员
func Test_OfferUnMarshal_3(t *testing.T) {
	// 输入 ,
	// 增加结构体没有的成分  "direct":1,
	offerMsg := `{
    "cmd": "offer",
	"appId": "123332323232",
    "roomId": "3232332",
    "uid": "32323232",
	"direct":1,
    "uname": "qingfu",
    "remoteUid": "12232323",
	"time": "1456449001000",
    "msg": {
        "sdp": "v=0\\r\\no=- 6871105499568....",
        "type": "offer"
		}
	}`
	//t2 := `{"type":"b", id:22222}`

	_byte := []byte(offerMsg)
	var offerStruct OfferCmd
	err :=OfferUnMarshal(_byte, &offerStruct)

	if (err != nil) {
		t.Error("OfferUnMarshal error: ", err)
	} else {
		OfferPrint(&offerStruct)
	}
}
// 测试序列化
func Test_OfferRespMarshal(t *testing.T) {
	// 输入
	offerMsg := `{
    "cmd": "offer",
	"appId": "123332323232",
    "roomId": "3232332",
    "uid": "32323232",
    "uname": "qingfu",
    "remoteUid": "12232323",
	"time": "1456449001000",
    "msg": {
        "sdp": "v=0\\r\\no=- 6871105499568....",
        "type": "offer"
		}
	}`

	_byte := []byte(offerMsg)
	var offerStruct OfferCmd
	err :=OfferUnMarshal(_byte, &offerStruct)

	if (err != nil) {
		t.Error("OfferUnMarshal error: ", err)
	} else {
		OfferPrint(&offerStruct)
	}

	respOffer := OfferRelayGenerate(&offerStruct)
	if respOffer.Cmd != CMD_RELAY_OFFER {
		t.Error("respOffer cmd:", respOffer.Cmd, " is not equal to ", CMD_RELAY_OFFER)
	}

	respOfferByte, err := OfferRespMarshal(&respOffer)
	if err != nil {
		t.Error("OfferResponseMarshal error")
	}
	t.Logf("respOffer:%s\n", string(respOfferByte))
}
