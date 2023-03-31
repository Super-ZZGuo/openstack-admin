package models

import (
	"go-admin/common/models"
)

type UseRange struct {
	RangeId      int    `json:"rangeId" gorm:"primaryKey;autoIncrement;comment:rangeid"`
	RangeName    string `json:"rangeName" gorm:"type:varchar(255);comment:RangeName"`
	Status       string `json:"status" gorm:"type:varchar(10);comment:Status"`
	ProjectName  string `json:"projectName" gorm:"type:varchar(100);comment:ProjectName"`
	RangeConsole string `json:"rangeConsole" gorm:"type:varchar(100);comment:"`
	Dept         string `json:"dept" gorm:"type:varchar(100);comment:"`
	Ipadress     string `json:"ipadress" gorm:"type:varchar(255);comment:"`
	models.ModelTime
	models.ControlBy
}

func (UseRange) TableName() string {
	return "sys_range"
}

func (e *UseRange) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *UseRange) GetId() interface{} {
	return e.RangeId
}
