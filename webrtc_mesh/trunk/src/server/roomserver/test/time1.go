package main

import (
	"fmt"
	"time"
)

func main()  {
	now := time.Now()
	local1, err1 := time.LoadLocation("") //等同于"UTC"
	if err1 != nil {
		fmt.Println(err1)
	}
	local2, err2 := time.LoadLocation("Local")//服务器设置的时区
	if err2 != nil {
		fmt.Println(err2)
	}
	local3, err3 := time.LoadLocation("America/Los_Angeles")
	if err3 != nil {
		fmt.Println(err3)
	}

	fmt.Println(now.In(local1))
	fmt.Println(now.In(local2))
	fmt.Println(now.In(local3))

	var tm time.Time
	var tm2 time.Time
	fmt.Println(tm)

	if tm == tm2 {
		fmt.Printf("tm == tm2")
	} else {
		fmt.Printf("tm != tm2")
	}
}

