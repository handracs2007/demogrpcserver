protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=../rpc/ --go-grpc_opt=paths=source_relative \
    demo.proto