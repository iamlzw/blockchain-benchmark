package tools

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"time"
)

func listenBlockEvent(ccp context.ChannelProvider,ch chan int){
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
	//timeTickerChan := time.Tick(time.Second * 10)
	//count := 0
	for {
		select {
		case bEvent = <-notifier:
			//count += len(bEvent.Block.Data.Data)
			ch <- len(bEvent.Block.Data.Data)
			//case <-timeTickerChan:
			//	ch<-count
			//	count = 0
			//case <-time.After(200*time.Second):
			//	return
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
