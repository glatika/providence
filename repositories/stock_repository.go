package repositories

import (
	"github.com/glatika/providence/model"

	"gorm.io/gorm"
)

type stockRepo struct {
	db *gorm.DB
}

func NewStockRepository(db *gorm.DB) model.StockRepository {
	return stockRepo{
		db: db,
	}
}

func (r stockRepo) Create(stock model.Stock) (int32, error) {
	result := r.db.Create(&stock)
	if result.Error != nil {
		return 0, result.Error
	}
	return stock.Id, nil
}

func (r stockRepo) FindStockById(id int32) (result *model.Stock, err error) {
	err = r.db.First(&result, id).Error
	return
}

func (r stockRepo) FindStockByHwid(hwid string) (result *model.Stock, err error) {
	err = r.db.First(&result, "hwid = ?", hwid).Error
	if err != nil {
		result = nil
	}
	return
}

func (r stockRepo) GetAllRepo(limit int, page int) ([]model.Stock, error) {
	var stocks = []model.Stock{}
	err := r.db.Offset((page - 1) * limit).Limit(limit).Find(&stocks).Error
	return stocks, err
}
