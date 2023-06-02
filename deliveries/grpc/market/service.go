package market

import (
	"context"
	"fmt"

	"github.com/glatika/providence/deliveries/grpc/market/market_pb"
	"github.com/glatika/providence/model"
)

type BarnMarketService struct {
	StockVariantCase model.StockVariantUsecase
	StockCase        model.StockUsecase
	TaskCase         model.TaskUsecase
	market_pb.UnimplementedBarnMarketProviderServer
}

func (s BarnMarketService) GetAllStock(req *market_pb.GetAllPagingRequest, res market_pb.BarnMarketProvider_GetAllStockServer) error {
	stocks, err := s.StockCase.GetAllStocks(int(req.GetPage()), int(req.GetSize()))
	if err != nil {
		return err
	}

	o := 0
	for o < len(stocks) {
		stock := (stocks)[o].ToMarketProto()
		if err = res.Send(&stock); err != nil {
			return err
		}
		o++
	}
	return nil
}

func (s BarnMarketService) GetAllStockTasks(req *market_pb.GetAllPagingRequest, res market_pb.BarnMarketProvider_GetAllStockTasksServer) error {
	stocks, err := s.TaskCase.GetAllTask(int(req.GetPage()), int(req.GetSize()))
	if err != nil {
		return err
	}

	o := 0
	for o < len(*stocks) {
		stock := (*stocks)[o].ToMarketProto()
		if err = res.Send(&stock); err != nil {
			return err
		}
		o++
	}
	return nil
}
func (s BarnMarketService) RegisterStockVariant(_ context.Context, req *market_pb.RegisterStockVariantRequest) (*market_pb.Empty, error) {
	return &market_pb.Empty{}, s.StockVariantCase.RegisterStockVariant(model.NewStockVariant{
		Variant:        req.Variant,
		Permissions:    req.Permission,
		Abilities:      req.Ability,
		EncryptionCert: req.Certificate,
	})
}

func (s BarnMarketService) RegisterTaskToStock(_ context.Context, req *market_pb.RegisterTaskToStockRequest) (*market_pb.RegisterTaskToStockResponse, error) {
	fmt.Println(req)

	err := s.TaskCase.RegisterTask(&model.NewTask{
		StockId:     req.GetStockId(),
		Instruction: req.GetInstruction(),
		Argument:    req.GetArgument(),
	})
	if err != nil {
		return &market_pb.RegisterTaskToStockResponse{
			Able:   false,
			Status: err.Error(),
		}, nil
	}
	return &market_pb.RegisterTaskToStockResponse{
		Able:   true,
		Status: "delivered",
	}, nil
}
