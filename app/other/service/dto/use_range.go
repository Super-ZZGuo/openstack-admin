package dto

import (
	"go-admin/common/dto"
)

type UseRangeGetPageReq struct {
	dto.Pagination `search:"-"`
	ProjectName    string `form:"projectName" search:"type:contains;column:project_name;table:sys_range" comment:""`
	Dept           string `form:"dept"  search:"type:contains;column:dept;table:sys_range"`
	UseRangeOrder
}

type UseRangeOrder struct {
	RangeId      string `form:"rangeIdOrder"  search:"type:order;column:range_id;table:sys_range"`
	RangeName    string `form:"rangeNameOrder"  search:"type:order;column:range_name;table:sys_range"`
	Status       string `form:"statusOrder"  search:"type:order;column:status;table:sys_range"`
	RangeConsole string `form:"rangeConsoleOrder"  search:"type:order;column:range_console;table:sys_range"`
	Ipadress     string `form:"ipadressOrder"  search:"type:order;column:ipadress;table:sys_range"`
	CreateBy     string `form:"createByOrder"  search:"type:order;column:create_by;table:sys_range"`
	UpdateBy     string `form:"updateByOrder"  search:"type:order;column:update_by;table:sys_range"`
	CreatedAt    string `form:"createdAtOrder"  search:"type:order;column:created_at;table:sys_range"`
	UpdatedAt    string `form:"updatedAtOrder"  search:"type:order;column:updated_at;table:sys_range"`
	DeletedAt    string `form:"deletedAtOrder"  search:"type:order;column:deleted_at;table:sys_range"`
}

func (m *UseRangeGetPageReq) GetNeedSearch() interface{} {
	return *m
}

// SysRangeGetReq 功能获取请求参数
type UseRangeGetReq struct {
	RangeId     int    `uri:"id"`
	ProjectName string `json:"projectName"`
}

func (s *UseRangeGetReq) GetId() interface{} {
	return s.RangeId
}

type RangeChild struct {
	RangeId      int    `json:"rangeId"`
	RangeName    string `json:"rangeName"`
	Ipadress     string `json:"ipadress"`
	RangeConsole string `json:"rangeConsole"`
}

type Response struct {
	ProjectName string       `json:"projectName"`
	Dept        string       `json:"dept"`
	Children    []RangeChild `json:"children"`
}
