package main

import (
	"encoding/json"
	"fmt"
)

type UserListJson struct {
	AppID     string `json:"appId"`
	SessionID string `json:"sessionId"`
	UID       string `json:"uid"`
	Uname     string `json:"uname"`
}

type Join struct {
	Cmd      string `json:"cmd"`
	Desc     string `json:"desc"`
	ConnectType     int    `json:"connectType"`
	Result   int    `json:"result"`
	RoomID   string `json:"roomId"`
	RoomName string `json:"roomName"`
	UID      string `json:"uid"`
	Uname    string `json:"uname"`
	UserList [] UserListJson  `json:"userList"`
}
/**
1. 封装json
 */
func encodeJoinJson() {
	user1 := [] UserListJson {
		{
			AppID: "13232323",
			SessionID: "xssfsdafsdf",
			UID:"xiafd",
			Uname:"zhangli",
		},
		{
			AppID: "11111",
			SessionID: "aaaa",
			UID:"xxxx",
			Uname:"liaoqingfu",
		},
	}
	//var user2  [2] UserListJso
	joinInfo := Join{
		Cmd: "join",
		UserList: user1	,
	}
	data, err := json.Marshal(joinInfo)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(data)
	fmt.Println(string(data))
}
func main()  {
	encodeJoinJson()
}