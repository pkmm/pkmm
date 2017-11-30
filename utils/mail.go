package utils

import (
	"github.com/astaxie/beego/utils"
	"github.com/astaxie/beego"
	"fmt"
	"time"
)

var (
	sendCh chan *utils.Email
	config string
)

func init() {
	queueSize, _ := beego.AppConfig.Int("mail.queue_size")
	host := beego.AppConfig.String("mail.host")
	port, _ := beego.AppConfig.Int("mail.port")

	username := beego.AppConfig.String("mail.user")
	password := beego.AppConfig.String("mail.password")
	from := beego.AppConfig.String("mail.from")
	if port == 0 {
		port = 25
	}

	config = fmt.Sprintf(`{"username":"%s","password":"%s","host":"%s","port":%d,"from":"%s"}`, username, password, host, port, from)
	sendCh = make(chan *utils.Email, queueSize)

	// 消费者，好好干活
	go func() {
		for {
			select {
			case m, ok := <-sendCh:
				if !ok {
					return
				}
				if err := m.Send(); err != nil {
					beego.Error("Send email:", err.Error())
				}

			}
		}
	}()
}

func SendMail(address, name, subject, content string, cc []string) bool {
	mail := utils.NewEMail(config)
	mail.To = []string{address}
	mail.Subject = subject
	mail.HTML = content
	if len(cc) > 0 {
		mail.Cc = cc
	}
	select {
	case sendCh <- mail:
		return true
	case <-time.After(time.Second * 3):
		return false
	}

}