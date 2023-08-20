package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"net/smtp"
	"strings"
)

type Todo struct {
	Task1 string
	Task2 string
	Task3 string
	Task4 string
	Task5 string
}
type senderInfo struct {
	User      string `jsob:user`
	Passwd    string `json:passwd`
	Host_port string `json:host_port`
	Mailaddr  string `json:mailaddr`
	Subject   string `json:subject`
}

func (self *senderInfo) SendMail(toList, body string) error {
	//head := fmt.Sprintf("To: %v\r\nSubject: %v\r\nContent-Type: text/plain;charset=UTF-8\r\n\r\n",
	//	toList, self.Subject)
	head := fmt.Sprintf("To: %v\r\nSubject: %v\r\nContent-Type: text/html;charset=UTF-8\r\n\r\n",
		toList, self.Subject)
	host := strings.Split(self.Host_port, ":")
	if len(host) != 2 {
		return fmt.Errorf("%v not a valid host_port", self.Host_port)
	}
	auth := smtp.PlainAuth("", self.User, self.Passwd, host[0])
	return smtp.SendMail(self.Host_port, auth, self.Mailaddr,
		strings.Split(toList, ";"), []byte(head+body))
}
func main() {
	tmpl := template.Must(template.ParseFiles("index.html"))
	todos := []Todo{
		{"Learn Go", "liao", "qing", "fu","li"},
		{"t Go", "liao", "qing", "fu","li"},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, struct{ Todos []Todo }{todos}, )
	})
	b := new(bytes.Buffer)
	tmpl.Execute(b, struct{ Todos []Todo }{todos})
	result := string(b.String())
	testMail := &senderInfo{
		User:"webim@maoyt.com",
		Passwd:"s9fE$8*etc2m#di0",
		Host_port:"smtp.dowindns.com:25",
		Mailaddr:"webim@maoyt.com",
		Subject:"报警"	}
	err := testMail.SendMail("liaoqingfu@maoyt.com", result)
	fmt.Printf("html:%s", tmpl.DefinedTemplates())
	http.ListenAndServe(":8080", nil)
}
