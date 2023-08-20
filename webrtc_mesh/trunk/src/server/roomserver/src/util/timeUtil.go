package util

import "time"

// 获取毫秒时间戳，线性增长
func GetMillisecond()  int64 {
	m := time.Now().UnixNano() / 1e6
	return m
}