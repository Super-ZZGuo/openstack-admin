package dto

import (
	"go-admin/app/admin/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type SysNetworkGetPageReq struct {
	dto.Pagination `search:"-"`
	ProjectName    string `form:"projectName"  search:"type:contains;column:project_name;table:sys_network"`
	NetworkName    string `form:"networkName"  search:"type:contains;column:network_name;table:sys_network"`
	SysNetworkOrder
}

type SysNetworkOrder struct {
	NetworkId string `form:"networkIdOrder"  search:"type:order;column:network_id;table:sys_network"`
	Cidr      string `form:"cidrOrder"  search:"type:order;column:cidr;table:sys_network"`
	PoolStart string `form:"poolStartOrder"  search:"type:order;column:pool_start;table:sys_network"`
	PoolEnd   string `form:"poolEndOrder"  search:"type:order;column:pool_end;table:sys_network"`
	Tag       string `form:"tagOrder"  search:"type:order;column:tag;table:sys_network"`
	CreateBy  string `form:"createByOrder"  search:"type:order;column:create_by;table:sys_network"`
	UpdateBy  string `form:"updateByOrder"  search:"type:order;column:update_by;table:sys_network"`
	CreatedAt string `form:"createdAtOrder"  search:"type:order;column:created_at;table:sys_network"`
	UpdatedAt string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:sys_network"`
	DeletedAt string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:sys_network"`
}

func (m *SysNetworkGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type SysNetworkInsertReq struct {
	NetworkId   int    `json:"-" comment:""` //
	NetworkName string `json:"networkName" comment:""`
	ProjectName string `json:"projectName" comment:""`
	Cidr        string `json:"cidr" comment:""`
	PoolStart   string `json:"poolStart" comment:""`
	PoolEnd     string `json:"poolEnd" comment:""`
	Tag         string `json:"tag" comment:""`
	common.ControlBy
}

func (s *SysNetworkInsertReq) Generate(model *models.SysNetwork) {
	if s.NetworkId == 0 {
		model.NetworkId = s.NetworkId
	}
	model.NetworkName = s.NetworkName
	model.Cidr = s.Cidr
	model.ProjectName = s.ProjectName
	model.PoolStart = s.PoolStart
	model.PoolEnd = s.PoolEnd
	model.Tag = s.Tag
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
}

func (s *SysNetworkInsertReq) GetId() interface{} {
	return s.NetworkId
}

type SysNetworkUpdateReq struct {
	NetworkId      int    `uri:"id" comment:""` //
	ProjectName    string `json:"projectName" comment:""`
	NetworkNewName string `json:"networkNewName" comment:""`
	NetworkOldName string `json:"networkOldName" comment:""`
	Newtag         string `json:"newTag" comment:""`
	OldTag         string `json:"oldTag" comment:""`
	common.ControlBy
}

func (s *SysNetworkUpdateReq) Generate(model *models.SysNetwork) {
	if s.NetworkId == 0 {
		model.NetworkId = s.NetworkId
	}
	if s.NetworkNewName == "" {
		model.NetworkName = s.NetworkOldName
	} else {
		model.NetworkName = s.NetworkNewName
	}
	if s.Newtag == "" {
		model.Tag = s.OldTag
	} else {
		model.Tag = s.Newtag
	}
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
}

func (s *SysNetworkUpdateReq) GetId() interface{} {
	return s.NetworkId
}

// SysNetworkGetReq 功能获取请求参数
type SysNetworkGetReq struct {
	NetworkId int `uri:"id"`
}

func (s *SysNetworkGetReq) GetId() interface{} {
	return s.NetworkId
}

type SysNetworkPutUpdateReq struct {
	NetworkId   int
	NetworkName string
	Tag         string
	common.ControlBy
}

func (s *SysNetworkPutUpdateReq) GetId() interface{} {
	return s.NetworkId
}

// SysNetworkDeleteReq 功能删除请求参数
type SysNetworkDeleteReq struct {
	ProjectName  string   `json:"projectName"`
	NetworkNames []string `json:"networkNames"`
	Ids          []int    `json:"ids"`
}

func (s *SysNetworkDeleteReq) GetId() interface{} {
	return s.Ids
}

type SysNetworkGetPageRespone struct {
	ProjectName string              `json:"projectName"`
	Children    []models.SysNetwork `json:"children"`
}
