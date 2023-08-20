package main

import (
	"encoding/json"
	"fmt"
	"github.com/json-iterator/go"
	"os"
)

func encodeJson() {
	type ColorGroup struct {
		ID     int
		Name   string
		Colors []string
	}
	group := ColorGroup{
		ID: 1,
		//Name:   "Reds",
		Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
	}
	b, err := json.Marshal(group)
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)
	fmt.Println("\n换行")

	var json_iterator = jsoniter.ConfigCompatibleWithStandardLibrary
	b, err = json_iterator.Marshal(group)
	os.Stdout.Write(b)
}

func decodeJson() {
	var jsonBlob = []byte(`[
        {"Name": "Platypus", "Order": "Monotremata"},
        {"Name": "Quoll",    "Order": "Dasyuromorphia"}
    ]`)
	type Animal struct {
		Name  string
		Order string
		Uid   string
	}
	var animals []Animal
	err := json.Unmarshal(jsonBlob, &animals)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v", animals)

	var animals2 []Animal
	var json_iterator = jsoniter.ConfigCompatibleWithStandardLibrary
	json_iterator.Unmarshal(jsonBlob, &animals2)
	fmt.Printf("%+v", animals2)
}

type Data struct {
	Type string      `json:"type"`
	Id   interface{} `json:"id"`
}

func decode(t string) {
	var x Data
	err := json.Unmarshal([]byte(t), &x)
	if err != nil {
		panic(err)
	}
	if x.Type == "a" {
		fmt.Println(x.Id.(string))
	} else {
		fmt.Println(x.Id.(float64)) //json解析中number默认作为float64解析
	}
}
func decodeJson2() {

	t1 := `{"type":"liao", id:"aaa"}`
	t2 := `{"protocol": "resp","cmd": "respJoin","result": 0,	"desc": "ok",
			"connectType":0}`
	t3 := `{
    "protocol": "resp",
	"cmd": "respJoin",
	"result": 0,
    "desc": "ok",
    "connectType": 0,
    "roomId": "21223232",
    "roomName": "test1",
    "uid": "12323232",
    "uname": "darren",
   	"userList": [
        {
            "uid": "2323232",
			"sessionId": "sadfsdfxsafsdf",
			"appId": "0123456789",
            "uname": "dongzige"            
        }
    ]
}
`
	//t2 := `{"type":"b", id:22222}`

	_byte := []byte(t1)
	json_data := jsoniter.Get(_byte, "type")
	fmt.Println("type", json_data.ToString())

	_byte2 := []byte(t2)
	json_data2 := jsoniter.Get(_byte2, "liao")
	fmt.Println("cmd", json_data2.ToString())
	if json_data2.ToString() == "" {
		fmt.Println("cmd is null")
	}

	_byte3 := []byte(t3)
	json_data3 := jsoniter.Get(_byte3, "userList")
	fmt.Println("userList", json_data3.ToString())
	// 怎么解析数组呢？
}
/**
功能要求：
建议使用https://tool.lu/json/ 进行json生成go 结构体
1. 解析复杂json，带数组json
2. 修改json内容
3. 封装复杂的json
 */


func main() {
	encodeJson()
	fmt.Printf("\n解码json\n")
	decodeJson2()
}

//---------------------
//作者：一蓑烟雨1989
//来源：CSDN
//原文：https://blog.csdn.net/wangshubo1989/article/details/78709802
//版权声明：本文为博主原创文章，转载请附上博文链接！
