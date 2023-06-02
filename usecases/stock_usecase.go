package usecase

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/glatika/providence/deliveries/grpc/stock/stock_pb"
	"github.com/glatika/providence/model"

	jwt "github.com/golang-jwt/jwt/v5"
)

type stockCase struct {
	stockRepo    model.StockRepository
	stockVarRepo model.StockVariantRepository
	jwtPrivKey   *ecdsa.PrivateKey
}

func NewStockUsecase(stockRepo model.StockRepository, stockVarRepo model.StockVariantRepository, privKeyPath string) (model.StockUsecase, error) {
	key, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failure reading private key : %s", err)
	}

	parsedKey, err := jwt.ParseECPrivateKeyFromPEM(key)
	if err != nil {
		return nil, fmt.Errorf("failure parsing private key : %s", err)
	}

	return stockCase{
		stockRepo:    stockRepo,
		stockVarRepo: stockVarRepo,
		jwtPrivKey:   parsedKey,
	}, nil
}

// RegisterStock :nodoc:
func (u stockCase) RegisterStock(c context.Context, req *stock_pb.RegisterRequest) (*stock_pb.RegisterResponse, error) {
	stock, err := u.stockRepo.FindStockByHwid(req.Hwid)
	if err != nil {
		if !(err.Error() == model.ErrRecordNotFound.Error()) {
			return nil, err
		}
	}
	if stock != nil {
		return nil, errors.New("unavailable")
	}

	variant, err := u.stockVarRepo.FindStockByVariant(req.Variant)
	if err != nil {
		if !(err.Error() == model.ErrRecordNotFound.Error()) {
			return nil, err
		}
	}

	if variant == nil {
		return nil, errors.New("uncompatible with server, can not register")
	}

	stock = &model.Stock{}
	stock.FromRegisterRequest(req)
	id, err := u.stockRepo.Create(*stock)
	if err != nil {
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES384, jwt.MapClaims{
		"mine":      id,
		"hwmine":    stock.Hwid,
		"signature": stock.Signature,
		"variant":   stock.Variant,
	})

	signedToken, err := token.SignedString(u.jwtPrivKey)

	return &stock_pb.RegisterResponse{
		Token: signedToken,
	}, err
}

func (u stockCase) GetAllStocks(page int, size int) ([]model.Stock, error) {
	return u.stockRepo.GetAllRepo(size, page)
}
