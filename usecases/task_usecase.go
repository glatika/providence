package usecase

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/glatika/providence/deliveries/grpc/stock/stock_pb"
	"github.com/glatika/providence/model"

	jwt "github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/metadata"
)

type taskCase struct {
	taskRepo     model.TaskRepository
	stockVarRepo model.StockVariantRepository
	stockRepo    model.StockRepository
	pubKey       *ecdsa.PublicKey
}

func NewTaskUsecase(taskRepo model.TaskRepository, stockVarRepo model.StockVariantRepository, stockRepo model.StockRepository, pubKeyPath string) (model.TaskUsecase, error) {
	key, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		return nil, err
	}

	parsedKey, err := jwt.ParseECPublicKeyFromPEM(key)
	if err != nil {
		return nil, err
	}

	return taskCase{
		taskRepo:     taskRepo,
		stockVarRepo: stockVarRepo,
		stockRepo:    stockRepo,
		pubKey:       parsedKey,
	}, nil
}

func (u taskCase) RequestTask(c context.Context, req *stock_pb.TaskRequest) (*stock_pb.TaskResponse, error) {
	md, ok := metadata.FromIncomingContext(c)
	if !ok {
		return nil, model.ErrFailedReadGRPCMetadata
	}

	tokens := md.Get("token")
	if len(tokens) == 0 {
		return nil, model.ErrFailedPrecondition
	}

	token, err := jwt.Parse(tokens[0], func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodECDSA)
		if !ok {
			return nil, model.ErrFailedPrecondition
		}

		return u.pubKey, nil

	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return nil, model.ErrFailedPrecondition
	}

	fmt.Println("THIS the mine id after requesting task", claims["mine"])
	id := int(claims["mine"].(float64))

	task, err := u.taskRepo.FindUndeliveredTaskByStockId(int32(id))
	if err != nil {
		return nil, err
	}

	_ = u.taskRepo.UpdateDeliveredByID(int(task.Id))

	gTask := stock_pb.TaskResponse{
		Taskid:      task.Id,
		Instruction: task.Instruction,
		Arg:         task.Argument,
	}

	return &gTask, nil
}

func (u taskCase) GetAllTask(page int, size int) (*[]model.Task, error) {
	stock, err := u.taskRepo.GetAllTasks(page, size)
	return &stock, err
}

func (u taskCase) RegisterTask(req *model.NewTask) error {
	fmt.Println(req)
	if req == nil {
		return errors.New("herre nil")
	}
	// TODO: checking if the stock exists
	stock, err := u.stockRepo.FindStockById(req.StockId)
	if err != nil {
		return err
	}

	if stock == nil {
		return model.ErrStockNotFound
	}

	// TODO: checking if the insturction is available for stock ability

	// TODO: checking if the stock are had permission to do ability in current host

	task := model.Task{
		Id:          0,
		StockId:     req.StockId,
		Instruction: req.Instruction,
		Argument:    req.Argument,
		Delivered:   false,
		Reported:    false,
		Success:     false,
		Report:      "",
	}

	_, err = u.taskRepo.Create(&task)
	return err
}

func writeToFile(f *os.File, data []byte) error {
	writtenSize := 0
	dataSize := len(data)
	for {
		writtenSize, err := f.Write(data[writtenSize:])
		if err != nil {
			return err
		}

		if writtenSize >= dataSize {
			return nil
		}
	}

}

func (u taskCase) ReportTask(req *stock_pb.StockLifecycle_ReportServer) error {

	reqData, err := (*req).Recv()

	if reqData.GetReport().GetIsFile() == stock_pb.ReportType_File {
		firstChunck := true
		var file *os.File
		for {
			if firstChunck {
				file, err = os.Open(path.Join("./test/", filepath.Base(reqData.GetReport().GetFilename())))
				if err != nil {
					return err
				}
			}
			if !firstChunck {
				reqData, err = (*req).Recv()
			}
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}

			err = writeToFile(file, reqData.GetFileContent())
			if err != nil {
				return err
			}
		}
	}
	task, err := u.taskRepo.FindUnreportedTaskById(reqData.GetReport().Taskid)
	if err != nil {
		return err
	}

	return u.taskRepo.UpdateReportedByID(int(task.Id), reqData.GetReport().Success, reqData.GetReport().Report)
}
