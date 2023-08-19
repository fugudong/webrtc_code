package util

import (
	log "github.com/cihub/seelog"
	"fmt"
	"net/smtp"
	"strings"
)

type MailInfo struct {
	SendIp	  string	// 发送者的外网ip
	Sender    string	// 发送人邮箱地址
	Receivers string	// 收件人邮箱地址
	User      string	// 发件人邮箱登录用户名，通常和Sender一致
	Passwd    string	// 发件人邮箱登录密码
	SmtpHostAndPort string
	Subject   string
}

func (self *MailInfo) SendMail(body string) error {
	head := fmt.Sprintf("To: %v\r\nFrom: %v\r\nSubject: %v\r\nContent-Type: text/html;charset=UTF-8\r\n\r\n",
		self.Receivers, self.User, self.Subject)
	host := strings.Split(self.SmtpHostAndPort, ":")
	if len(host) != 2 {
		return log.Errorf("%v not a valid SmtpHostAndPort", self.SmtpHostAndPort)
	}
	auth := smtp.PlainAuth("", self.User, self.Passwd, host[0])
	return smtp.SendMail(self.SmtpHostAndPort, auth, self.Sender,
		strings.Split(self.Receivers, ";"), []byte(head+body))
}