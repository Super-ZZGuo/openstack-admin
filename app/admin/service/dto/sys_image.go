package dto

import (
	"go-admin/app/admin/models"
	"go-admin/common/dto"
	common "go-admin/common/models"
)

type SysImageGetPageReq struct {
	dto.Pagination `search:"-"`
	ImageName      string `form:"imageName"  search:"type:contains;column:image_name;table:sys_image"`
	Type           string `form:"type"  search:"type:contains;column:type;table:sys_image"`
	SysImageOrder
}

type SysImageOrder struct {
	ImageId     string `form:"imageIdOrder"  search:"type:order;column:image_id;table:sys_image"`
	OpenstackId string `form:"openstackIdOrder"  search:"type:order;column:openstack_id;table:sys_image"`
	Tag         string `form:"tagOrder"  search:"type:order;column:tag;table:sys_image"`
	CreateBy    string `form:"createByOrder"  search:"type:order;column:create_by;table:sys_image"`
	UpdateBy    string `form:"updateByOrder"  search:"type:order;column:update_by;table:sys_image"`
	CreatedAt   string `form:"createdAtOrder"  search:"type:order;column:created_at;table:sys_image"`
	UpdatedAt   string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:sys_image"`
	DeletedAt   string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:sys_image"`
}

func (m *SysImageGetPageReq) GetNeedSearch() interface{} {
	return *m
}

type SysImageInsertReq struct {
	ImageId   int    `json:"-" comment:""` //
	ImageName string `json:"imageName" comment:""`
	Tag       string `json:"tag" comment:""`
	common.ControlBy
}

func (s *SysImageInsertReq) Generate(model *models.SysImage) {
	if s.ImageId == 0 {
		model.ImageId = s.ImageId
	}
	model.ImageName = s.ImageName
	model.Tag = s.Tag
	model.CreateBy = s.CreateBy // 添加这而，需要记录是被谁创建的
}

func (s *SysImageInsertReq) GetId() interface{} {
	return s.ImageId
}

type SysImageGetUpdateReq struct {
	ImageId      int    `uri:"id" comment:""` //
	ImageNewName string `json:"imageNewName" comment:""`
	ImageOldName string `json:"imageOldName" comment:""`
	Newtag       string `json:"newTag" comment:""`
	OldTag       string `json:"oldTag" comment:""`
	common.ControlBy
}

func (s *SysImageGetUpdateReq) Generate(model *models.SysImage) {
	if s.ImageId == 0 {
		model.ImageId = s.ImageId
	}
	if s.ImageNewName == "" {
		model.ImageName = s.ImageOldName
	} else {
		model.ImageName = s.ImageNewName
	}
	if s.Newtag == "" {
		model.ImageName = s.OldTag
	} else {
		model.ImageName = s.Newtag
	}
	model.UpdateBy = s.UpdateBy // 添加这而，需要记录是被谁更新的
}

func (s *SysImageGetUpdateReq) GetId() interface{} {
	return s.ImageId
}

type SysImagePutUpdateReq struct {
	ImageId   int
	ImageName string
	Tag       string
	common.ControlBy
}

func (s *SysImagePutUpdateReq) GetId() interface{} {
	return s.ImageId
}

type UploadSysImage struct {
	ImageId      int    `form:"imageId" comment:"用户ID" vd:"len($)>0"` // 用户ID
	ImageName    string `form:"imageName" comment:"头像" vd:"len($)>0"`
	Type         string `form:"type" comment:"类型" vd:"len($)>0"`
	Path         string `json:"-" comment:"存储位置" vd:"len($)>0"`
	Openstack_id string `json:"-"`
	common.ControlBy
}

func (s *UploadSysImage) GetId() interface{} {
	return s.ImageId
}

func (s *UploadSysImage) Generate(model *models.SysImage) {
	if s.ImageId != 0 {
		model.ImageId = s.ImageId
	}
	model.ImageName = s.ImageName

}

// SysImageGetReq 功能获取请求参数
type SysImageGetReq struct {
	ImageId int `uri:"id"`
}

func (s *SysImageGetReq) GetId() interface{} {
	return s.ImageId
}

// SysImageDeleteReq 功能删除请求参数
type SysImageDeleteReq struct {
	Ids      []int    `json:"ids"`
	Names    []string `json:"names"`
	IsUpload bool     `json:"isUpload"`
}

func (s *SysImageDeleteReq) GetId() interface{} {
	return s.Ids
}

type SysImageDelete struct {
	Ids []int `json:"ids"`
}

func (s *SysImageDelete) GetId() interface{} {
	return s.Ids
}
