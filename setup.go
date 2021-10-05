package main

import (
	context1 "context"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	_ "io/ioutil"
	"time"

	//"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/policydsl"
)

var defaultInitCCArgs = [][]byte{[]byte("init")}
//var defaultInitCCArgs = [][]byte{[]byte("init")}
const (
	peer1 = "peer0.org1.example.com"
	peer2 = "peer0.org2.example.com"
)

//init the sdk
func initSDK() *fabsdk.FabricSDK {
	//// Initialize the SDK with the configuration file
	configProvider := config.FromFile("config/config_e2e.yaml")
	sdk, err := fabsdk.New(configProvider)
	if err != nil {
		_ = fmt.Errorf("failed to create sdk: %v", err)
	}
	return sdk
}

func initCCP(sdk *fabsdk.FabricSDK) context.ChannelProvider{
	ccp := sdk.ChannelContext("mychannel", fabsdk.WithUser("User1"),fabsdk.WithOrg("Org1"))
	return ccp
}

func createCC(sdk *fabsdk.FabricSDK) {

	ccPkg, err := packager.NewCCPackage("github.com/benchmark", "./chaincode")
	if err != nil {
		fmt.Println(err)
	}
	adminContext := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg("Org1"))

	// Org resource management client
	orgResMgmt, err := resmgmt.New(adminContext)
	if err != nil {
		fmt.Printf("Failed to create new resource management client: %s\n", err)
	}
	// Install example cc to org peers
	installCCReq := resmgmt.InstallCCRequest{Name: "benchmark", Path: "github.com/benchmark", Version: "1.0", Package: ccPkg}
	_, err = orgResMgmt.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		fmt.Println(err)
	}

	adminContextOrg2 := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg("Org2"))

	// Org resource management client
	org2ResMgmt, err := resmgmt.New(adminContextOrg2)
	if err != nil {
		fmt.Printf("Failed to create new resource management client: %s\n", err)
	}
	// Install example cc to org peers
	installCCReqOrg2 := resmgmt.InstallCCRequest{Name: "benchmark", Path: "github.com/benchmark", Version: "1.0", Package: ccPkg}
	_, err = org2ResMgmt.InstallCC(installCCReqOrg2, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		fmt.Println(err)
	}

	// Set up chaincode policy
	ccPolicy := policydsl.SignedByAnyMember([]string{"Org1MSP"})
	// Org resource manager will instantiate 'example_cc' on channel
	resp, err := orgResMgmt.InstantiateCC(
		"mychannel",
		resmgmt.InstantiateCCRequest{Name: "benchmark", Path: "github.com/benchmark", Version: "1.0", Args: defaultInitCCArgs,Policy: ccPolicy},
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
	)

	fmt.Println(resp.TransactionID)
}

func invokeChaincode(ctx context.ChannelProvider){
	cc,err:= channel.New(ctx)
	if err != nil {
		fmt.Errorf("Failed to create new event client: %s", err)
	}

	args := [][]byte{[]byte("invoke"),[]byte("A0001"),[]byte("A5001"),[]byte("1")}
	resp,err := cc.Execute(channel.Request{ChaincodeID: "testchaincode", Fcn: "invoke", Args: args},
		channel.WithRetry(retry.DefaultChannelOpts),

	)
	if err != nil {
		fmt.Println("Failed to query funds: %s", err)
	}
	fmt.Println(resp)
}
func invoke(ccp context.ChannelProvider,records [][]string,ctx context1.Context){
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
			for ; i < l ; i++ {
				args := [][]byte{[]byte(records[i][0]),[]byte(records[i][1]),[]byte("1")}
				_, _ = cc.Execute(channel.Request{ChaincodeID: "benchmark", Fcn: "invoke", Args: args},
					channel.WithRetry(retry.DefaultChannelOpts),
					channel.WithTargetEndpoints(peer1),
				)
			}
	}
}

func query(ccp context.ChannelProvider,records [][]string, ch chan int,stopSignal chan bool){
	//ccp := sdk.ChannelContext("mychannel", fabsdk.WithUser("User1"),fabsdk.WithOrg("Org1"))
	cc,err := channel.New(ccp)
	if err != nil {
		fmt.Println(err)
	}
	l := len(records)
	var i int

	select {
		case <-stopSignal:
			return
		default:
			for ; i < l ; i++ {
				args := [][]byte{[]byte(records[i][0])}
				_,_ = cc.Query(channel.Request{ChaincodeID: "benchmark", Fcn: "query", Args: args},
					channel.WithRetry(retry.DefaultChannelOpts),
					channel.WithTargetEndpoints(peer1),
					channel.WithTimeout(fab.PeerConnection,1*time.Second))
				ch<-1
				//fmt.Println(string(resp.Payload))
			}
	}
}


func listenBlockEvent1(ccp context.ChannelProvider){
	ec,err := event.New(ccp,event.WithBlockEvents())
	if err !=nil {
		_ = fmt.Errorf("init event client error %s", err)
	}

	reg, notifier, err :=ec.RegisterBlockEvent()

	if err != nil {
		fmt.Printf("Failed to register block event: %s", err)
		return
	}
	defer ec.Unregister(reg)

	var bEvent *fab.BlockEvent
	timeTickerChan := time.Tick(time.Second * 2)
	count := 0
	for {
		select {
		case bEvent = <-notifier:
			//fmt.Println(len(bEvent.Block.Data.Data))
			count += len(bEvent.Block.Data.Data)
		case <-timeTickerChan:
			fmt.Println(count/2)
			count = 0
		case <-time.After(200*time.Second):
			return
		}
	}

}
//
//func listenCCEvent(ccp context.ChannelProvider){
//	ec,err := event.New(ccp,event.WithBlockEvents())
//
//	if err !=nil {
//		fmt.Errorf("init event client error %s",err)
//	}
//
//	eventID := "event"
//
//	reg, notifier, err :=ec.RegisterChaincodeEvent("mycc",eventID)
//
//	if err != nil {
//		fmt.Printf("Failed to register block event: %s", err)
//		return
//	}
//	defer ec.Unregister(reg)
//
//	var ccEvent *fab.CCEvent
//	select {
//	case ccEvent = <-notifier:
//		fmt.Printf("receive block event %v",ccEvent)
//	case <-time.After(time.Second * 200):
//		fmt.Printf("Did NOT receive block event\n")
//	}
//}

//type commitTxHandler struct {
//	eventch chan *fab.TxStatusEvent
//}
//
//func newSubmitHandler(eventch chan *fab.TxStatusEvent) invoke.Handler {
//	return invoke.NewSelectAndEndorseHandler(
//		invoke.NewEndorsementValidationHandler(
//			invoke.NewSignatureValidationHandler(&commitTxHandler{eventch}),
//		),
//	)
//}
//
//func listenTxEvent(client *event.Client,chlclient *channel.Client,args [][]byte){
//	//args := [][]byte{[]byte("invoke"),[]byte("A0001"),[]byte("A5001")}
//
//	client, err := event.New(ctx,event.WithBlockEvents())
//	cc,err:= channel.New(ctx)
//	if err != nil {
//		fmt.Errorf("Failed to create new event client: %s", err)
//	}
//
//	var options []channel.RequestOption
//	options = append(options, channel.WithTimeout(fab.Execute, 1000))
//	var notifier chan *fab.TxStatusEvent
//
//	request := channel.Request{ChaincodeID: "testchaincode", Fcn: "invoke"}
//
//	response, err := cc.InvokeHandler(
//		newSubmitHandler(notifier),
//		request,
//		options...,
//	)
//
//
//	//resp,err := cc.Execute(channel.Request{ChaincodeID: "fabcar02", Fcn: "createCar", Args: [][]byte{[]byte("CAR105"), []byte("VM"),[]byte("Polo"),[]byte("Grey"),[]byte("Mary")}},
//	//	channel.WithRetry(retry.DefaultChannelOpts),
//	//)
//	//if err != nil {
//	//	fmt.Println("Failed to query funds: %s", err)
//	//}
//	//
//	fmt.Print(response)
//
//
//	if err != nil {
//		fmt.Println("failed to register tx event")
//	}
//	defer client.Unregister(registration)
//
//	//fmt.Println("tx event registered successfully")
//
//	var txEvent *fab.TxStatusEvent
//	select {
//	case txEvent = <-notifier:
//		fmt.Printf("Received block event: %#v\n", txEvent)
//	case <-time.After(time.Second * 200):
//		fmt.Printf("Did NOT receive tx event\n")
//	}
//}




