package util

import (
	"fmt"
	"testing"
	"time"
)

func Test_token(t *testing.T) {
	userId := "2882341273"
	appId := "970CA35de60c44645bbae8a215061b33"
	appCertificate := "5CFd2fd1755d40ecb72977518be15d3b"
	now := time.Now().Unix()		//返回秒
	validTimeInSeconds := int64(3600*24)
	expiredTsInSeconds := now + validTimeInSeconds

	token, err:= GenerateSignalingToken(userId, appId, appCertificate, expiredTsInSeconds)
	if err != nil {
		t.Error("GenerateSignalingToken failed: ",err)
	}
	fmt.Println("token = " + token)


	expectedExpired, err := CheckSignalingToken(token, userId, appCertificate)
	if err != nil {
		t.Error("CheckSignalingToken failed: ", err)
	}
	if expectedExpired != expiredTsInSeconds {
		t.Error("expectedExpired is no equal to expiredTsInSeconds ")
	}
}

func Test_token2(t *testing.T) {
	userId := "2882341273"
	appId := "970CA35de60c44645bbae8"
	appCertificate := "5CFd2fd1755d40ecb72977518"
	now := time.Now().Unix()		//返回秒
	validTimeInSeconds := int64(3600*24)
	expiredTsInSeconds := now + validTimeInSeconds

	token, err:= GenerateSignalingToken(userId, appId, appCertificate, expiredTsInSeconds)
	if err != nil {
		t.Error("GenerateSignalingToken failed: ",err)
	}
	fmt.Println("token = " + token)

	expectedExpired, err := CheckSignalingToken(token, userId, appCertificate)
	if err != nil {
		t.Error("CheckSignalingToken failed: ", err)
	}
	if expectedExpired != expiredTsInSeconds {
		t.Error("expectedExpired is no equal to expiredTsInSeconds ")
	}

	expectedExpired, err = CheckSignalingToken(token, "liaoqingfu", appCertificate)
	if err == nil {
		t.Error("CheckSignalingToken failed: ", err)
	}
	if expectedExpired == expiredTsInSeconds {
		t.Error("expectedExpired is no equal to expiredTsInSeconds ")
	}
}