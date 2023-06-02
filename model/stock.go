package model

import (
	"context"

	"github.com/glatika/providence/deliveries/grpc/market/market_pb"
	"github.com/glatika/providence/deliveries/grpc/stock/stock_pb"
)

type Stock struct {
	Id              int32  `gorm:"primary_key,column:id" json:"id"`
	Hwid            string `gorm:"column:hwid" json:"hwid"`
	Variant         string `gorm:"column:variant" json:"variant"`
	OperatingSystem string `gorm:"column:operating_system" json:"operating_system"`
	Signature       string `gorm:"column:signature" json:"signature"`
}

func (s Stock) ToMarketProto() market_pb.Stock {
	return market_pb.Stock{
		Stockid: s.Id,
		Variant: s.Variant,
		Os:      s.OperatingSystem,
	}
}

type NewStock struct {
	Hwid            string
	Variant         string
	OperatingSystem string
	Signature       string
}

type StockRepository interface {
	Create(n Stock) (int32, error)
	FindStockById(id int32) (*Stock, error)
	FindStockByHwid(hwid string) (*Stock, error)
	GetAllRepo(page, limit int) ([]Stock, error)
}

type StockUsecase interface {
	RegisterStock(context.Context, *stock_pb.RegisterRequest) (*stock_pb.RegisterResponse, error)
	GetAllStocks(page int, size int) ([]Stock, error)
}

func (n NewStock) ToStock() Stock {
	return Stock{
		Id:              0,
		Hwid:            n.Hwid,
		Signature:       n.Signature,
		Variant:         n.Variant,
		OperatingSystem: n.OperatingSystem,
	}
}

func (n *Stock) FromRegisterRequest(proto *stock_pb.RegisterRequest) {
	n.Hwid = proto.GetHwid()
	n.OperatingSystem = proto.GetOs()
	n.Variant = proto.GetVariant()
	n.Signature = proto.GetSignature()
}
