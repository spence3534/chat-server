package models

import "chat-server/models/common"

type GroupBasic struct {
	common.GlobalModel
	Name      string
	OwnerID   uint
	GroupHead string `json:"groupHead"`
	Type      int    `json:"type"`
	Desc      string `json:"desc"`
}

func (b GroupBasic) TableName() string {
	return "group_basic"
}
