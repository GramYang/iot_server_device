package router

import (
	"github.com/gin-gonic/gin"

	g "github.com/GramYang/gylog"

	sc "iot_server8/sqlx_client"

	"iot_server8/util"
)

type loginInfomation struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func loginRouter(r *gin.Engine, base string) {
	if base != "" {
		rg := r.Group("/" + base)
		rg.POST("/login", login)
	} else {
		r.POST("/login", login)
	}
}

func login(c *gin.Context) {
	var li loginInfomation
	if err := c.ShouldBindJSON(&li); err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	g.Debugln("login:", li)
	if li.Username == "" || li.Password == "" {
		c.JSON(400, gin.H{
			"message": "username or password invalid",
		})
		return
	}
	//检查用户名是否注册
	existed, err := sc.IsUserExisted(li.Username)
	if err != nil {
		g.Errorln(err)
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	if !existed {
		c.JSON(400, gin.H{
			"message": "username isn't existed",
		})
		return
	}
	//用户名注册后再验证密码
	user, err := sc.GetUserInfo(li.Username)
	if err != nil {
		g.Errorln(err)
		c.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}
	if user.Password == li.Password {
		token := util.CreateToken()
		c.JSON(200, gin.H{
			"message": "login success",
			"token":   token,
		})
		g.Debugln("200", token)
	} else {
		c.JSON(400, gin.H{
			"message": "password wrong",
		})
	}
}
