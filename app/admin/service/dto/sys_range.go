package dto

import (
	"go-admin/app/admin/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type SysRangeGetPageReq struct {
	dto.Pagination `search:"-"`
	ProjectName    string `form:"projectName" search:"type:contains;column:project_name;table:sys_range" comment:""`
	SysRangeOrder
}

type SysRangeOrder struct {
	RangeId          string `form:"rangeIdOrder"  search:"type:order;column:range_id;table:sys_range"`
	RangeName        string `form:"rangeNameOrder"  search:"type:order;column:range_name;table:sys_range"`
	Status           string `form:"statusOrder"  search:"type:order;column:status;table:sys_range"`
	Image            string `form:"imageOrder"  search:"type:order;column:image;table:sys_range"`
	Flavor           string `form:"flavorOrder"  search:"type:order;column:flavor;table:sys_range"`
	RangeOpenstackId string `form:"rangeOpenstackIdOrder"  search:"type:order;column:range_openstack_id;table:sys_range"`
	CreateBy         string `form:"createByOrder"  search:"type:order;column:create_by;table:sys_range"`
	UpdateBy         string `form:"updateByOrder"  search:"type:order;column:update_by;table:sys_range"`
	CreatedAt        string `form:"createdAtOrder"  search:"type:order;column:created_at;table:sys_range"`
	UpdatedAt        string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:sys_range"`
	DeletedAt        string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:sys_range"`
}

type ProjectJoin struct {
	ProjectName string `search:"type:contains;column:project_name;table:sys_project" form:"projectName"`
}

func (m *SysRangeGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type SysRangeNetwork struct {
	Ipadress    string `json:"ipadress" comment:""`
	NetworkName string `json:"networkName" comment:""`
}

type SysRangeInsertReq struct {
	RangeId          int               `json:"-" comment:""` //
	RangeName        string            `json:"rangeName" comment:""`
	Status           string            `json:"status" comment:""`
	Image            string            `json:"image" comment:""`
	Flavor           string            `json:"flavor" comment:""`
	ProjectId        int               `json:"projectId" comment:""`
	ProjectName      string            `json:"projectName"`
	RangeOpenstackId string            `json:"-"`
	Network          []SysRangeNetwork `json:"network"`
	common.ControlBy
}

func (s *SysRangeInsertReq) Generate(model *models.SysRange) {
	if s.RangeId != 0 {
		model.RangeId = s.RangeId
	}
	model.RangeName = s.RangeName
	model.Status = s.Status
	model.Image = s.Image
	model.Flavor = s.Flavor
	model.RangeOpenstackId = s.RangeOpenstackId
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.ProjectId = s.ProjectId
	model.ProjectName = s.ProjectName
}

func (s *SysRangeInsertReq) GetId() interface{} {
	return s.RangeId
}

type SysRangeUpdateReq struct {
	RangeId     int    `uri:"id" comment:""`
	RangeName   string `json:"rangeName" comment:""`
	Image       string `json:"image" comment:""`
	ProjectName string `json:"projectName"`
	Option      string `json:"option"`
	common.ControlBy
}

func (s *SysRangeUpdateReq) Generate(model *models.SysRange) {
	if s.RangeId == 0 {
		model.RangeId = s.RangeId
	}
	model.RangeName = s.RangeName
	model.Image = s.Image
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
	model.ProjectName = s.ProjectName
}

func (s *SysRangeUpdateReq) GetId() interface{} {
	return s.RangeId
}

type SysRangeStatusUpadteReq struct {
	RangeId int
	Status  string
	common.ControlBy
}

func (s *SysRangeStatusUpadteReq) GetId() interface{} {
	return s.RangeId
}

// SysRangeGetReq 功能获取请求参数
type SysRangeGetReq struct {
	RangeId     int    `uri:"id"`
	ProjectName string `json:"projectName"`
}

func (s *SysRangeGetReq) GetId() interface{} {
	return s.RangeId
}

// SysRangeDeleteReq 功能删除请求参数
type SysRangeDeleteReq struct {
	Ids         []int    `json:"ids"`
	ProjectName string   `json:"projectName"`
	RangeNames  []string `json:"rangeNames" comment:""`
}

func (s *SysRangeDeleteReq) GetId() interface{} {
	return s.Ids
}
