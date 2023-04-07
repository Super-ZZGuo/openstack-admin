package dto

import (
	"go-admin/app/admin/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type SysProjectGetPageReq struct {
	dto.Pagination `search:"-"`
	ProjectId      string `form:"projectId"  search:"type:exact;column:project_id;table:sys_project"`
	ProjectName    string `form:"projectName"  search:"type:contains;column:project_name;table:sys_project"`
	Status         string `form:"status" search:"type:exact;column:status;table:sys_project" comment:"状态"`
	SysProjectOrder
}

type SysProjectOrder struct {
	ProjectOpenstackId string `form:"projectOpenstackIdOrder"  search:"type:order;column:project_openstack_id;table:sys_project"`
	Tag                string `form:"tagOrder"  search:"type:order;column:tag;table:sys_project"`
	CreateBy           string `form:"createByOrder"  search:"type:order;column:create_by;table:sys_project"`
	UpdateBy           string `form:"updateByOrder"  search:"type:order;column:update_by;table:sys_project"`
	CreatedAt          string `form:"createdAtOrder"  search:"type:order;column:created_at;table:sys_project"`
	UpdatedAt          string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:sys_project"`
	DeletedAt          string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:sys_project"`
}

func (m *SysProjectGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type SysProjectInsertReq struct {
	ProjectId          int    `json:"-" comment:""` //
	ProjectName        string `json:"projectName" comment:""`
	Tag                string `json:"tag" comment:""`
	Status             string `json:"-" comment:""`
	ProjectOpenstackId string `json:"-" comment:""`
	common.ControlBy
}

func (s *SysProjectInsertReq) Generate(model *models.SysProject) {
	if s.ProjectId == 0 {
		model.ProjectId = s.ProjectId
	}
	model.ProjectName = s.ProjectName
	model.Status = s.Status
	model.Tag = s.Tag
	model.ProjectOpenstackId = s.ProjectOpenstackId
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
}

func (s *SysProjectInsertReq) GetId() interface{} {
	return s.ProjectId
}

type SysProjectUpdateReq struct {
	ProjectId          int    `uri:"id" comment:""` //
	NewProjectName     string `json:"newProjectName" comment:""`
	Status             string `json:"status" comment:""`
	OldProjectName     string `json:"oldProjectName" comment:""`
	OldTag             string `json:"oldTag" comment:""`
	NewTag             string `json:"newTag" comment:""`
	Option             string `json:"option" comment:""`
	ProjectOpenstackId string `json:"-" comment:""`
	common.ControlBy
}

func (s *SysProjectUpdateReq) Generate(model *models.SysProject) {
	if s.ProjectId == 0 {
		model.ProjectId = s.ProjectId
	}
	model.ProjectName = s.NewProjectName
	if s.NewTag == "" && s.OldTag != "" {
		model.Tag = s.OldTag
	}
	if s.NewTag != "" {
		model.Tag = s.NewTag
	}
	model.Status = s.Status
	model.ProjectOpenstackId = s.ProjectOpenstackId
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
}

func (s *SysProjectUpdateReq) GetId() interface{} {
	return s.ProjectId
}

// SysProjectPutUpdate 更新数据库
type SysProjectPutUpdate struct {
	ProjectId          int
	ProjectName        string
	ProjectOpenstackId string
	Tag                string
	Status             string
}

func (s *SysProjectPutUpdate) GetId() interface{} {
	return s.ProjectId
}

// SysProjectGetReq 功能获取请求参数
type SysProjectGetReq struct {
	ProjectId int `uri:"id"`
}

func (s *SysProjectGetReq) GetId() interface{} {
	return s.ProjectId
}

// SysProjectDeleteReq 功能删除请求参数
type SysProjectDeleteReq struct {
	Ids          []int    `json:"ids"`
	ProjectNames []string `json:"projectNames" comment:""`
}

func (s *SysProjectDeleteReq) GetId() interface{} {
	return s.Ids
}
