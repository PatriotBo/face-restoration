package dao

import (
	"context"
	"face-restoration/internal/conf"
	"face-restoration/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBDao interface {
	GetPredictResult(ctx context.Context, opts PredictOption) ([]*model.PredictRecord, error)
	CreatePredictRecord(ctx context.Context, r *model.PredictRecord) error
	UpdatePredictRecord(ctx context.Context, r *model.PredictRecord) error
}

type dbDao struct {
	db *gorm.DB
}

func NewDao() DBDao {
	eng, err := gorm.Open(mysql.Open(conf.GetDSN()), nil)
	if err != nil {
		panic(err)
	}
	return &dbDao{
		db: eng,
	}
}

type PredictOption struct {
	OpenID string
	Status int
}

// GetPredictResult get predict result which status is processing
func (d *dbDao) GetPredictResult(ctx context.Context, opts PredictOption) ([]*model.PredictRecord, error) {
	var list []*model.PredictRecord
	return list, d.db.WithContext(ctx).
		Where("open = ? AND status = ?", opts.OpenID, opts.Status).
		Find(&list).Error
}

// CreatePredictRecord insert record
func (d *dbDao) CreatePredictRecord(ctx context.Context, r *model.PredictRecord) error {
	return d.db.WithContext(ctx).Create(r).Error
}

// UpdatePredictRecord update record
func (d *dbDao) UpdatePredictRecord(ctx context.Context, r *model.PredictRecord) error {
	return d.db.WithContext(ctx).Model(new(model.PredictRecord)).Updates(r).Error
}
