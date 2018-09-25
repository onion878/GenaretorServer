package utils

import (
	"bytes"
	"crypto/tls"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"math/rand"
	"strings"
	"time"
)

//RandomString(8, "A") 大写
//RandomString(8, "a0") 小写
//RandomString(20, "Aa0") 混合
func RandomString(randLength int, randType string) (result string) {
	var num string = "0123456789"
	var lower string = "abcdefghijklmnopqrstuvwxyz"
	var upper string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := bytes.Buffer{}
	if strings.Contains(randType, "0") {
		b.WriteString(num)
	}
	if strings.Contains(randType, "a") {
		b.WriteString(lower)
	}
	if strings.Contains(randType, "A") {
		b.WriteString(upper)
	}
	var str = b.String()
	var strLen = len(str)
	if strLen == 0 {
		result = ""
		return
	}

	rand.Seed(time.Now().UnixNano())
	b = bytes.Buffer{}
	for i := 0; i < randLength; i++ {
		b.WriteByte(str[rand.Intn(strLen)])
	}
	result = b.String()
	return
}

func ReSuccess(c *gin.Context, data interface{})  {
	c.JSON(200, gin.H{
		"success": true,
		"time": time.Now().Format(time.RFC3339),
		"data": data,
	})
}

func ReFail(c *gin.Context, val interface{})  {
	c.JSON(200, gin.H{
		"success": false,
		"time": time.Now().Format(time.RFC3339),
		"val": val,
	})
}

func SendMail(email string, title string, content string) {
	d := gomail.NewDialer("smtp.qq.com", 587, "genaretor@qq.com", "nbvlluxakyzgebji")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	m := gomail.NewMessage()
	m.SetHeader("From", "genaretor@qq.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", title)
	m.SetBody("text/html", content)
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}