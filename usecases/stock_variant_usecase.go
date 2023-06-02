package usecase

import (
	"errors"

	"github.com/glatika/providence/model"
)

type stockVariantCase struct {
	stockVarRepo model.StockVariantRepository
}

func NewStockVariantCase(stockVarRepo model.StockVariantRepository) model.StockVariantUsecase {
	return stockVariantCase{
		stockVarRepo: stockVarRepo,
	}
}

func (u stockVariantCase) RegisterStockVariant(stockVariant model.NewStockVariant) error {
	alreadyStockVar, err := u.stockVarRepo.FindStockByVariant(stockVariant.Variant)
	if err != nil {
		if !(err.Error() == model.ErrRecordNotFound.Error()) {
			return err
		}
	}

	if alreadyStockVar != nil {
		return errors.New("stock variant already registered")
	}

	_, err = u.stockVarRepo.Create(&model.StockVariant{
		Id:          0,
		Variant:     stockVariant.Variant,
		Abilities:   stockVariant.Abilities,
		Permissions: stockVariant.Permissions,
	})

	return err
}

func (u stockVariantCase) GetAllRegisteredStockVariant() ([]model.StockVariant, error) {
	return []model.StockVariant{}, nil
}
