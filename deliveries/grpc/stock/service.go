package stock

import (
	context "context"

	"github.com/glatika/providence/deliveries/grpc/stock/stock_pb"
	"github.com/glatika/providence/model"
)

type StockLifecycle struct {
	TaskCase  model.TaskUsecase
	StockCase model.StockUsecase
	stock_pb.UnimplementedStockLifecycleServer
}

func (s StockLifecycle) Request(c context.Context, r *stock_pb.TaskRequest) (*stock_pb.TaskResponse, error) {
	// TODO: add token verification from gRPC meta, is user use token and registered
	return s.TaskCase.RequestTask(c, r)
}

func (s StockLifecycle) Register(c context.Context, r *stock_pb.RegisterRequest) (*stock_pb.RegisterResponse, error) {
	return s.StockCase.RegisterStock(c, r)
}

func (s StockLifecycle) Report(r stock_pb.StockLifecycle_ReportServer) error {
	// TODO: add token verification from gRPC meta, is user use token and registered
	return s.TaskCase.ReportTask(&r)
}
