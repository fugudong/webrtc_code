package main

import (
	"fmt"
	"github.com/json-iterator/go"
	"github.com/rs/xid"
	"golang.org/x/net/websocket"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type JoinCmd struct {
	Cmd      string `json:"cmd"`
	AppID    string `json:"appId"`
	ConnectType int `json:"connectType"`		// 默认为P2P方式
	Token    string `json:"token"`
	RoomID   string `json:"roomId"`
	RoomName string `json:"roomName"`
	UID      string `json:"uid"`
	Uname    string `json:"uname"`
	UserType int   	`json:"userType"`
	TalkType int   	`json:"talkType"`
	Time     string `json:"time"`		// 时间同样使用string，在处理时转换为int64
	IP       string `json:"ip"`
	OsName   string `json:"osName"`
	Browser   string `json:"browser"`
	SdkInfo  string `json:"sdkInfo"`
}

type LeaveCmd struct {
	Cmd       string `json:"cmd"`
	AppID     string `json:"appId"`
	RoomID    string `json:"roomId"`
	UID       string `json:"uid"`
	Uname    string `json:"uname"`
	UserType int   	`json:"userType"`
	Time      string `json:"time"`
}

type RoomUser struct {
	AppID     string `json:"appId"`
	UID       string `json:"uid"`
	Uname     string `json:"uname"`
	UserType int   	`json:"userType"`
	TalkType int   	`json:"talkType"`
}

type JoinRespCmd struct {
	Cmd      string `json:"cmd"`
	Result   int    `json:"result"`
	Desc     string `json:"desc"`
	ConnectType int `json:"connectType"`
	RoomID   string `json:"roomId"`
	RoomName string `json:"roomName"`
	UID      string `json:"uid"`
	Uname    string `json:"uname"`
	UserList [] RoomUser `json:"userList"`
}

// 根据房间范围进行测试
const RoomIdRange = 100
//
const UserIdRange = 100
// client数量,测试最大并发数量
const ClientMaxNumber =5000

// 生成user id, 可以先用随机数而不是uid进行测试，以便检验uid冲突
func createUid() string {
	value := rand.Intn(UserIdRange)
	str := strconv.Itoa(value)
	return str
}

//随机数生成room id
func createRoomId()  string {
	value := rand.Intn(RoomIdRange)
	str := strconv.Itoa(value)
	return str
}

func createRand(maxRange int )  int {
	value := rand.Intn(maxRange)
	return value
}
// 用户ID
func GenerateUserID() string {
	id := xid.New()
	return id.String()
}

// 房间ID
func GenerateRoomID() string {
	id := xid.New()
	return id.String()
}

func createJoin(uid string, roomId string)  [] byte{
	curTime := time.Now().Unix()
	strTime := strconv.FormatInt(curTime, 10)
	join := JoinCmd{
		Cmd:"join",
		AppID:"12345",
		ConnectType:0,
		Token:"1223232423432432432432535435345544444444444444444",
		RoomID:roomId,
		RoomName:roomId,
		UID:uid,
		Uname:uid,
		UserType:0,
		TalkType:1,
		Time:strTime,
		OsName:"hw os",
		Browser:"chrome windows",
		SdkInfo:"0voice rtc 1.0"}
	json,err := jsoniter.Marshal(join)
	if err != nil {
		log.Println("jsoniter.Marshal(join) failed")
	}
	return json
}

func createLeave(uid string, roomId string) [] byte  {
	curTime := time.Now().Unix()
	strTime := strconv.FormatInt(curTime, 10)
	leave := LeaveCmd{
		Cmd:"leave",
		AppID:"12345",
		RoomID:roomId,
		UID:uid,
		Uname:uid,
		UserType:0,
		Time:strTime}
	json,err := jsoniter.Marshal(leave)
	if err != nil {
		log.Println("jsoniter.Marshal(leave) failed")
	}
	return json
}

func JoinRespUnMarshal(data [] byte, respJoin *JoinRespCmd) ( error){
	err :=jsoniter.Unmarshal(data, respJoin)

	return  err
}

var ClientCountFaild int = 0
var ClientCountOk int = 0

func client()  {
	var wsurl = "wss://easywebrtc.com:8088/ws"
	var origin = "https://easywebrtc.com/"
	var sleepTime  time.Duration
	sleepTime = time.Duration(createRand(100))
	time.Sleep(time.Second * sleepTime)

	//var wsurl = "ws://129.204.197.215:9000/ws"
	//var origin = "http://129.204.197.215:9000"
	ws, err := websocket.Dial(wsurl, "", origin)
	if err != nil {
		ClientCountFaild++
		fmt.Printf("count->ok:%d,failed:%d , err:%s\n", ClientCountOk, ClientCountFaild, err.Error())
		time.Sleep(time.Second * 2)
		return
	}
	ClientCountOk++
	for ; ;  {
		var sleepTime  time.Duration
		sleepTime = time.Duration(createRand(50))
		time.Sleep(time.Second * sleepTime)
		// join
		joinByte := createJoin(GenerateUserID(), GenerateRoomID())
		n, err := ws.Write(joinByte)
		if n != len(joinByte) {
			fmt.Printf("only send %d byte, but it have %d byte", n, len(joinByte))
		}
		// 读取join结果
		data := make([]byte, 1024*10)
		n, err = ws.Read(data)
		if err != nil {
			fmt.Println(err)
		}
		// 解析join结果
		var respJoin JoinRespCmd
		JoinRespUnMarshal(data, &respJoin)
		//fmt.Println("resp: ", string(data))
		if respJoin.Result != 0 {
			fmt.Printf("join resp faild:%s", respJoin.Desc)
		}
		// 离开
		leaveByte := createLeave(respJoin.UID, respJoin.RoomID)
		n, err = ws.Write(leaveByte)
		if n != len(leaveByte) {
			fmt.Printf("only send %d byte, but it have %d byte", n, len(leaveByte))
		}

		data = make([]byte, 1024*10)
		n, err = ws.Read(data)
		if err != nil {
			fmt.Println(err)
		}

		//break

	}
	ws.Close()
}
func main()  {
	ClientCountFaild = 0
	ClientCountOk = 0
	for i := 0; i < ClientMaxNumber; i++ {
		go client()
	}

	for ; ;  {
		time.Sleep(time.Second * 1)
	}
}