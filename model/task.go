package model

import (
	"context"
	"time"

	"github.com/glatika/providence/deliveries/grpc/market/market_pb"
	"github.com/glatika/providence/deliveries/grpc/stock/stock_pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type Task struct {
	Id          int32     `gorm:"primary_key" json:"id"`
	StockId     int32     `gorm:"column:stock_id" json:"stockid"`
	Instruction string    `gorm:"column:instruction" json:"instruction"`
	Argument    string    `gorm:"column:argument" json:"argument"`
	Delivered   bool      `gorm:"column:delivered" json:"delivered"`
	DeliveredAt time.Time `gorm:"type:datetime,column:delivered_at,default:CURRENT_TIMESTAMP()" json:"delivered_at"`
	Reported    bool      `gorm:"column:reported" json:"reported"`
	ReportedAt  time.Time `gorm:"type:datetime,column:reported_at,default:CURRENT_TIMESTAMP()" json:"reported_at"`
	Success     bool      `gorm:"column:success" json:"success"`
	Report      string    `gorm:"column:report" json:"report"`
}

type NewTask struct {
	StockId     int32
	Instruction string
	Argument    string
}

func (t Task) ToMarketProto() market_pb.Task {
	// var deliverAt *timestamppb.Timestamp = nil
	// if t.DeliveredAt != nil {
	// 	deliverAt = timestamppb.New(*t.DeliveredAt)
	// }

	// var reportAt *timestamppb.Timestamp = nil
	// if t.ReportedAt != nil {
	// 	reportAt = timestamppb.New(*t.ReportedAt)
	// }

	deliverAt := timestamppb.New(t.DeliveredAt)
	reportAt := timestamppb.New(t.ReportedAt)

	return market_pb.Task{
		Id:          t.Id,
		Instruction: t.Instruction,
		Argument:    t.Argument,
		Delivered:   t.Delivered,
		DeliveredAt: deliverAt,
		Reported:    t.Reported,
		ReportedAt:  reportAt,
		Success:     t.Success,
		Report:      t.Report,
	}
}

type TaskRepository interface {
	Create(t *Task) (int32, error)
	FindTaskById(id int) (*Task, error)
	GetAllTasks(page, limit int) ([]Task, error)
	FindUndeliveredTaskByStockId(stockid int32) (*Task, error)
	FindUnreportedTaskById(taskid int32) (*Task, error)
	UpdateDeliveredByID(id int) error
	UpdateReportedByID(id int, success bool, report string) error
}

type TaskUsecase interface {
	RequestTask(context.Context, *stock_pb.TaskRequest) (*stock_pb.TaskResponse, error)
	GetAllTask(page int, size int) (*[]Task, error)
	RegisterTask(req *NewTask) error
	ReportTask(req *stock_pb.StockLifecycle_ReportServer) error
}
