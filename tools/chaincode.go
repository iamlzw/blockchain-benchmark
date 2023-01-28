package tools

import (
	"fmt"
	mb "github.com/hyperledger/fabric-protos-go/msp"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	_ "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	_ "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	lcpackager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/lifecycle"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/policydsl"
)

func CreateCCLifecycle(sdk *fabsdk.FabricSDK) {

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

	// Approve cc
	approveCC(packageID, orgResMgmt2,peer2)

	// Query approve cc
	queryApprovedCC(orgResMgmt1,peer1)

	// Query approve cc
	queryApprovedCC(orgResMgmt2,peer2)

	// Check commit readiness
	checkCCCommitReadiness(orgResMgmt1,peer1)

	// Check commit readiness
	checkCCCommitReadiness(orgResMgmt2,peer2)

	// Commit cc
	commitCC(orgResMgmt1,peer1)

	// Commit cc
	//commitCC(orgResMgmt2,peer2)

	// Query committed cc
	queryCommittedCC(orgResMgmt1,peer1)
	//
	// Query committed cc
	queryCommittedCC(orgResMgmt2,peer2)

	// Init cc
	initCC(sdk)

}

func packageCC() (string, []byte) {
	desc := &lcpackager.Descriptor{
		Path:  ccPath,
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

	_, err := orgResMgmt.LifecycleApproveCC(channelID, approveCCReq, resmgmt.WithTargetEndpoints(peer), resmgmt.WithOrdererEndpoint(orderer), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
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
	_, err = client.Execute(channel.Request{ChaincodeID: ccID, Fcn: "InitLedger", Args: initArgs, IsInit: true},
		channel.WithRetry(retry.DefaultChannelOpts),channel.WithTargetEndpoints(peer1,peer2))
	if err != nil {
		fmt.Printf("Failed to init: %s", err)
	}
}

