package service

import (
	"chat-server/models"
	"chat-server/models/common/response"
	"chat-server/utils"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// GetUserList godoc
// @Tags 用户模块
// @Summary 查询所有用户信息
// @Success 200 {object} response.JsonResult
// @Router /getUserList [get]
func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()
	response.OkWithData(data, "请求成功", c)
}

// GetUserInfo godoc
// @Tags 用户模块
// @Summary 查询用户信息
// @Param id query string true "用户id"
// @Success 200 {object} response.JsonResult
// @Router /user/getUserInfo [get]
func GetUserInfo(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	data := models.GetUserInfo(user)
	response.OkWithData(data, "获取成功", c)
}

// CreateUser godoc
// @Tags 用户模块
// @Summary 注册接口
// @Produce application/json
// @Param name formData string true "用户名"
// @Param phone formData string true "手机号"
// @Param password formData string true "密码"
// @Param repassword formData string true "确认密码"
// @Success 200 {object} response.JsonResult
// @Router /createUser [post]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.PostForm("name")
	password := c.PostForm("password")
	repassword := c.PostForm("repassword")
	phone := c.PostForm("phone")

	salt := fmt.Sprintf("%06d", rand.Int31())

	data := models.FindUserByName(user.Name)
	userPhone := models.FindUserByPhone(phone)
	if data.Name != "" {
		response.FailWithMessage("用户名已注册", c)
		return
	} else if userPhone.Phone != "" {
		response.FailWithMessage("手机号已注册", c)
		return
	} else if password != repassword {
		response.FailWithMessage("两次密码不一致", c)
		return
	}

	user.Password = utils.MakePassword(password, salt)
	user.Salt = salt
	user.Phone = phone
	models.CreateUser(user)
	response.OkWithMessage("注册成功", c)
}

// DeleteUser godoc
// @Tags 用户模块
// @Summary 删除用户
// @Produce application/json
// @Param id query string true "用户id"
// @Success 200 {object} response.JsonResult
// @Router /user/deleteUser [get]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}

	userId, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(userId)
	models.DeleteUser(user)
	response.OkWithMessage("删除成功", c)
}

// UpdateUserInfo godoc
// @Tags 用户模块
// @Summary 更新用户信息
// @Produce application/json
// @Param id formData string true "id"
// @Param name formData string false "name"
// @Param password formData string false "password"
// @Param headPic formData string false "headPic"
// @Param email formData string false "email"
// @Param phone formData string false "phone"
// @Success 200 {object} response.JsonResult
// @Router /user/updateUserInfo [post]
func UpdateUserInfo(c *gin.Context) {
	user := models.UserBasic{}

	id, _ := strconv.Atoi(c.PostForm("id"))
	name := c.PostForm("name")
	password := c.PostForm("password")
	email := c.PostForm("email")
	headPic := c.PostForm("headPic")
	phone := c.PostForm("phone")
	fmt.Println(headPic)
	user.ID = uint(id)
	user.Name = name
	user.Password = password
	user.Email = email
	user.Phone = phone
	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{
			"message": "修改参数不匹配",
			"code":    400,
		})
	} else {
		models.UpdateUser(user)
		c.JSON(200, gin.H{
			"message": "更新成功",
			"code":    200,
			"data":    true,
		})
	}

}

// Login godoc
// @Tags 用户模块
// @Summary 登录接口
// @Produce application/json
// @Param name formData string false "用户名"
// @Param password formData string false "密码"
// @Success 200 {object} response.JsonResult
// @Router /login [post]
func Login(c *gin.Context) {
	data := models.UserBasic{}
	name := c.PostForm("name")
	password := c.PostForm("password")

	user := models.FindUserByName(name)

	if user.Name == "" {
		response.FailWithMessage("用户不存在", c)
		return
	}

	// 解密: 用数据库中存的里面和传过来的密码进行加密，再判断是否相等
	flag := utils.ValidPassword(password, user.Salt, user.Password)

	if !flag {
		response.FailWithMessage("登录失败", c)
		return
	}

	// 密码加密
	pwd := utils.MakePassword(password, user.Salt)
	// 查询用户信息
	data = models.FindUserByNameAndPwd(name, pwd)

	token := utils.GetToken(name)
	if token == "" {
		response.FailWithMessage("内部系统错误", c)
	} else {
		data.Token = token
		response.OkWithData(data, "登录成功", c)
	}
}

// 防止跨域站点请求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMessage(c *gin.Context) {
	conn, err := upGrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer func(conn *websocket.Conn) {
		err = conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(conn)

	MsgHandler(conn, c)
}

func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	msg, err := utils.Subscribe(c, utils.PublishKey)
	if err != nil {
		fmt.Println(err)
	}

	tm := time.Now().Format("2006-01-02 15:04:05")
	m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
	err = ws.WriteMessage(1, []byte(m))
	if err != nil {
		fmt.Println(err)
	}
}

// SendUserMsg godoc
// @Tags 用户模块
// @Summary websocket发送消息
// @Produce application/json
// @Param userId formData string false "userId"
// @Success 200 {object} response.JsonResult
// @Router /user/sendUserMsg [get]
func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}

// GetFriendList godoc
// @Tags 用户模块
// @Summary 查询好友列表
// @Param userId query string true "用户id"
// @Success 200 {object} response.JsonResult
// @Router /user/getFriendList [get]
func GetFriendList(c *gin.Context) {
	id := c.Query("userId")
	userId, _ := strconv.Atoi(id)
	users := models.FriendList(uint(userId))
	response.OkWithData(users, "查询成功", c)
}

// AddFriend godoc
// @Tags 用户模块
// @Summary 添加好友
// @Param userId query string true "用户id"
// @Param targetId query string true "目标id"
// @Success 200 {object} response.JsonResult
// @Router /user/AddFriend [post]
func AddFriend(c *gin.Context) {
	id := c.PostForm("userId")
	targetId := c.PostForm("targetId")
	uId, _ := strconv.Atoi(id)
	uTargetId, _ := strconv.Atoi(targetId)
	code := models.AddFriend(uint(uId), uint(uTargetId))
	if code == 1 {
		response.OkWithMessage("添加成功", c)
	} else {
		response.FailWithMessage("添加好友失败", c)
	}
}

// SearchFriend godoc
// @Tags 用户模块
// @Summary 搜索好友
// @Param phone query string true "手机号"
// @Success 200 {object} response.JsonResult
// @Router /user/getFriendList [get]
func SearchFriend(c *gin.Context) {
	phone := c.Query("phone")
	user := models.FindUserByPhone(phone)
	if user.Phone != "" {
		response.OkWithData(user, "查询成功", c)
	}
}

// CreateGroup godoc
// @Tags 用户模块
// @Summary 创建群
// @Param userId query string true "用户id"
// @Param name query string true "群名"
// @Success 200 {object} response.JsonResult
// @Router /user/createGroup [post]
func CreateGroup(c *gin.Context) {
	userId := c.PostForm("userId")
	name := c.PostForm("name")

	uid, _ := strconv.Atoi(userId)

	group := models.Group{}
	group.GroupName = name
	group.OwnerId = uint(uid)
	group.GroupLeaderId = uint(uid)

	code, msg := models.CreateGroup(group)

	if code != -1 {
		response.OkWithMessage(msg, c)
	} else {
		response.FailWithMessage(msg, c)
	}

}

// GetGroup godoc
// @Tags 用户模块
// @Summary 获取群列表
// @Param id query string true "用户id"
// @Success 200 {object} response.JsonResult
// @Router /user/getGroup [get]
func GetGroup(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("userId"))
	data, msg := models.GetGroup(uint(id))
	if len(data) == 0 {
		response.OkWithData(data, "查询失败", c)
	} else {
		response.OkWithData(data, msg, c)
	}
}

func JoinGroup(c *gin.Context) {
	id, _ := strconv.Atoi(c.PostForm("userId"))
	groupNameOrId := c.PostForm("groupNameOrId")
	code, msg := models.JoinGroup(groupNameOrId, uint(id))
	if code == 0 {
		response.OkWithMessage(msg, c)
	} else {
		response.FailWithMessage(msg, c)
	}
}
