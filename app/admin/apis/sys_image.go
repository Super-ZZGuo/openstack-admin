package apis

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	_ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
)

type SysImage struct {
	api.Api
}

// GetPage 获取SysImage列表
// @Summary 获取SysImage列表
// @Description 获取SysImage列表
// @Tags SysImage
// @Param pageSize query int false "页条数"
// @Param pageIndex query int false "页码"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.SysImage}} "{"code": 200, "data": [...]}"
// @Router /api/v1/sys-image [get]
// @Security Bearer
func (e SysImage) GetPage(c *gin.Context) {
	req := dto.SysImageGetPageReq{}
	s := service.SysImage{}
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
	list := make([]models.SysImage, 0)
	var count int64

	err = s.GetPage(&req, p, &list, &count)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取SysImage失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查询成功")
}

// Get 获取SysImage
// @Summary 获取SysImage
// @Description 获取SysImage
// @Tags SysImage
// @Param id path int false "id"
// @Success 200 {object} response.Response{data=models.SysImage} "{"code": 200, "data": [...]}"
// @Router /api/v1/sys-image/{id} [get]
// @Security Bearer
func (e SysImage) Get(c *gin.Context) {
	req := dto.SysImageGetReq{}
	s := service.SysImage{}
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
	var object models.SysImage

	p := actions.GetPermissionFromContext(c)
	err = s.Get(&req, p, &object)
	// object.Path = ""
	if err != nil {
		e.Error(500, err, fmt.Sprintf("获取SysImage失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(object, "查询成功")
}

// Insert 创建SysImage
// @Summary 创建SysImage
// @Description 创建SysImage
// @Tags SysImage
// @Accept application/json
// @Product application/json
// @Param data body dto.SysImageInsertReq true "data"
// @Success 200 {object} response.Response	"{"code": 200, "message": "添加成功"}"
// @Router /api/v1/sys-image [post]
// @Security Bearer
func (e SysImage) Insert(c *gin.Context) {
	req := dto.SysImageInsertReq{}
	s := service.SysImage{}
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

	err = s.Insert(&req)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("创建SysImage失败，\r\n失败信息 %s", err.Error()))
		return
	}

	e.OK(req.GetId(), "创建成功")
}

// Update 修改SysImage
// @Summary 修改SysImage
// @Description 修改SysImage
// @Tags SysImage
// @Accept application/json
// @Product application/json
// @Param id path int true "id"
// @Param data body dto.SysImageGetUpdateReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "修改成功"}"
// @Router /api/v1/sys-image/{id} [put]
// @Security Bearer
func (e SysImage) Update(c *gin.Context) {
	req := dto.SysImageGetUpdateReq{}
	new := dto.SysImagePutUpdateReq{}
	s := service.SysImage{}
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

	new.ImageId = req.ImageId
	if req.ImageNewName != "" {
		client := models.CreateImageClient(models.CreateImageProvider("admin"))
		err = models.UpadteImage(client, req.ImageNewName, models.GetImageId(client, req.ImageOldName))
		if err != nil {
			e.Error(500, err, fmt.Sprintf("修改SysImage失败，\r\n失败信息 %s", err.Error()))
			return
		}
		new.ImageName = req.ImageNewName
	} else {
		new.ImageName = req.ImageOldName
	}
	if req.Newtag != "" {
		new.Tag = req.Newtag
	} else {
		new.Tag = req.OldTag
	}

	err = s.Update(&new, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("修改SysImage失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "修改成功")
}

// Delete 删除SysImage
// @Summary 删除SysImage
// @Description 删除SysImage
// @Tags SysImage
// @Param data body dto.SysImageDeleteReq true "body"
// @Success 200 {object} response.Response	"{"code": 200, "message": "删除成功"}"
// @Router /api/v1/sys-image [delete]
// @Security Bearer
func (e SysImage) Delete(c *gin.Context) {
	s := service.SysImage{}
	req := dto.SysImageDeleteReq{}
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
	client := models.CreateImageClient(models.CreateImageProvider("admin"))
	var openstackId string
	for _, v := range req.Names {
		openstackId = models.GetImageId(client, v)
		err = models.DeleteImage(client, openstackId)
		if err != nil {
			e.Error(500, err, fmt.Sprintf("删除SysImage失败，\r\n失败信息 %s", err.Error()))
			return
		}
	}

	err = s.Remove(&req, p)
	if err != nil {
		e.Error(500, err, fmt.Sprintf("删除SysImage失败，\r\n失败信息 %s", err.Error()))
		return
	}
	e.OK(req.GetId(), "删除成功")
}

//上传镜像
func (e SysImage) UploadSysImage(c *gin.Context) {
	s := service.SysImage{}
	req := dto.UploadSysImage{}
	err := e.MakeContext(c).
		MakeOrm().
		MakeService(&s.Service).
		Errors
	if err != nil {
		e.Logger.Error(err)
		e.Error(500, err, err.Error())
		return
	}
	// 数据权限检查
	p := actions.GetPermissionFromContext(c)
	form, _ := c.MultipartForm()
	file := form.File["upload"][0]
	filPath := "static/uploadfile/image/" + form.Value["type"][0] + "/" + form.Value["imageName"][0] + "." + form.Value["type"][0]
	fmt.Println(filPath)
	fmt.Println(filPath)
	fmt.Println(filPath)
	fmt.Println(filPath)

	e.Logger.Debugf("upload image file: %s", file.Filename)
	// 上传文件至指定目录
	err = c.SaveUploadedFile(file, filPath)
	if err != nil {
		e.Logger.Errorf("save file error, %s", err.Error())
		e.Error(500, err, "")
		return
	}

	req.ImageId, _ = strconv.Atoi(form.Value["imageId"][0])
	req.Type = form.Value["type"][0]
	req.ImageName = form.Value["imageName"][0]
	req.Path = filPath

	if strings.EqualFold(req.Type, "img") {
		req.Type = "qcow2"
	}
	client := models.CreateImageClient(models.CreateImageProvider("admin"))
	err = models.CreateImage(client, req.ImageName, req.Type, req.Path)
	if err != nil {
		e.Error(500, err, "")
		e.Logger.Error(err)
		return
	}
	req.Openstack_id = models.GetImageId(client, req.ImageName)
	err = s.UpdateImagePath(&req, p)
	if err != nil {
		e.Error(500, err, "上传失败")
		e.Logger.Error(err)
		return
	}
	e.OK(filPath, "上传成功")
}
