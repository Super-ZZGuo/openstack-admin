package dto

import (
	"go-admin/app/admin/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type SysRangeGetPageReq struct {
	dto.Pagination `search:"-"`
	SysRangeOrder
}

type SysRangeOrder struct {
	RangeId          string `form:"rangeIdOrder"  search:"type:order;column:range_id;table:sys_range"`
	TenantName       string `form:"tenantNameOrder"  search:"type:order;column:tenant_name;table:sys_range"`
	RangeName        string `form:"rangeNameOrder"  search:"type:order;column:range_name;table:sys_range"`
	Status           string `form:"statusOrder"  search:"type:order;column:status;table:sys_range"`
	Image            string `form:"imageOrder"  search:"type:order;column:image;table:sys_range"`
	Flavor           string `form:"flavorOrder"  search:"type:order;column:flavor;table:sys_range"`
	CreateBy         string `form:"createByOrder"  search:"type:order;column:create_by;table:sys_range"`
	UpdateBy         string `form:"updateByOrder"  search:"type:order;column:update_by;table:sys_range"`
	CreatedAt        string `form:"createdAtOrder"  search:"type:order;column:created_at;table:sys_range"`
	UpdatedAt        string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:sys_range"`
	DeletedAt        string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:sys_range"`
	RangeOpenstackID string `form:"rangeOpenstackIDOrder"  search:"type:order;column:range_openstack_ID;table:sys_range"`
}

func (m *SysRangeGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type SysRangeInsertReq struct {
	RangeId          int    `json:"rangeId" comment:""` //
	TenantName       string `json:"tenantName" comment:""`
	RangeName        string `json:"rangeName" comment:""`
	Status           string `json:"status" comment:""`
	Image            string `json:"image" comment:""`
	Flavor           string `json:"flavor" comment:""`
	RangeOpenstackID string `json:"rangeOpenstackID" comment:""`
	common.ControlBy
}

func (s *SysRangeInsertReq) Generate(model *models.SysRange) {
	if s.RangeId != 0 {
		model.RangeId = s.RangeId
	}
	model.TenantName = s.TenantName
	model.RangeName = s.RangeName
	model.Status = s.Status
	model.Image = s.Image
	model.Flavor = s.Flavor
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
	model.RangeOpenstackID = s.RangeOpenstackID
}

func (s *SysRangeInsertReq) GetId() interface{} {
	return s.RangeId
}

type SysRangeUpdateReq struct {
	RangeId          int    `uri:"rangeId" comment:""` //
	TenantName       string `json:"tenantName" comment:""`
	RangeName        string `json:"rangeName" comment:""`
	Status           string `json:"status" comment:""`
	Image            string `json:"image" comment:""`
	Flavor           string `json:"flavor" comment:""`
	RangeOpenstackID string `json:"rangeOpenstackID" comment:""`
	common.ControlBy
}

func (s *SysRangeUpdateReq) Generate(model *models.SysRange) {
	if s.RangeId == 0 {
		model.RangeId = s.RangeId
	}
	model.TenantName = s.TenantName
	model.RangeName = s.RangeName
	model.Status = s.Status
	model.Image = s.Image
	model.Flavor = s.Flavor
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
	model.RangeOpenstackID = s.RangeOpenstackID
}

func (s *SysRangeUpdateReq) GetId() interface{} {
	return s.RangeId
}

// SysRangeGetReq 功能获取请求参数
type SysRangeGetReq struct {
	RangeId int `uri:"rangeId"`
}

func (s *SysRangeGetReq) GetId() interface{} {
	return s.RangeId
}

// SysRangeDeleteReq 功能删除请求参数
type SysRangeDeleteReq struct {
	Ids              []int    `json:"ids"`
	RangeOpenstackID []string `json:"rangeOpenstackID" comment:""`
}

func (s *SysRangeDeleteReq) GetOpenstackId() interface{} {
	return s.RangeOpenstackID
}

func (s *SysRangeDeleteReq) GetId() interface{} {
	return s.Ids
}
