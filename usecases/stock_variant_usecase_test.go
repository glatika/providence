package usecase

import (
	"errors"
	"log"
	"os"
	"testing"

	"github.com/glatika/providence/model"
	"github.com/glatika/providence/model/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRegisterStockVariant(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	privkey := os.Getenv("PRIVKEYPATH")
	if privkey == "" {
		log.Fatalln("Set PRIVKEYPATH env var first")
	}

	stockVariant := model.NewStockVariant{
		Variant:        "",
		Permissions:    "",
		Abilities:      "",
		EncryptionCert: "",
	}

	stockVarItem := model.StockVariant{
		Id:          0,
		Variant:     stockVariant.Variant,
		Permissions: stockVariant.Permissions,
		Abilities:   stockVariant.Abilities,
	}

	t.Run("Success", func(t *testing.T) {
		stockVarRepo := mock.NewMockStockVariantRepository(ctrl)
		stockVarRepo.EXPECT().
			FindStockByVariant(stockVariant.Variant).
			Times(1).
			Return(nil, model.ErrRecordNotFound)

		stockVarRepo.EXPECT().
			Create(&stockVarItem).
			Times(1).
			Return(int32(1), nil)

		usecase := NewStockVariantCase(stockVarRepo)
		assert.NoError(t, usecase.RegisterStockVariant(stockVariant))
	})

	t.Run("Already existed", func(t *testing.T) {
		stockVarRepo := mock.NewMockStockVariantRepository(ctrl)
		stockVarRepo.EXPECT().
			FindStockByVariant(stockVariant.Variant).
			Times(1).
			Return(&stockVarItem, nil)

		usecase := NewStockVariantCase(stockVarRepo)
		assert.Error(t, errors.New("stock variant already registered"), usecase.RegisterStockVariant(stockVariant))
	})

	t.Run("Failed to find", func(t *testing.T) {
		stockVarRepo := mock.NewMockStockVariantRepository(ctrl)
		stockVarRepo.EXPECT().
			FindStockByVariant(stockVariant.Variant).
			Times(1).
			Return(nil, errors.New("Hehe"))

		usecase := NewStockVariantCase(stockVarRepo)
		assert.Error(t, usecase.RegisterStockVariant(stockVariant))
	})

	t.Run("Failed to create", func(t *testing.T) {
		stockVarRepo := mock.NewMockStockVariantRepository(ctrl)
		stockVarRepo.EXPECT().
			FindStockByVariant(stockVariant.Variant).
			Times(1).
			Return(nil, model.ErrRecordNotFound)

		stockVarRepo.EXPECT().
			Create(&stockVarItem).
			Times(1).
			Return(int32(0), errors.New("Hehe"))

		usecase := NewStockVariantCase(stockVarRepo)
		assert.Error(t, usecase.RegisterStockVariant(stockVariant))
	})

}
