package router

import (
	"../structs"
	"../utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func CreateEmailCode(c *gin.Context) {
	var code structs.RegisterCache
	c.ShouldBind(&code)
	engine := utils.GetCon()
	engine.Get(&code)
	if len(code.Code) > 0 {
		utils.ReFail(c, "该邮箱已注册!")
	} else {
		code.Code = utils.RandomString(6, "A")
		engine.Insert(code)
		utils.SendMail(code.Id, "GenaretorTool注册码", fmt.Sprintf("您的验证码:<h3>%s</h3>", code.Code))
		utils.ReSuccess(c, "获取验证码成功!")
	}
}

func Register(c *gin.Context) {
	var user structs.User
	c.ShouldBind(&user)
	var code structs.RegisterCache
	code.Id = user.Username
	engine := utils.GetCon()
	engine.Get(&code)
	if user.Code == code.Code && code.Created.Add(60 * 60 * 30).Before(time.Now()) {
		user.Id = user.Username
		engine.Insert(user)
		utils.ReSuccess(c, "注册成功!")
	} else {
		utils.ReFail(c, "验证码错误!")
	}
}

func Login(c *gin.Context) bool {
	var user structs.User
	c.ShouldBind(&user)
	engine := utils.GetCon()
	has, _ :=engine.Cols("password").Get(&user)
	return has
}