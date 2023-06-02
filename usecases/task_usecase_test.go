package usecase

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/glatika/providence/deliveries/grpc/stock/stock_pb"
	"github.com/glatika/providence/model"
	"github.com/glatika/providence/model/mock"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestRequestTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	privkey := os.Getenv("PRIVKEYPATH")
	if privkey == "" {
		log.Fatalln("Set PRIVKEYPATH env var first")
	}

	pubkey := os.Getenv("PUBKEYPATH")
	if pubkey == "" {
		log.Fatalln("Set PUBKEYPATH env var first")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES384, jwt.MapClaims{
		"mine":      1,
		"hwmine":    "123",
		"signature": "Emirue",
		"variant":   "A",
	})

	key, err := ioutil.ReadFile(privkey)
	assert.NoError(t, err)

	parsedKey, err := jwt.ParseECPrivateKeyFromPEM(key)
	assert.NoError(t, err)

	signedToken, err := token.SignedString(parsedKey)
	assert.NoError(t, err)

	req := stock_pb.TaskRequest{}

	md := metadata.MD{
		"token": []string{signedToken},
	}

	ctx := metadata.NewIncomingContext(context.TODO(), md)

	taskItem := model.Task{
		Id:          0,
		StockId:     1,
		Instruction: "Launch Missiles",
		Argument:    "Nowhere",
		Delivered:   false,
		Success:     false,
		Reported:    false,
	}

	t.Run("Succes", func(t *testing.T) {
		taskRepo := mock.NewMockTaskRepository(ctrl)
		taskRepo.EXPECT().
			FindUndeliveredTaskByStockId(int32(1)).
			Times(1).
			Return(&taskItem, nil)

		taskRepo.EXPECT().
			UpdateDeliveredByID(int(taskItem.Id)).
			Times(1).
			Return(nil)

		stockVarRepo := mock.NewMockStockVariantRepository(ctrl)
		stockRepo := mock.NewMockStockRepository(ctrl)

		usecase, err := NewTaskUsecase(taskRepo, stockVarRepo, stockRepo, pubkey)
		assert.NoError(t, err)
		task, err := usecase.RequestTask(ctx, &req)
		assert.NoError(t, err)
		assert.Equal(t, taskItem.Id, task.Taskid)
	})

	t.Run("Error GRPC Metadata", func(t *testing.T) {
		taskRepo := mock.NewMockTaskRepository(ctrl)
		stockVarRepo := mock.NewMockStockVariantRepository(ctrl)
		stockRepo := mock.NewMockStockRepository(ctrl)

		usecase, err := NewTaskUsecase(taskRepo, stockVarRepo, stockRepo, pubkey)
		assert.NoError(t, err)
		task, err := usecase.RequestTask(context.TODO(), &req)
		assert.Equal(t, err, model.ErrFailedReadGRPCMetadata)
		assert.Nil(t, task)
	})

	t.Run("Error Precondition Token", func(t *testing.T) {
		taskRepo := mock.NewMockTaskRepository(ctrl)
		stockVarRepo := mock.NewMockStockVariantRepository(ctrl)
		stockRepo := mock.NewMockStockRepository(ctrl)

		mmd := metadata.MD{
			"token": []string{},
		}

		mctx := metadata.NewIncomingContext(context.TODO(), mmd)

		usecase, err := NewTaskUsecase(taskRepo, stockVarRepo, stockRepo, pubkey)
		assert.NoError(t, err)
		task, err := usecase.RequestTask(mctx, &req)
		assert.Equal(t, err, model.ErrFailedPrecondition)
		assert.Nil(t, task)
	})

	t.Run("Error Precondition Token Claim", func(t *testing.T) {
		taskRepo := mock.NewMockTaskRepository(ctrl)
		stockVarRepo := mock.NewMockStockVariantRepository(ctrl)
		stockRepo := mock.NewMockStockRepository(ctrl)

		mmd := metadata.MD{
			"token": []string{"Hehe"},
		}

		mctx := metadata.NewIncomingContext(context.TODO(), mmd)

		usecase, err := NewTaskUsecase(taskRepo, stockVarRepo, stockRepo, pubkey)
		assert.NoError(t, err)
		task, err := usecase.RequestTask(mctx, &req)
		assert.Error(t, err)
		assert.Nil(t, task)
	})

}
