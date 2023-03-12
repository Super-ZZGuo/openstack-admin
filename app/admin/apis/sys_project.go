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

type SysProject struct {
	api.Api
}

// GetPage 获取SysProject列表
// @Summary 获取SysProject列表
// @Description 获取SysProject列表
// @Tags SysProject
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.SysProject}} "{"code": 200, "data": [...]}"
// @Router /api/v1/sys-project [get]
// @Security Bearer
func (e SysProject) GetPage(c *gin.Context) {
	req := dto.SysProjectGetPageReq{}
	s := service.SysProject{}
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
	list := make([]models.SysProject, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取SysProject失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取SysProject
// @Summary 获取SysProject
// @Description 获取SysProject
// @Tags SysProject
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.SysProject} "{"code": 200, "data": [...]}"
// @Router /api/v1/sys-project/{id} [get]
// @Security Bearer
func (e SysProject) Get(c *gin.Context) {
	req := dto.SysProjectGetReq{}
	s := service.SysProject{}
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
	var object models.SysProject

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取SysProject失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建SysProject
// @Summary 创建SysProject
// @Description 创建SysProject
// @Tags SysProject
// @Accept application/json
// @Product application/json
// @Param data body dto.SysProjectInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/sys-project [post]
// @Security Bearer
func (e SysProject) Insert(c *gin.Context) {
	req := dto.SysProjectInsertReq{}
	s := service.SysProject{}
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

	client := models.CreateIdentityClient(models.CreateIdentityProvider("admin"))
	_, err = models.CreateProject(client, req.ProjectName, req.Tag)
	if err != nil {
		return
	}
	req.ProjectOpenstackId = models.GetProjectId(client, req.ProjectName)

	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建SysProject失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改SysProject
// @Summary 修改SysProject
// @Description 修改SysProject
// @Tags SysProject
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.SysProjectUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/sys-project/{id} [put]
// @Security Bearer
func (e SysProject) Update(c *gin.Context) {
	req := dto.SysProjectUpdateReq{}
	s := service.SysProject{}
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

	client := models.CreateIdentityClient(models.CreateIdentityProvider(req.OldProjectName))
	req.ProjectOpenstackId = models.GetProjectId(client, req.OldProjectName)
	_, err = models.UpateProject(client, req.NewProjectName, req.OldProjectName)
	if err != nil {
		return
	}

	err = s.Update(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改SysProject失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "修改成功")
}

// Delete 删除SysProject
// @Summary 删除SysProject
// @Description 删除SysProject
// @Tags SysProject
// @Param data body dto.SysProjectDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/sys-project [delete]
// @Security Bearer
func (e SysProject) Delete(c *gin.Context) {
	s := service.SysProject{}
	req := dto.SysProjectDeleteReq{}
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

	client := models.CreateIdentityClient(models.CreateIdentityProvider(req.ProjectName))
	err = models.DelteProject(client, req.ProjectName)
	if err != nil {
		return
	}

	err = s.Remove(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除SysProject失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}
