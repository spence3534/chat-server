package models

import (
	"chat-server/models/common"
	"chat-server/utils"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type UserBasic struct {
	common.GlobalModel
	Name          string `json:"name"`
	Password      string `json:"password"`
	Phone         string `json:"phone" valid:"matches(^1[3-9]{1}\\d{9}$)"`
	ClientIp      string `json:"clientIp"`
	ClientPort    string `json:"clientPort"`
	Email         string `json:"email" valid:"email"`
	Salt          string
	LoginTime     *time.Time `json:"loginTime"`
	HeartbeatTime *time.Time `json:"heartbeatTime"`
	LoginOutTime  *time.Time `gorm:"column:login_out_time" json:"login_out_time"`
	IsLogout      bool       `json:"isLogout"`
	DeviceInfo    string     `json:"deviceInfo"`
	Token         string     `json:"token"`
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	utils.DB.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}

func GetUserInfo(user UserBasic) *UserBasic {
	u := user
	utils.DB.First(&u)
	return &u
}

func FindUserByNameAndPwd(name, password string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ? AND password = ?", name, password).First(&user)
	return user
}

func FindUserByName(name string) UserBasic {
	data := UserBasic{}
	utils.DB.Where("name = ?", name).First(&data)
	return data
}

func FindUserByPhone(phone string) UserBasic {
	data := UserBasic{}
	utils.DB.Where("phone = ?", phone).First(&data)
	return data
}

func FindUserByEmail(email string) UserBasic {
	data := UserBasic{}
	utils.DB.Where("email = ?", email).First(&data)
	return data
}

func CreateUser(user UserBasic) *gorm.DB {
	return utils.DB.Create(&user)
}

func DeleteUser(user UserBasic) *gorm.DB {
	return utils.DB.Delete(&user)
}

func UpdateUser(user UserBasic) *gorm.DB {
	return utils.DB.Model(&user).Where("id = ?", user.ID).Updates(UserBasic{
		Name:     user.Name,
		Password: user.Password,
		Phone:    user.Phone,
		Email:    user.Email,
	})
}
