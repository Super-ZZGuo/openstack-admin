package apis

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	models_admin "go-admin/app/admin/models"
	dto_admin "go-admin/app/admin/service/dto"
	models_other "go-admin/app/other/models/use"
	service_other "go-admin/app/other/service"
	dto_other "go-admin/app/other/service/dto"
	"go-admin/common/actions"
)

type UseRange struct {
	api.Api
}

// GetPage 获取UseRange列表
// @Summary 获取UseRange列表
// @Description 获取UseRange列表
// @Tags UseRange
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models_other.UseRange}} "{"code": 200, "data": [...]}"
// @Router /api/v1/sys-image [get]
// @Security Bearer
func (e UseRange) GetPage(c *gin.Context) {
	req := dto_other.UseRangeGetPageReq{}
	dReq := dto_admin.SysDeptGetReq{}
	s := service_other.UseRange{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}

	p := actions.GetPermissionFromContext(c)

	//获取部门信息
	dReq.Id = p.DeptId
	var object models_admin.SysDept
	err = s.GetDept(&dReq, &object)
	if err != nil {
		e.Error(500, err, "查询失败")
		return
	}

	if object.DeptName == "监管人员" {
		object.DeptName = ""
	}
	req.Dept = object.DeptName

	pReq := dto_admin.SysProjectGetPageReq{
		Status: "2", //2为开启状态
	}
	pList := make([]models_admin.SysProject, 0)
	var countP int64
	err = s.GetProjectPage(&pReq, p, &pList, &countP)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取SysProject失败，\r\n失败信息 %s", err.Error()))
		return
	}
	responseList := make([]dto_other.Response, 0)
	if countP == 0 {
		if req.Dept == "" {
			responseList = append(responseList, dto_other.Response{
				ProjectName: "",
				Dept:        "监管人员",
				Children:    []dto_other.RangeChild{},
			})
		} else {
			responseList = append(responseList, dto_other.Response{
				ProjectName: "",
				Dept:        req.Dept,
				Children:    []dto_other.RangeChild{},
			})
		}
		e.PageOK(responseList, 0, req.GetPageIndex(), req.GetPageSize(), "查询成功")
		return

	}

	req.ProjectName = pList[0].ProjectName

	list := make([]models_other.UseRange, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取UseRange失败，\r\n失败信息 %s", err.Error()))
		return
	}

	responseList = append(responseList, dto_other.Response{
		ProjectName: req.ProjectName,
		Dept:        req.Dept,
		Children:    []dto_other.RangeChild{},
	})
	for _, item := range list {
		item.ProjectName = ""
		responseList[0].Children = append(responseList[0].Children, dto_other.RangeChild{
			RangeName:    item.RangeName,
			RangeId:      item.RangeId,
			Ipadress:     item.Ipadress,
			RangeConsole: item.RangeConsole,
		})

	}

	if req.Dept == "" && len(responseList) != 0 {
		responseList[0].Dept = "监管人员"
	}

	e.PageOK(responseList, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取UseRange
// @Summary 获取UseRange
// @Description 获取UseRange
// @Tags UseRange
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models_other.UseRange} "{"code": 200, "data": [...]}"
// @Router /api/v1/sys-image/{id} [get]
// @Security Bearer
func (e UseRange) Get(c *gin.Context) {
	req := dto_other.UseRangeGetReq{}
	s := service_other.UseRange{}
	err := e.MakeContext(c).
		MakeOrm().
		Bind(&req).
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	var object models_other.UseRange

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	// object.Path = ""
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取UseRange失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}
