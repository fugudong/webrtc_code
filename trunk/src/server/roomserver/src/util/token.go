package util

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)
//const appId             = "C5D15F8FD394285DA5227B533302A518" //App ID
//const appCertificate    = "fe1a0437bf217bdd34cd65053fb0fe1d" //App Certificate
//const expiredTime       = "1546271999" // 授权时间戳
//const userId           = "2323232x" //客户端定义的用户 ID
const DefaultAppCertificate    = "fe1a0437bf217bdd34cd65053fb0fe1d"
type TokenRequest struct {
	AppID string  `json:"appId"`
	UID string  `json:"uid"`
}

type TokenInfo struct {
	Ver string `json:"ver"`
	AppID string  `json:"appId"`
	Hash string  `json:"hash"`
	Expired string `json:"expired"`
}


func GenerateSignalingToken(userId string, appID string, appCertificate string, expiredTsInSeconds int64) (string, error) {

	version := "1"
	expired := fmt.Sprint(expiredTsInSeconds)	// 转成字符串
	content := userId + appID + appCertificate + expired

	hasher := md5.New()
	hasher.Write([]byte(content))
	md5sum := hex.EncodeToString(hasher.Sum(nil))

	token := &TokenInfo{
		version,
		appID,
		md5sum,
		expired,
	}

	buf, err := json.Marshal(token)
	if err != nil {
		return "marshal error", err
	}

	return string(buf), err
}

func CheckSignalingToken(tokenStr string, userId string, appCertificate string) (int64, error) {
	var token TokenInfo
	err := json.Unmarshal([]byte(tokenStr), &token)	// 解析出结构体
	if err != nil {
		return 0, errors.New("Unmarshal token failed")
	}
	// 校验字MD5
	content := userId + token.AppID + appCertificate + token.Expired

	hasher := md5.New()
	hasher.Write([]byte(content))
	md5sum := hex.EncodeToString(hasher.Sum(nil))
	if md5sum == token.Hash {
		expired, err := strconv.ParseInt(token.Expired,10, 64)
		if err != nil {
			return 0, errors.New("expired time parse failed")
		}
		return expired, nil
	}
	return 0,errors.New("token check failed")
}