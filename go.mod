module github.com/lifegoeson/blockchain-benchmark

go 1.15

require (
	github.com/Shopify/sarama v1.28.0 // indirect
	github.com/fsouza/go-dockerclient v1.7.2 // indirect
	github.com/go-echarts/go-echarts/v2 v2.2.4
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/hashicorp/go-version v1.3.0 // indirect
	github.com/hyperledger/fabric v1.4.2
	github.com/hyperledger/fabric-amcl v0.0.0-20210319225857-000ace5745f9 // indirect
	github.com/hyperledger/fabric-sdk-go v1.0.0
	github.com/spf13/viper v1.1.1
	github.com/sykesm/zap-logfmt v0.0.4 // indirect
	go.uber.org/zap v1.19.1 // indirect
)

//replace github.com/hyperledger/fabric-sdk-go v1.0.0 => ../fabric-sdk-go
