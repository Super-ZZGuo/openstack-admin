package dto

import (
	"go-admin/app/admin/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type SysFlavorGetPageReq struct {
	dto.Pagination `search:"-"`
	SysFlavorOrder
}

type SysFlavorOrder struct {
	FlavorId   string `form:"flavorIdOrder"  search:"type:order;column:flavor_id;table:sys_flavor"`
	FlavorName string `form:"flavorNameOrder"  search:"type:order;column:flavor_name;table:sys_flavor"`
	Disk       string `form:"diskOrder"  search:"type:order;column:disk;table:sys_flavor"`
	Vcpu       string `form:"vcpuOrder"  search:"type:order;column:vcpu;table:sys_flavor"`
	Ram        string `form:"ramOrder"  search:"type:order;column:ram;table:sys_flavor"`
	Tag        string `form:"tagOrder"  search:"type:order;column:tag;table:sys_flavor"`
	CreateBy   string `form:"createByOrder"  search:"type:order;column:create_by;table:sys_flavor"`
	UpdateBy   string `form:"updateByOrder"  search:"type:order;column:update_by;table:sys_flavor"`
	CreatedAt  string `form:"createdAtOrder"  search:"type:order;column:created_at;table:sys_flavor"`
	UpdatedAt  string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:sys_flavor"`
	DeletedAt  string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:sys_flavor"`
}

func (m *SysFlavorGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type SysFlavorInsertReq struct {
	FlavorId   int    `json:"flavorId" comment:""` //
	FlavorName string `json:"flavorName" comment:""`
	Disk       int    `json:"disk" comment:""`
	Vcpu       int    `json:"vcpu" comment:""`
	Ram        int    `json:"ram" comment:""`
	Tag        string `json:"tag" comment:""`
	common.ControlBy
}

func (s *SysFlavorInsertReq) Generate(model *models.SysFlavor) {
	if s.FlavorId == 0 {
		model.FlavorId = s.FlavorId
	}
	model.FlavorName = s.FlavorName
	model.Disk = s.Disk
	model.Vcpu = s.Vcpu
	model.Ram = s.Ram
	model.Tag = s.Tag
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
}

func (s *SysFlavorInsertReq) GetId() interface{} {
	return s.FlavorId
}

// type SysFlavorUpdateReq struct {
// 	FlavorId   int    `uri:"flavorId" comment:""` //
// 	FlavorName string `json:"flavorName" comment:""`
// 	Disk       string `json:"disk" comment:""`
// 	Vcpu       string `json:"vcpu" comment:""`
// 	Ram        string `json:"ram" comment:""`
// 	Tag        string `json:"tag" comment:""`
// 	common.ControlBy
// }

// func (s *SysFlavorUpdateReq) Generate(model *models.SysFlavor) {
// 	if s.FlavorId == 0 {
// 		model.FlavorId = s.FlavorId
// 	}
// 	model.FlavorName = s.FlavorName
// 	model.Disk = s.Disk
// 	model.Vcpu = s.Vcpu
// 	model.Ram = s.Ram
// 	model.Tag = s.Tag
// 	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
// }

// func (s *SysFlavorUpdateReq) GetId() interface{} {
// 	return s.FlavorId
// }

// SysFlavorGetReq 功能获取请求参数
type SysFlavorGetReq struct {
	FlavorId int `uri:"id"`
}

func (s *SysFlavorGetReq) GetId() interface{} {
	return s.FlavorId
}

// SysFlavorDeleteReq 功能删除请求参数
type SysFlavorDeleteReq struct {
	Ids []int `json:"ids"`
}

func (s *SysFlavorDeleteReq) GetId() interface{} {
	return s.Ids
}
