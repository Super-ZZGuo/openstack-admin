package models

import (
	"go-admin/common/models"
)

type SysProject struct {
	ProjectId          int    `json:"projectId" gorm:"primaryKey;autoIncrement;comment:projectId"`
	ProjectName        string `json:"projectName" gorm:"type:varchar(10);comment:ProjectName"`
	ProjectOpenstackId string `json:"projectOpenstackId" gorm:"type:varchar(100);comment:ProjectOpenstackId"`
	Tag                string `json:"tag" gorm:"type:varchar(100);comment:tag"`
	models.ModelTime
	models.ControlBy
}

func (SysProject) TableName() string {
	return "sys_project"
}

func (e *SysProject) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysProject) GetId() interface{} {
	return e.ProjectId
}
