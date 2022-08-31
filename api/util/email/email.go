package email

import (
	"bytes"
	"crypto/tls"
	"html/template"
	"net"
	"net/smtp"

	"github.com/jordan-wright/email"
	"github.com/spf13/viper"

	"zendea/util/log"
)

var emailTemplate = `
<div style="background-color:white;border-top:2px solid #12ADDB;box-shadow:0 1px 3px #AAAAAA;line-height:180%;padding:0 15px 12px;width:500px;margin:50px auto;color:#555555;font-family:'Century Gothic','Trebuchet MS','Hiragino Sans GB',微软雅黑,'Microsoft Yahei',Tahoma,Helvetica,Arial,'SimSun',sans-serif;font-size:12px;">
    <h2 style="border-bottom:1px solid #DDD;font-size:14px;font-weight:normal;padding:13px 0 10px 8px;">
        <span style="color: #12ADDB;font-weight:bold;">
            {{.Title}}
        </span>
    </h2>
    <div style="padding:0 12px 0 12px; margin-top:18px;">
        {{if .Content}}
		<p>
            {{.Content}}
        </p>
		{{end}}
		{{if .QuoteContent}}
		<div style="background-color: #f5f5f5;padding: 10px 15px;margin:18px 0;word-wrap:break-word;">
            {{.QuoteContent}}
        </div>
		{{end}}
       
		{{if .Url}}
        <p>
            <a style="text-decoration:none; color:#12addb" href="{{.Url}}" target="_blank" rel="noopener">点击查看详情</a>
        </p>
		{{end}}
    </div>
</div>
`

// 发送模版邮件
func SendTemplateEmail(to, subject, title, content, quoteContent, url string) {
	tpl, err := template.New("emailTemplate").Parse(emailTemplate)
	if err != nil {
		log.Error(err.Error())
		return
	}
	var b bytes.Buffer
	err = tpl.Execute(&b, map[string]interface{}{
		"Title":        title,
		"Content":      content,
		"QuoteContent": quoteContent,
		"Url":          url,
	})
	if err != nil {
		log.Error(err.Error())
		return
	}
	html := b.String()
	SendEmail(to, subject, html)
}

// 发送邮件
func SendEmail(to string, subject, html string) {
	var (
		host      = viper.GetString("smtp.host")
		port      = viper.GetString("smtp.port")
		username  = viper.GetString("smtp.username")
		password  = viper.GetString("smtp.password")
		ssl       = viper.GetBool("smtp.ssl")
		addr      = net.JoinHostPort(host, port)
		auth      = smtp.PlainAuth("", username, password, host)
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         host,
		}
	)

	e := email.NewEmail()
	e.From = viper.GetString("smtp.username")
	e.To = []string{to}
	e.Subject = subject
	e.HTML = []byte(html)

	if ssl {
		if err := e.SendWithTLS(addr, auth, tlsConfig); err != nil {
			log.Error("发送邮件异常")
		}
	} else {
		if err := e.Send(addr, auth); err != nil {
			log.Error("发送邮件异常")
		}
	}
}
