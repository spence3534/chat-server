package models

import "chat-server/models/common"

type GroupBasic struct {
	common.GlobalModel
	Name      string
	OwnerID   uint
	GroupHead string
	Type      int
	Desc      string
}

func (b GroupBasic) TableName() string {
	return "group_basic"
}
