package main

import (
	"fmt"
	"github.com/gin-gonic/contrib/static"
	"log"
	"net/http"
	"os"
	"time"

	"./service"
	"./structs"
	"./utils"
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

var identityKey = "id"

func helloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get(identityKey)
	c.JSON(200, gin.H{
		"userID":   claims["id"],
		"userName": user.(*User).UserName,
		"text":     "Hello World.",
	})
}

// User demo
type User struct {
	UserName  string
	FirstName string
	LastName  string
}

func main() {
	utils.StartPool()
	port := os.Getenv("PORT")
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	if port == "" {
		port = "8000"
	}
	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				UserName: claims["id"].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals structs.User
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}

			if utils.Login(loginVals) {
				userID := loginVals.UserName

				// 设置token
				return &User{
					UserName: userID,
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			//用于判断是否有权限
			if v, ok := data.(*User); ok && v.UserName == "2214839296@qq.com" {
				return true
			}
			//
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"auth":    false,
				"message": "未登录!",
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	r.POST("/login", authMiddleware.LoginHandler)

	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "ROUTE_NOT_FOUND", "msg": "该接口不存在!"})
	})

	t := new(service.Template)

	r.GET("/getTemplate", func(c *gin.Context) {
		page := utils.GetPage(c)
		c.JSON(200, t.ListByPage(page))
	})

	r.GET("/getDetailTemplate", func(c *gin.Context) {
		page := utils.GetPage(c)
		c.JSON(200, t.ListDetailByPage(page, string(c.Query("pid"))))
	})

	r.POST("/getNewest/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(200, t.GetNewestDetail(id))
	})

	r.POST("/checkNew", func(c *gin.Context) {
		var b structs.DetailData
		c.BindJSON(&b)
		c.JSON(200, t.CheckNew(b.Pid, b.Id))
	})

	r.Use(static.Serve("/download/", static.LocalFile("./template", true)))

	DownloadsPath := "template/"
	r.Use(authMiddleware.MiddlewareFunc())
	{

		r.POST("/auth", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"data": jwt.ExtractClaims(c)["id"],
			})
		})

		r.POST("/upload", func(c *gin.Context) {
			// Source
			file, err := c.FormFile("file")
			if err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
				return
			}
			fileId := utils.NewKeyId()
			userName := jwt.ExtractClaims(c)["id"].(string)
			filename := fileId + ".zip"
			if err := c.SaveUploadedFile(file, filename); err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
				return
			}
			oldLocation := filename
			newLocation := DownloadsPath + userName + "/" + filename
			os.MkdirAll(DownloadsPath+userName, os.ModePerm)
			e := os.Rename(oldLocation, newLocation)
			if e != nil {
				log.Fatal(err)
			}
			info := c.Request.FormValue("info")
			name := c.Request.FormValue("name")
			template := structs.Template{}
			template.Info = info
			template.Detail = info
			template.Id = utils.NewKeyId()
			template.Name = name
			template.User = userName
			id := t.Create(template)

			detail := structs.TemplateDetail{}
			detail.Id = fileId
			detail.Pid = id
			detail.Name = file.Filename
			detail.User = template.User
			detail.Info = info
			t.CreateDetail(detail)
			c.JSON(200, gin.H{
				"success":  true,
				"message":  "上传成功!",
				"serveId":  id,
				"detailId": fileId,
			})
		})

		user := r.Group("/user")
		{
			u := new(service.User)

			user.GET("/listByPage", func(c *gin.Context) {
				page := utils.GetPage(c)
				c.JSON(200, u.ListByPage(page))
			})

			user.POST("/update", func(c *gin.Context) {
				var user structs.User
				c.ShouldBindJSON(&user)
				c.JSON(200, u.Update(user))
			})

			user.POST("/create", func(c *gin.Context) {
				var data structs.User
				c.ShouldBindJSON(&data)
				if u.Create(data) {
					c.JSON(200, true)
				} else {
					c.JSON(500, utils.SendError("登录名重复,保存失败!"))
				}
			})

			user.POST("/remove", func(c *gin.Context) {
				var data structs.User
				c.ShouldBindJSON(&data)
				c.JSON(200, u.Remove(data))
			})
		}
	}

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}

}
