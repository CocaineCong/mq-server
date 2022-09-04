package dao

import (
	"context"
	"gorm.io/gorm"
	"mq-server/model"
)

type SkillGoodsDao struct {
	*gorm.DB
}

func NewSkillGoodsDao(ctx context.Context) *SkillGoodsDao {
	return &SkillGoodsDao{NewDBClient(ctx)}
}

func (dao *SkillGoodsDao) Create(in *model.SkillGoods) error {
	return dao.Model(&model.SkillGoods{}).Create(&in).Error
}

func (dao *SkillGoodsDao) CreateByList(in []*model.SkillGoods) error {
	return dao.Model(&model.SkillGoods{}).Create(&in).Error
}

func (dao *SkillGoodsDao) ListSkillGoods() (resp []*model.SkillGoods, err error) {
	err = dao.Model(&model.SkillGoods{}).Where("num > 0").Find(&resp).Error
	return
}

func (dao *SkillGoodsDao) UpdateNumByRedis(id uint, num int) error {
	return dao.Model(&model.SkillGoods{}).Where("id=?", id).
		Update("num", num).Error
}
