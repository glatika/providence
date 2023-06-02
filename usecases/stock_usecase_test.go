package usecase

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"

	"github.com/glatika/providence/deliveries/grpc/stock/stock_pb"
	"github.com/glatika/providence/model"
	"github.com/glatika/providence/model/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestRegisterStock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	privkey := os.Getenv("PRIVKEYPATH")
	if privkey == "" {
		log.Fatalln("Set PRIVKEYPATH env var first")
	}

	req := stock_pb.RegisterRequest{
		Hwid:      "{123-456-123}",
		Signature: "Emirue",
		Variant:   "A-Winx64",
		Os:        "Windows",
	}

	stock := model.Stock{
		Id:              0,
		Hwid:            req.Hwid,
		Signature:       req.Signature,
		Variant:         req.Variant,
		OperatingSystem: req.Os,
	}

	variant := model.StockVariant{
		Id:          1,
		Variant:     req.Variant,
		Permissions: "",
		Abilities:   "",
	}

	t.Run("Success", func(t *testing.T) {
		stockRep := mock.NewMockStockRepository(ctrl)
		stockRep.EXPECT().
			FindStockByHwid(req.Hwid).
			Times(1).
			Return(nil, model.ErrRecordNotFound)

		stockVarRep := mock.NewMockStockVariantRepository(ctrl)
		stockVarRep.EXPECT().
			FindStockByVariant(req.Variant).
			Times(1).
			Return(&variant, nil)

		stockRep.EXPECT().
			Create(stock).
			Times(1).
			Return(stock.Id, nil)

		u, err := NewStockUsecase(stockRep, stockVarRep, privkey)
		require.NoError(t, err)
		response, err := u.RegisterStock(context.TODO(), &req)
		require.NoError(t, err)
		require.NotEqual(t, response, nil)
	})

	t.Run("Stock unavailable", func(t *testing.T) {
		stockRep := mock.NewMockStockRepository(ctrl)
		stockRep.EXPECT().
			FindStockByHwid(req.Hwid).
			Times(1).
			Return(&stock, nil)

		stockVarRep := mock.NewMockStockVariantRepository(ctrl)

		u, err := NewStockUsecase(stockRep, stockVarRep, privkey)
		require.NoError(t, err)
		response, err := u.RegisterStock(context.TODO(), &req)
		require.Nil(t, response)
		require.Equal(t, errors.New("unavailable"), err)
	})

	t.Run("Stock Variant Not Found", func(t *testing.T) {
		stockRep := mock.NewMockStockRepository(ctrl)
		stockRep.EXPECT().
			FindStockByHwid(req.Hwid).
			Times(1).
			Return(nil, model.ErrRecordNotFound)

		stockVarRep := mock.NewMockStockVariantRepository(ctrl)
		stockVarRep.EXPECT().
			FindStockByVariant(req.Variant).
			Times(1).
			Return(nil, nil)

		u, err := NewStockUsecase(stockRep, stockVarRep, privkey)
		require.NoError(t, err)
		response, err := u.RegisterStock(context.TODO(), &req)
		require.Error(t, err, errors.New("variant not found"))
		require.Nil(t, response)
	})

	t.Run("Failed to create", func(t *testing.T) {
		stockRep := mock.NewMockStockRepository(ctrl)
		stockRep.EXPECT().
			FindStockByHwid(req.Hwid).
			Times(1).
			Return(nil, model.ErrRecordNotFound)

		stockVarRep := mock.NewMockStockVariantRepository(ctrl)
		stockVarRep.EXPECT().
			FindStockByVariant(req.Variant).
			Times(1).
			Return(&variant, nil)

		stockRep.EXPECT().
			Create(stock).
			Times(1).
			Return(int32(0), errors.New("Hehe"))

		u, err := NewStockUsecase(stockRep, stockVarRep, privkey)
		require.NoError(t, err)
		response, err := u.RegisterStock(context.TODO(), &req)
		require.Error(t, err, errors.New("Hehe"))
		require.Nil(t, response)
	})

	t.Run("DB Query Error FindStock By Variant", func(t *testing.T) {
		stockRep := mock.NewMockStockRepository(ctrl)
		stockRep.EXPECT().
			FindStockByHwid(req.Hwid).
			Times(1).
			Return(nil, model.ErrRecordNotFound)

		stockVarRep := mock.NewMockStockVariantRepository(ctrl)
		stockVarRep.EXPECT().
			FindStockByVariant(req.Variant).
			Times(1).
			Return(nil, errors.New("variant not found"))

		u, err := NewStockUsecase(stockRep, stockVarRep, privkey)
		require.NoError(t, err)
		response, err := u.RegisterStock(context.TODO(), &req)
		require.Error(t, err, errors.New("variant not found"))
		require.Nil(t, response)
	})

	t.Run("DB Query Error FindSTockByHWID", func(t *testing.T) {
		stockRep := mock.NewMockStockRepository(ctrl)
		stockRep.EXPECT().
			FindStockByHwid(req.Hwid).
			Times(1).
			Return(nil, errors.New("hehe"))

		stockVarRep := mock.NewMockStockVariantRepository(ctrl)

		u, err := NewStockUsecase(stockRep, stockVarRep, privkey)
		require.NoError(t, err)
		response, err := u.RegisterStock(context.TODO(), &req)
		require.Nil(t, response)
		require.Equal(t, errors.New("hehe"), err)
	})

}
