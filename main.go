package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/glatika/providence/deliveries/grpc/market"
	"github.com/glatika/providence/deliveries/grpc/market/market_pb"
	"github.com/glatika/providence/deliveries/grpc/stock"
	"github.com/glatika/providence/deliveries/grpc/stock/stock_pb"
	"github.com/glatika/providence/deliveries/panel_rpc"
	"github.com/glatika/providence/repositories"
	usecase "github.com/glatika/providence/usecases"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	Full  string = "full"
	Proxy        = "proxy"
	Panel        = "panel"
)

var (
	port       = flag.Int("port", 42069, "The server port")
	pannelPort = flag.Int("panelPort", 4877, "The panel port")
	mode       = flag.String("mode", Full, "Providence start mode")
	privKey    = flag.String("privateKey", "./priv.pem", "The private key path")
	pubKey     = flag.String("publicKey", "./pub.pem", "The public key path")
	mariaDsn   = flag.String("mariaDsn", "root@tcp(localhost:3306)/providence?parseTime=true", "mariaDb DSN")
)

func main() {
	flag.Parse()

	db, err := gorm.Open(mysql.Open(*mariaDsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	stockRepo := repositories.NewStockRepository(db)
	stockVarRepo := repositories.NewStockVariantRepository(db)
	taskRepo := repositories.NewTaskRepository(db)

	_taskCase, err := usecase.NewTaskUsecase(taskRepo, stockVarRepo, stockRepo, *pubKey)
	if err != nil {
		log.Fatalln("[!!] ERROR public key : ", err)
	}
	_stockCase, err := usecase.NewStockUsecase(stockRepo, stockVarRepo, *privKey)
	if err != nil {
		log.Fatalln("[!!] ERROR private key : ", err)
	}
	_stockVariantCase := usecase.NewStockVariantCase(stockVarRepo)

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	stock_pb.RegisterStockLifecycleServer(grpcServer, stock.StockLifecycle{
		StockCase: _stockCase,
		TaskCase:  _taskCase,
	})

	switch *mode {
	case Full:
		reflection.Register(grpcServer)
		fmt.Printf("Service running on %d", *port)
		pannel := panel_rpc.PanelRPC{
			TaskUsecase:     _taskCase,
			StockUsecase:    _stockCase,
			StockVarUsecase: _stockVariantCase,
		}
		go pannel.Run(*pannelPort)
		grpcServer.Serve(lis)
	case Proxy:
		market_pb.RegisterBarnMarketProviderServer(grpcServer, market.BarnMarketService{
			StockVariantCase: _stockVariantCase,
			StockCase:        _stockCase,
			TaskCase:         _taskCase,
		})

		reflection.Register(grpcServer)
		fmt.Printf("Service running on %d \n", *port)
		grpcServer.Serve(lis)
	case Panel:
		log.Fatalln("This feature unimplemented")
	}
}
