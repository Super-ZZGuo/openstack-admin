package apis

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

type SysRange struct {
	api.Api
}

// GetPage 获取SysRange列表
// @Summary 获取SysRange列表
// @Description 获取SysRange列表
// @Tags SysRange
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.SysRange}} "{"code": 200, "data": [...]}"
// @Router /api/v1/sys-range [get]
// @Security Bearer
func (e SysRange) GetPage(c *gin.Context) {
	req := dto.SysRangeGetPageReq{}
	s := service.SysRange{}
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
	list := make([]models.SysRange, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取SysRange失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取SysRange
// @Summary 获取SysRange
// @Description 获取SysRange
// @Tags SysRange
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.SysRange} "{"code": 200, "data": [...]}"
// @Router /api/v1/sys-range/{id} [get]
// @Security Bearer
func (e SysRange) Get(c *gin.Context) {
	req := dto.SysRangeGetReq{}
	s := service.SysRange{}
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
	var object models.SysRange

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取SysRange失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建SysRange
// @Summary 创建SysRange
// @Description 创建SysRange
// @Tags SysRange
// @Accept application/json
// @Product application/json
// @Param data body dto.SysRangeInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/sys-range [post]
// @Security Bearer
func (e SysRange) Insert(c *gin.Context) {
	req := dto.SysRangeInsertReq{}
	s := service.SysRange{}
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

	client := models.CreateComputeClient()
	createOpts := servers.CreateOpts{
		Name:      req.RangeName,
		ImageRef:  req.Image,
		FlavorRef: req.Flavor,
	}
	_, err = servers.Create(client, createOpts).Extract()
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建SysRange失败，\r\n失败信息 %s", err.Error()))
		return
	}

	// 设置创建人
	req.SetCreateBy(user.GetUserId(c))

	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建SysRange失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改SysRange
// @Summary 修改SysRange
// @Description 修改SysRange
// @Tags SysRange
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.SysRangeUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/sys-range/{id} [put]
// @Security Bearer
func (e SysRange) Update(c *gin.Context) {
	req := dto.SysRangeUpdateReq{}
	s := service.SysRange{}
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
	req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)

	err = s.Update(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改SysRange失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除SysRange
// @Summary 删除SysRange
// @Description 删除SysRange
// @Tags SysRange
// @Param data body dto.SysRangeDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/sys-range [delete]
// @Security Bearer
func (e SysRange) Delete(c *gin.Context) {
	s := service.SysRange{}
	req := dto.SysRangeDeleteReq{}
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
	client := models.CreateComputeClient()
	for _, serverID := range req.RangeOpenstackID {
		servers.Delete(client, serverID)
	}

	// req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)

	err = s.Remove(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除SysRange失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
