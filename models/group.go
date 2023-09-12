package models

import (
	"chat-server/models/common"
	"chat-server/utils"
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
