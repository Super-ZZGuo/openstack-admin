package service

import (
	"errors"
	"fmt"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type SysRange struct {
	service.Service
}

// GetPage 获取SysRange列表
func (e *SysRange) GetPage(c *dto.SysRangeGetPageReq, p *actions.DataPermission, list *[]models.SysRange, count *int64) error {
	var err error
	var data models.SysRange

	err = e.Orm.Model(&data).Preload("Project").
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("SysRangeService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取SysRange对象
func (e *SysRange) Get(d *dto.SysRangeGetReq, p *actions.DataPermission, model *models.SysRange) error {
	var data models.SysRange

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetSysRange error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// Insert 创建SysRange对象
func (e *SysRange) Insert(c *dto.SysRangeInsertReq) error {
	var err error
	var data models.SysRange
	c.Generate(&data)
	err = e.Orm.Create(&data).Error
	if err != nil {
		e.Log.Errorf("SysRangeService Insert error:%s \r\n", err)
		return err
	}
	return nil
}

// Update 修改SysRange对象
func (e *SysRange) Update(c *dto.SysRangeUpdateReq, p *actions.DataPermission) error {
	var err error
	var data = models.SysRange{}
	e.Orm.Scopes(
		actions.Permission(data.TableName(), p),
	).First(&data, c.GetId())
	c.Generate(&data)

	db := e.Orm.Save(&data)
	if err = db.Error; err != nil {
		e.Log.Errorf("SysRangeService Save error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// Remove 删除SysRange
func (e *SysRange) Remove(d *dto.SysRangeDeleteReq, p *actions.DataPermission) error {
	var data models.SysRange

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	fmt.Println(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveSysRange error:%s \r\n", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
	}
	return nil
}
