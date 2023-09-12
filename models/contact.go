package models

import (
	"chat-server/models/common"
	"chat-server/utils"
	"fmt"
)

type Contact struct {
	common.GlobalModel
	OwnerID  uint // 用户id
	TargetID uint // 好友id
	Type     int  // 对应的类型 0 1 3
	Desc     string
}

func (c *Contact) TableName() string {
	return "contact"
}

// 添加好友

func AddFriend(userId uint, targetId uint) int {
	data := Contact{}
	if targetId != 0 {
		utils.DB.Where("owner_id = ? AND target_id = ?", userId, targetId).First(&data)
		if data.OwnerID == 0 {
			if data.TargetID == 0 {
				fmt.Println("不是好友")
				// 开启事务
				tx := utils.DB.Begin()

				// 事务一旦开始，不管什么异常最终都会Rollback
				defer func() {
					if r := recover(); r != nil {
						tx.Rollback()
					}
				}()
				var friend = []Contact{
					{
						OwnerID:  userId,
						TargetID: targetId,
						Type:     1,
					},
					{
						OwnerID:  targetId,
						TargetID: userId,
						Type:     1,
					},
				}
				if err := tx.Create(&friend).Error; err != nil {
					fmt.Println(err)
					tx.Rollback()
					return 0
				}

				return 1
			}
		}
		return 0
	}
	return 0
}

// 好友列表

func FriendList(userId uint) []UserBasic {
	contacts := make([]Contact, 0)
	objIds := make([]uint64, 0)
	utils.DB.Where("owner_id = ? AND type = 1", userId).Find(&contacts)
	for _, v := range contacts {
		fmt.Println(">>>>>>>>>>>>>>", v)
		objIds = append(objIds, uint64(v.TargetID))
	}
	users := make([]UserBasic, 0)
	utils.DB.Where("id IN ?", objIds).Find(&users)
	return users
}
