package models

import "gorm.io/gorm"

// 人之间关系
type Contact struct {
	gorm.Model
	OwnerId  int64 //所有者id
	TargetId uint  //联系人id
	Type     int   //对应的类型 0  1  3

}

//func (table *Message) TableName() string {
//	return "contact"
//}
