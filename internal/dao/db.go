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
	ListProcessingRecords(ctx context.Context) ([]*model.PredictRecord, error)

	CreateOrder(ctx context.Context, order *model.Orders) error
	UpdateOrder(ctx context.Context, order *model.Orders) error
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

// ListProcessingRecords get records which status is processing,to fetch results for them.
func (d *dbDao) ListProcessingRecords(ctx context.Context) ([]*model.PredictRecord, error) {
	var list []*model.PredictRecord
	return list, d.db.WithContext(ctx).
		Where("status = ?", model.Processing).
		Order("create_time ASC").
		Find(&list).Error
}

func (d *dbDao) CreateOrder(ctx context.Context, order *model.Orders) error {
	return d.db.WithContext(ctx).Create(order).Error
}

func (d *dbDao) UpdateOrder(ctx context.Context, order *model.Orders) error {
	return d.db.WithContext(ctx).
		Model(new(model.Orders)).
		Where("subscribe_id = ?", order.SubscribeID).
		Updates(order).Error
}
