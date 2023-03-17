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
)

type SysNetwork struct {
	api.Api
}

// GetPage 获取SysNetwork列表
// @Summary 获取SysNetwork列表
// @Description 获取SysNetwork列表
// @Tags SysNetwork
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.SysNetwork}} "{"code": 200, "data": [...]}"
// @Router /api/v1/sys-network [get]
// @Security Bearer
func (e SysNetwork) GetPage(c *gin.Context) {
	req := dto.SysNetworkGetPageReq{}
	s := service.SysNetwork{}
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
	list := make([]models.SysNetwork, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取SysNetwork失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取SysNetwork
// @Summary 获取SysNetwork
// @Description 获取SysNetwork
// @Tags SysNetwork
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.SysNetwork} "{"code": 200, "data": [...]}"
// @Router /api/v1/sys-network/{id} [get]
// @Security Bearer
func (e SysNetwork) Get(c *gin.Context) {
	req := dto.SysNetworkGetReq{}
	s := service.SysNetwork{}
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
	var object models.SysNetwork

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取SysNetwork失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建SysNetwork
// @Summary 创建SysNetwork
// @Description 创建SysNetwork
// @Tags SysNetwork
// @Accept application/json
// @Product application/json
// @Param data body dto.SysNetworkInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/sys-network [post]
// @Security Bearer
func (e SysNetwork) Insert(c *gin.Context) {
	req := dto.SysNetworkInsertReq{}
	s := service.SysNetwork{}
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
	// 设置创建人
	req.SetCreateBy(user.GetUserId(c))

	client := models.CreateNetworkClient(models.CreateNetworkProvider(req.ProjectName))
	err = models.CreateNetwork(client, req.NetworkName, req.Cidr, req.PoolStart, req.PoolEnd)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建SysNetwork失败，\r\n失败信息 %s", err.Error()))
		return
	}

	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建SysNetwork失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改SysNetwork
// @Summary 修改SysNetwork
// @Description 修改SysNetwork
// @Tags SysNetwork
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.SysNetworkUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/sys-network/{id} [put]
// @Security Bearer
func (e SysNetwork) Update(c *gin.Context) {
	req := dto.SysNetworkUpdateReq{}
	new := dto.SysNetworkPutUpdateReq{}
	s := service.SysNetwork{}
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

	new.NetworkId = req.NetworkId
	if req.NetworkNewName != "" {
		client := models.CreateNetworkClient(models.CreateNetworkProvider(req.ProjectName))
		err = models.UpadteNetwork(client, req.NetworkNewName, req.NetworkOldName)
		if err != nil {
			e.Error(500, err, fmt.Sprintf("修改SysImage失败，\r\n失败信息 %s", err.Error()))
			return
		}
		new.NetworkName = req.NetworkNewName
	} else {
		new.NetworkName = req.NetworkOldName
	}
	if req.Newtag != "" {
		new.Tag = req.Newtag
	} else {
		new.Tag = req.OldTag
	}

	err = s.Update(&new, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改SysNetwork失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除SysNetwork
// @Summary 删除SysNetwork
// @Description 删除SysNetwork
// @Tags SysNetwork
// @Param data body dto.SysNetworkDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/sys-network [delete]
// @Security Bearer
func (e SysNetwork) Delete(c *gin.Context) {
	s := service.SysNetwork{}
	req := dto.SysNetworkDeleteReq{}
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

	// req.SetUpdateBy(user.GetUserId(c))
	p := actions.GetPermissionFromContext(c)

	client := models.CreateNetworkClient(models.CreateFlavorProvider(req.ProjectName))
	for _, name := range req.NetworkNames {
		err = models.DeleteNetwork(client, name)
		if err != nil {
			e.Error(500, err, fmt.Sprintf("删除SysFlavor失败，\r\n失败信息 %s", err.Error()))
			return
		}
	}

	err = s.Remove(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除SysNetwork失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
