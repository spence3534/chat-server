package models

import (
	"chat-server/models/common"
	"chat-server/utils"
	"fmt"
)

type Group struct {
	common.GlobalModel
	OwnerId       uint   `json:"ownerId"`
	Image         string `json:"image"`
	Desc          string `json:"desc"`
	GroupName     string `json:"groupName"`
	GroupLeaderId uint   `json:"groupLeader"`
}

func CreateGroup(group Group) (int, string) {
	user := UserBasic{}
	utils.DB.Where("id = ?", group.GroupLeaderId).First(&user)
	if user.ID == 0 {
		return -1, "用户不存在"
	}
	if group.GroupLeaderId == 0 {
		return -1, "创建群的用户id为空"
	}

	if group.GroupName == "" {
		return -1, "群名称不能为空"
	}

	utils.DB.Create(&group)
	return 0, "创建成功"
}

func GetGroup(id uint) ([]Group, string) {
	group := []Group{}
	utils.DB.Where("owner_id = ?", id).Find(&group)
	return group, "查询成功"
}

func JoinGroup(groupNameAndId interface{}, userId uint) (int, string) {
	contact := Contact{}
	group := Group{}
	contact.OwnerID = userId
	contact.Type = 2
	utils.DB.Where("id = ? or group_name = ?", groupNameAndId, groupNameAndId).First(&group)
	if group.GroupName == "" {
		return -1, "该群不存在"
	}
	contact.TargetID = group.ID
	utils.DB.Where("owner_id = ? and target_id = ? and type = 2", userId, group.ID).First(&contact)
	if contact.CreatedAt.IsZero() {
		utils.DB.Create(&contact)
		return 0, "加群成功"
	} else {
		fmt.Println(contact)
		return -1, "不能重复加群"
	}
}
