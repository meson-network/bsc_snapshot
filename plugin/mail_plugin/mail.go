package mail_plugin

import (
	"fmt"
	"net/smtp"
	"net/textproto"
	"strconv"
	"time"

	"github.com/jordan-wright/email"
	"github.com/meson-network/bsc-data-file-utils/basic"
)

type Config struct {
	Host     string
	Port     int
	UserName string
	Password string
}

type Sender struct {
	host     string
	port     int
	userName string
	password string
	pool     *email.Pool
}

var instanceMap = map[string]*Sender{}

func GetInstance() *Sender {
	return GetInstance_("default")
}

func GetInstance_(name string) *Sender {
	sender := instanceMap[name]
	if sender == nil {
		basic.Logger.Errorln(name + "email plugin null")
	}
	return sender
}

func Init(config *Config) error {
	return Init_("default", config)
}

func Init_(name string, config *Config) error {

	if name == "" {
		name = "default"
	}

	_, exist := instanceMap[name]
	if exist {
		return fmt.Errorf("email instance <%s> has already been initialized", name)
	}

	if config.Port == 0 {
		config.Port = 587
	}

	sender := &Sender{
		host:     config.Host,
		port:     config.Port,
		userName: config.UserName,
		password: config.Password,
	}

	email_pool, pool_err := email.NewPool(
		sender.host+":"+strconv.Itoa(sender.port),
		4,
		smtp.PlainAuth("", sender.userName, sender.password, sender.host),
	)

	if pool_err != nil {
		return pool_err
	}

	sender.pool = email_pool
	instanceMap[name] = sender
	return nil
}

func (s *Sender) Send(from_text string, to_address string, subject string, body string) error {

	e := &email.Email{
		To:      []string{to_address},
		From:    from_text,
		Subject: subject,
		Text:    []byte(body),
		Headers: textproto.MIMEHeader{},
	}

	var err error
	for i := 0; i < 3; i++ {
		err = s.pool.Send(e, 8*time.Second)
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	return err
}
