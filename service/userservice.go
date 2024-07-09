package service

import (
	"fmt"
	valid "github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/http"
	"project2/models"
	"project2/utils"
	"strconv"
	"time"
)

// GetUserList
// @Summary 所有用户
// @Tags 首页
// @Success 200 {string} json{"code", "message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()

	c.JSON(200, gin.H{
		"code":    -1, //0为成功，-1为失败
		"message": "获取用户列表成功",
		"data":    data,
	})
}

// CreateUser
// @Summary 添加用户
// @Tags 用户注册
// @param name query string true "用户名"
// @param password query string true "密码"
// @param repassword query string false "确认密码"
// @Success 200 {string} json{"code", "message"}
// @Router /register [get]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.Query("name")
	password := c.Query("password")
	repassword := c.Query("repassword")

	//判存
	data := models.FindUserByName(user.Name) //该函数返回一个UserBasic对象
	if data.Name != "" {
		c.JSON(200, gin.H{
			"code":    -1, //0为成功，-1为失败
			"message": "用户已注册",
			"data":    data,
		})
		return
	}
	if repassword != password {
		c.JSON(200, gin.H{
			"code":    -1, //0为成功，-1为失败
			"message": "两次密码不一致",
			"data":    data,
		})
		return
	}

	//加密
	salt := fmt.Sprintf("%06d", rand.Int31())
	//user.Password = password
	user.Password = utils.MakePassword(password, salt) //加密操作
	user.Salt = salt
	models.CreateUser(user)
	c.JSON(200, gin.H{
		"code":    0, //0为成功，-1为失败
		"message": "注册成功",
		"data":    user,
	})
}

// DeleteUser
// @Summary 删除用户
// @Tags 注销
// @param id query string false "用户id"
// @Success 200 {string} json{"code", "message"}
// @Router /logout [get]
func DeleteUser(c *gin.Context) {
	//user := models.UserBasic{}
	//user.Name = c.Query("name")
	//data := make([]*models.UserBasic, 10)
	//data = models.GetUserList()
	//for i := 0; i < len(data); i++ {
	//	if data[i].Name == user.Name {
	//		models.DeleteUser(user)
	//	}
	//}
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	models.DeleteUser(user)
	fmt.Println("----delete---:", user)
	c.JSON(200, gin.H{
		"code":    0, //0为成功，-1为失败
		"message": "注销成功",
		"data":    user,
	})
}

// UpdateUser
// @Summary 更新信息
// @Tags 更新用户信息
// @param id formData string false "id"
// @param name formData string false "name"
// @param password formData string false "password"
// @param phone formData string false "phone"
// @param identity formData string false "identity"
// @param email formData string false "email"
// @param clientIp formData string false "clientIp"
// @Success 200 {string} json{"code", "message"}
// @Router /updateInfo [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.Password = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Identity = c.PostForm("indetity")
	user.Email = c.PostForm("email")
	user.ClentIp = c.PostForm("clientIp")

	_, err := valid.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"message": "信息格式不匹配",
		})
	} else {
		res := models.UpdateUser(user)
		fmt.Println("-----res------:", res)
		c.JSON(200, gin.H{
			"message": "更新成功",
		})
	}

}

// FindUserByNameAndPwd
// @Summary 登录
// @Tags 登录
// @param name query string true "用户名"
// @param password query string true "密码"
// @Success 200 {string} json{"code", "message"}
// @Router /findUserByNameAndPwd [post]
func FindUserByNameAndPwd(c *gin.Context) {
	data := models.UserBasic{}
	name := c.Query("name")
	password := c.Query("password")

	user := models.FindUserByName(name)
	if user.Name == "" {
		c.JSON(200, gin.H{
			"code":    -1, //0为成功，-1为失败
			"message": "用户不存在",
			"data":    user,
		})
		return
	}
	//解密
	fmt.Println("user:::", user)
	flag := utils.ValidMakePassword(password, user.Salt, user.Password)
	if !flag {
		c.JSON(200, gin.H{
			"code":    -1, //0为成功，-1为失败
			"message": "密码不正确",
			"data":    flag,
		})
		return
	}
	//加密打印
	pwd := utils.MakePassword(password, user.Salt)
	data = models.FindUserByNameAndPwd(name, pwd)

	c.JSON(200, gin.H{
		"code":    0, //0为成功，-1为失败
		"message": "登录成功",
		"data":    data,
	})

}

// 防止跨域站点伪造请求
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 信息发送接收具体业务逻辑实现
func SendMsg(c *gin.Context) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(ws *websocket.Conn) { //为甚
		err := ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	MsgHandler(ws, c)
}
func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	msg, err := utils.Subscribe(c, utils.PublishKey)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("发送消息：", msg)
	tm := time.Now().Format("2006-01-02 15:04:05")
	m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
	err = ws.WriteMessage(1, []byte(m))
	if err != nil {
		fmt.Println(err)
	}
}
func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}
