package main

import (
	context1 "context"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"

	//"github.com/hyperledger/fabric-sdk-go/pkg/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	_ "strconv"
	_ "testing"
	"time"

	mb "github.com/hyperledger/fabric-protos-go/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	_ "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"

	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/policydsl"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"

	//mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	//
	//"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	//packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	lcpackager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/lifecycle"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

//var defaultInitCCArgs = [][]byte{[]byte("init"),[]byte("a"),[]byte("200"),[]byte("b"),[]byte("200")}
var defaultInitCCArgs = [][]byte{[]byte("init")}
var initArgs = [][]byte{[]byte("a"), []byte("100"), []byte("b"), []byte("200")}
var ccID="example_cc14"
const (
	peer1 = "peer0.org1.example.com"
	peer2 = "peer0.org2.example.com"
	peer3 = "peer0.org3.example.com"
	peer4 = "peer0.org4.example.com"
	channelID      = "mychannel"
	org1Name        = "Org1"
	org2Name        = "Org2"
	orgAdmin       = "Admin"
	ordererOrgName = "OrdererOrg"
)


//init the sdk
func initSDK() *fabsdk.FabricSDK {
	//// Initialize the SDK with the configuration file
	configProvider := config.FromFile("D:\\workspace\\go\\src\\github.com\\lifegoeson\\blockchain-benchmark\\config\\config_e2e.yaml")
	sdk, err := fabsdk.New(configProvider)
	if err != nil {
		_ = fmt.Errorf("failed to create sdk: %v", err)
	}
	return sdk
}

func initCCP(sdk *fabsdk.FabricSDK) context.ChannelProvider{
	ccp := sdk.ChannelContext("mychannel", fabsdk.WithUser("Admin"),fabsdk.WithOrg(org1Name))
	return ccp
}

func createCC(sdk *fabsdk.FabricSDK) {
	//prepare context
	createCCLifecycle(sdk)

	//ccPkg, err := packager.NewCCPackage("github.com/example_cc", "D:\\workspace\\go")
	//if err != nil {
	//	fmt.Println(err)
	//}
	////在将链码打包后将其tar.gz保存到本地
	//err = ioutil.WriteFile("cc.tar.gz",ccPkg.Code,0644)
	//if err != nil{
	//	fmt.Println("write to file failure:",err)
	//}
	//// Install example cc to org peers
	//installCCReq := resmgmt.InstallCCRequest{Name: "example_cc", Path: "github.com/example_cc", Version: "0", Package: ccPkg}
	//_, err = orgResMgmt.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	//if err != nil {
	//	fmt.Println(err)
	//}
	//// Set up chaincode policy
	//ccPolicy := cauthdsl.SignedByAnyMember([]string{"Org1MSP"})
	//// Org resource manager will instantiate 'example_cc' on channel
	//resp, err := orgResMgmt.InstantiateCC(
	//	channelID,
	//	resmgmt.InstantiateCCRequest{Name: "example_cc", Path: "github.com/example_cc", Version: "0", Args: initArgs, Policy: ccPolicy},
	//	resmgmt.WithRetry(retry.DefaultResMgmtOpts),
	//)
	//
	//fmt.Println(resp.TransactionID)
}

func invokeChaincode(ctx context.ChannelProvider){
	cc,err:= channel.New(ctx)
	if err != nil {
		fmt.Printf("Failed to create new event client: %s", err)
	}

	args := [][]byte{[]byte("move"),[]byte("a"),[]byte("b"),[]byte("10")}
	resp,err := cc.Execute(channel.Request{ChaincodeID: ccID, Fcn: "invoke", Args: args},
		channel.WithRetry(retry.DefaultChannelOpts),
		channel.WithTargetEndpoints(peer1,peer2),

	)
	if err != nil {
		fmt.Printf("Failed to query funds: %s\n", err)
	}
	fmt.Println(string(resp.Payload))
}

func queryChaincode(ctx context.ChannelProvider){
	cc,err:= channel.New(ctx)
	if err != nil {
		fmt.Printf("Failed to create new event client: %s\n", err)
	}

	args := [][]byte{[]byte("query"),[]byte("a")}
	resp,err := cc.Execute(channel.Request{ChaincodeID: ccID, Fcn: "invoke", Args: args},
		channel.WithRetry(retry.DefaultChannelOpts),
		channel.WithTargetEndpoints(peer1,peer2),

	)
	if err != nil {
		fmt.Printf("Failed to query funds: %s\n\n", err)
	}
	fmt.Println(string(resp.Payload))
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
			resp, _ := cc.Execute(channel.Request{ChaincodeID: "benchmark1", Fcn: "invoke", Args: args},
				channel.WithRetry(retry.DefaultChannelOpts),
				channel.WithTargetEndpoints(peer1),
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
//		fmt.Println("Failed to register block event: %s", err)
//		return
//	}
//	defer ec.Unregister(reg)
//
//	var ccEvent *fab.CCEvent
//	select {
//	case ccEvent = <-notifier:
//		fmt.Println("receive block event %v",ccEvent)
//	case <-time.After(time.Second * 200):
//		fmt.Println("Did NOT receive block event\n")
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
//		fmt.Println("Received block event: %#v\n", txEvent)
//	case <-time.After(time.Second * 200):
//		fmt.Println("Did NOT receive tx event\n")
//	}
//}

func createCCLifecycle(sdk *fabsdk.FabricSDK) {

	adminContext1 := sdk.Context(fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(org1Name))

	// Org resource management client
	orgResMgmt1, err := resmgmt.New(adminContext1)
	if err != nil {
		fmt.Printf("Failed to create new resource management client: %s\n", err)
	}

	adminContext2 := sdk.Context(fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(org2Name))

	// Org resource management client
	orgResMgmt2, err := resmgmt.New(adminContext2)
	if err != nil {
		fmt.Printf("Failed to create new resource management client: %s\n\n", err)
	}
	// Package cc
	label, ccPkg := packageCC()
	packageID := lcpackager.ComputePackageID(label, ccPkg)

	// Install cc
	installCC(label, ccPkg, orgResMgmt1,peer1)

	// Install cc
	installCC(label, ccPkg, orgResMgmt2,peer2)

	// Get installed cc package
	getInstalledCCPackage( packageID, ccPkg, orgResMgmt1,peer1)

	// Get installed cc package
	getInstalledCCPackage( packageID, ccPkg, orgResMgmt2,peer2)

	// Query installed cc
	queryInstalled(label, packageID, orgResMgmt1,peer1)

	// Query installed cc
	queryInstalled(label, packageID, orgResMgmt2,peer2)

	// Approve cc
	approveCC(packageID, orgResMgmt1,peer1)
	//
	//// Approve cc
	approveCC(packageID, orgResMgmt2,peer2)
	//
	//// Query approve cc
	queryApprovedCC(orgResMgmt1,peer1)

	// Query approve cc
	queryApprovedCC(orgResMgmt2,peer2)
	////
	////// Check commit readiness
	checkCCCommitReadiness(orgResMgmt1,peer1)

	// Check commit readiness
	checkCCCommitReadiness(orgResMgmt2,peer2)
	////
	////// Commit cc
	commitCC(orgResMgmt1,peer1)
	//////
	//////// Commit cc
	//////commitCC(orgResMgmt2,peer2)
	//////
	//////// Query committed cc
	queryCommittedCC(orgResMgmt1,peer1)
	//
	//// Query committed cc
	queryCommittedCC(orgResMgmt2,peer2)
	//////
	//////// Init cc
	initCC(sdk)

}

func packageCC() (string, []byte) {
	desc := &lcpackager.Descriptor{
		Path:  "D:\\workspace\\go\\src\\github.com\\example_cc",
		Type:  pb.ChaincodeSpec_GOLANG,
		Label: ccID,
	}
	ccPkg, err := lcpackager.NewCCPackage(desc)
	if err != nil {
		fmt.Println(err)
	}
	return desc.Label, ccPkg
}

func installCC(label string, ccPkg []byte, orgResMgmt *resmgmt.Client,peer string) {
	installCCReq := resmgmt.LifecycleInstallCCRequest{
		Label:   label,
		Package: ccPkg,
	}

	resp := lcpackager.ComputePackageID(installCCReq.Label, installCCReq.Package)

	_, err := orgResMgmt.LifecycleInstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts),resmgmt.WithTargetEndpoints(peer))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp)
}

func getInstalledCCPackage(packageID string, ccPkg []byte, orgResMgmt *resmgmt.Client,peer string) {
	resp, err := orgResMgmt.LifecycleGetInstalledCCPackage(packageID, resmgmt.WithTargetEndpoints(peer), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
}

func queryInstalled(label string, packageID string, orgResMgmt *resmgmt.Client,peer string) {
	_, err := orgResMgmt.LifecycleQueryInstalledCC(resmgmt.WithTargetEndpoints(peer), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		fmt.Println(err)
	}
}

func approveCC(packageID string, orgResMgmt *resmgmt.Client, peer string) {
	ccPolicy := policydsl.SignedByNOutOfGivenRole(2,mb.MSPRole_MEMBER,[]string{"Org1MSP","Org2MSP"})
	approveCCReq := resmgmt.LifecycleApproveCCRequest{
		Name:              ccID,
		Version:           "0",
		PackageID:         packageID,
		Sequence:          1,
		EndorsementPlugin: "escc",
		ValidationPlugin:  "vscc",
		SignaturePolicy:   ccPolicy,
		InitRequired:      true,
	}

	_, err := orgResMgmt.LifecycleApproveCC(channelID, approveCCReq, resmgmt.WithTargetEndpoints(peer), resmgmt.WithOrdererEndpoint("orderer.example.com"), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		fmt.Println(err)
	}
}

func queryApprovedCC(orgResMgmt *resmgmt.Client, peer string) {
	queryApprovedCCReq := resmgmt.LifecycleQueryApprovedCCRequest{
		Name:     ccID,
		Sequence: 1,
	}
	_, err := orgResMgmt.LifecycleQueryApprovedCC(channelID, queryApprovedCCReq, resmgmt.WithTargetEndpoints(peer), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		fmt.Println(err)
	}
}

func checkCCCommitReadiness(orgResMgmt *resmgmt.Client,peer string) {
	ccPolicy := policydsl.SignedByNOutOfGivenRole(2,mb.MSPRole_MEMBER,[]string{"Org1MSP","Org2MSP"})
	req := resmgmt.LifecycleCheckCCCommitReadinessRequest{
		Name:              ccID,
		Version:           "0",
		EndorsementPlugin: "escc",
		ValidationPlugin:  "vscc",
		SignaturePolicy:   ccPolicy,
		Sequence:          1,
		InitRequired:      true,
	}
	resp, err := orgResMgmt.LifecycleCheckCCCommitReadiness(channelID, req, resmgmt.WithTargetEndpoints(peer1,peer2), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(resp.Approvals)
}

func commitCC(orgResMgmt *resmgmt.Client, peer string) {
	ccPolicy := policydsl.SignedByNOutOfGivenRole(2,mb.MSPRole_MEMBER,[]string{"Org1MSP","Org2MSP"})
	req := resmgmt.LifecycleCommitCCRequest{
		Name:              ccID,
		Version:           "0",
		Sequence:          1,
		EndorsementPlugin: "escc",
		ValidationPlugin:  "vscc",
		SignaturePolicy:   ccPolicy,
		InitRequired:      true,
	}
	resp, err := orgResMgmt.LifecycleCommitCC(channelID, req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithTargetEndpoints(peer1,peer2), resmgmt.WithOrdererEndpoint("orderer.example.com"))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
}

func queryCommittedCC(orgResMgmt *resmgmt.Client, peer string) {
	req := resmgmt.LifecycleQueryCommittedCCRequest{
		Name: ccID,
	}
	_, err := orgResMgmt.LifecycleQueryCommittedCC(channelID, req, resmgmt.WithTargetEndpoints(peer), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		fmt.Println(err)
	}
}

func initCC(sdk *fabsdk.FabricSDK) {
	//prepare channel client context using client context
	clientChannelContext := sdk.ChannelContext(channelID, fabsdk.WithUser("Admin"), fabsdk.WithOrg(org1Name))
	// Channel client is used to query and execute transactions (Org1 is default org)
	client, err := channel.New(clientChannelContext)
	if err != nil {
		fmt.Printf("Failed to create new channel client: %s", err)
	}

	// init
	_, err = client.Execute(channel.Request{ChaincodeID: ccID, Fcn: "init", Args: initArgs, IsInit: true},
		channel.WithRetry(retry.DefaultChannelOpts),channel.WithTargetEndpoints(peer1,peer2))
	if err != nil {
		fmt.Printf("Failed to init: %s", err)
	}

}




