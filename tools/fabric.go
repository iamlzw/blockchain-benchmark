package tools

import (
	context1 "context"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

//init the sdk
func InitSDK() *fabsdk.FabricSDK {
	//// Initialize the SDK with the configuration file
	configProvider := config.FromFile(configPath)
	sdk, err := fabsdk.New(configProvider)
	if err != nil {
		_ = fmt.Errorf("failed to create sdk: %v", err)
	}
	return sdk
}

func InitCCP(sdk *fabsdk.FabricSDK) context.ChannelProvider{
	ccp := sdk.ChannelContext(channelID, fabsdk.WithUser(orgAdmin),fabsdk.WithOrg(org1Name))
	return ccp
}

func invoke(ccp context.ChannelProvider,records [][]string,ctx context1.Context,c chan int){
	//ccp := sdk.ChannelContext("mychannel", fabsdk.WithUser("User1"),fabsdk.WithOrg("Org1"))
	cc,err := channel.New(ccp)
	if err != nil {
		fmt.Println(err)
	}
	l := len(records)
	var i int
	select {
	case <-ctx.Done():
		return
	default:
		for i = 0; i < l ; i++ {
			args := [][]byte{[]byte(records[i][0]),[]byte(records[i][1]),[]byte("1")}
			resp, _ := cc.Execute(channel.Request{ChaincodeID: "benchmark4", Fcn: "Transfer", Args: args},
				channel.WithRetry(retry.DefaultChannelOpts),
				channel.WithTargetEndpoints(peer1,peer2),
			)

			if resp.TxValidationCode.String() == "VALID" {
				c<-1
			}
		}
	}
}

func query(ctx context1.Context,ccp context.ChannelProvider,records [][]string, ch chan int){
	cc,err := channel.New(ccp)
	if err != nil {
		fmt.Println(err)
	}
	l := len(records)
	var i int

	select {
	case <-ctx.Done():
		return
	default:
		for ; i < l ; i++ {
			args := [][]byte{[]byte(records[i][0])}
			_,_ = cc.Query(channel.Request{ChaincodeID: "benchmark", Fcn: "query", Args: args},
				channel.WithRetry(retry.DefaultChannelOpts),
				channel.WithTargetEndpoints(peer1))
			ch<-1
		}
	}
}

func InvokeChaincode(ctx context.ChannelProvider){
	cc,err:= channel.New(ctx)
	if err != nil {
		fmt.Printf("Failed to create new event client: %s", err)
	}

	args := [][]byte{[]byte("A0000"),[]byte("A0001"),[]byte("10")}
	resp,err := cc.Execute(channel.Request{ChaincodeID: ccID, Fcn: "Transfer", Args: args},
		channel.WithRetry(retry.DefaultChannelOpts),
		channel.WithTargetEndpoints(peer1,peer2),

	)
	if err != nil {
		fmt.Printf("Failed to query funds: %s\n", err)
	}
	fmt.Println(string(resp.Payload))
}

func QueryChaincode(ctx context.ChannelProvider){
	cc,err:= channel.New(ctx)
	if err != nil {
		fmt.Printf("Failed to create new event client: %s\n", err)
	}

	args := [][]byte{[]byte("A0000")}
	resp,err := cc.Execute(channel.Request{ChaincodeID: ccID, Fcn: "Query", Args: args},
		channel.WithRetry(retry.DefaultChannelOpts),
		channel.WithTargetEndpoints(peer1,peer2),

	)
	if err != nil {
		fmt.Printf("Failed to query funds: %s\n\n", err)
	}
	fmt.Println(string(resp.Payload))
}
