override MODE := "full"
override DBDSN := "root@tcp(localhost:3306)/providence?parseTime=true"
override PRIVKEYPATH := "./test_pkcs8.pem"
override PUBKEYPATH := "./test_pub.pem"

grpc-gen:
	protoc --go_out=. --go_opt=paths=source_relative \
    	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		-I=. \
    	deliveries/grpc/market/market_pb/service.proto && \
	protoc --go_out=. --go_opt=paths=source_relative \
    	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
    	deliveries/grpc/stock/stock_pb/service.proto

up-migrate:
	sql-migrate up 

down-migrate:
	sql-migrate down

model/mock/mock_stock_usecase.go:
	mockgen -destination=model/mock/mock_stock_usecase.go -package=mock github.com/glatika/providence/model StockUsecase 
	
model/mock/mock_stock_variant_usecase.go:
	mockgen -destination=model/mock/mock_stock_variant_usecase.go -package=mock github.com/glatika/providence/model StockVariantUsecase 

model/mock/mock_task_usecase.go:
	mockgen -destination=model/mock/mock_task_usecase.go -package=mock github.com/glatika/providence/model TaskUsecase 

model/mock/mock_stock_repository.go:
	mockgen -destination=model/mock/mock_stock_repository.go -package=mock github.com/glatika/providence/model StockRepository 

model/mock/mock_stock_variant_repository.go:
	mockgen -destination=model/mock/mock_stock_variant_repository.go -package=mock github.com/glatika/providence/model StockVariantRepository

model/mock/mock_task_repository.go:
	mockgen -destination=model/mock/mock_task_repository.go -package=mock github.com/glatika/providence/model TaskRepository

mockgen: model/mock/mock_stock_usecase.go \
	model/mock/mock_stock_variant_usecase.go \
	model/mock/mock_task_usecase.go \
	model/mock/mock_stock_repository.go \
	model/mock/mock_stock_variant_repository.go \
	model/mock/mock_task_repository.go

clean:
	rm -rf model/mock/*

keypair:
	openssl genpkey -algorithm EC -out test_base.pem -pkeyopt ec_paramgen_curve:P-384 && \
	openssl pkcs8 -topk8  -in test_base.pem -nocrypt -out test_pkcs8.pem && \\
	openssl ec -pubout -in test_base.pem -out test_pub.pem

only-test:
	SVC_DISABLE_CACHING=true PRIVKEYPATH=$(PRIVKEYPATH) PUBKEYPATH=$(PUBKEYPATH) richgo test ./... -v --cover

run-server:
	go run main.go --privateKey ./test_pkcs8.pem --publicKey ./test_pub.pem \
	--mariaDsn $(DBDSN) --mode $(MODE)

run-panel: 
	cd panel && mint start