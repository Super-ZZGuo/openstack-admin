package service

import (
	"errors"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type SysImage struct {
	service.Service
}

// GetPage 获取SysImage列表
func (e *SysImage) GetPage(c *dto.SysImageGetPageReq, p *actions.DataPermission, list *[]models.SysImage, count *int64) error {
	var err error
	var data models.SysImage

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("SysImageService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取SysImage对象
func (e *SysImage) Get(d *dto.SysImageGetReq, p *actions.DataPermission, model *models.SysImage) error {
	var data models.SysImage

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetSysImage error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建SysImage对象
func (e *SysImage) Insert(c *dto.SysImageInsertReq) error {
	var err error
	var data models.SysImage
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("SysImageService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改SysImage对象
func (e *SysImage) Update(c *dto.SysImagePutUpdateReq, p *actions.DataPermission) error {
	var err error
	var model models.SysImage
	db := e.Orm.Scopes(
		actions.Permission(model.TableName(), p),
	).First(&model, c.GetId())
	if err = db.Error; err != nil {
		e.Log.Errorf("Service UpdateSysImage path error: %s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")

	}
	err = e.Orm.Table(model.TableName()).Where("image_id =? ", c.ImageId).Updates(c).Error
	if err != nil {
		e.Log.Errorf("Service UpdateSysImage path error: %s", err)
		return err
	}
	return nil
}

// Remove 删除SysImage
func (e *SysImage) Remove(d *dto.SysImageDeleteReq, p *actions.DataPermission) error {
	var data models.SysImage

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveSysImage error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}

func (e *SysImage) UpdateImagePath(c *dto.UploadSysImage, p *actions.DataPermission) error {
	var err error
	var model models.SysImage
	db := e.Orm.Scopes(
		actions.Permission(model.TableName(), p),
	).First(&model, c.GetId())
	if err = db.Error; err != nil {
		e.Log.Errorf("Service UpdateSysImage path error: %s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")

	}
	err = e.Orm.Table(model.TableName()).Where("image_id =? ", c.ImageId).Updates(c).Error
	if err != nil {
		e.Log.Errorf("Service UpdateSysImage path error: %s", err)
		return err
	}
	return nil
}
