package service

import (
	"errors"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	models_admin "go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	models_other "go-admin/app/other/models/use"
	dto_other "go-admin/app/other/service/dto"
	"go-admin/common/actions"
	cDto "go-admin/common/dto"
)

type UseRange struct {
	service.Service
}

// GetPage 获取UseRange列表
func (e *UseRange) GetPage(c *dto_other.UseRangeGetPageReq, p *actions.DataPermission, list *[]models_other.UseRange, count *int64) error {
	var err error
	var data models_other.UseRange

	err = e.Orm.Model(&data).
		Scopes(
			cDto.MakeCondition(c.GetNeedSearch()),
			cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
			actions.Permission(data.TableName(), p),
		).
		Find(list).Limit(-1).Offset(-1).
		Count(count).Error
	if err != nil {
		e.Log.Errorf("UseRangeService GetPage error:%s \r\n", err)
		return err
	}
	return nil
}

// Get 获取UseRange对象
func (e *UseRange) Get(d *dto_other.UseRangeGetReq, p *actions.DataPermission, model *models_other.UseRange) error {
	var data models_other.UseRange

	err := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).
		First(model, d.GetId()).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("Service GetUseRange error:%s \r\n", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

func (e *UseRange) GetDept(d *dto.SysDeptGetReq, model *models_admin.SysDept) error {
	var err error
	var data models_admin.SysDept

	db := e.Orm.Model(&data).
		First(model, d.GetId())
	err = db.Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("db error:%s", err)
		return err
	}
	if err = db.Error; err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}
