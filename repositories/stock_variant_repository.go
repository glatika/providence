package repositories

import (
	"github.com/glatika/providence/model"

	"gorm.io/gorm"
)

type stockVariantRepo struct {
	db *gorm.DB
}

func NewStockVariantRepository(db *gorm.DB) model.StockVariantRepository {
	return stockVariantRepo{
		db: db,
	}
}

func (r stockVariantRepo) Create(n *model.StockVariant) (int32, error) {
	if err := r.db.Create(n).Error; err != nil {
		return 0, err
	}
	return n.Id, nil
}

func (r stockVariantRepo) FindStockByVariant(v string) (result *model.StockVariant, err error) {
	result = &model.StockVariant{}
	if err = r.db.Table("stock_variants").First(&result, "variant = ?", v).Error; err != nil {
		return nil, err
	}
	return
}
