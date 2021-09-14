package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func main(){
	//testpath()
	sdk := initSDK()
	ccp := initCCP(sdk)
	cc := initCC()
	//listenBlockEvent(ccp)

	//client, err := event.New(ccp,event.WithBlockEvents())
	//cc,err:= channel.New(ccp)
	//if err != nil {
	//	fmt.Errorf("Failed to create new event client: %s", err)
	//}
	//args := [][]byte{[]byte("A0001"),[]byte("A5001"),[]byte("1")}
	//
	//listenTxEvent(client,cc,args)
	go listenBlockEvent(ccp)

	records := readTestData()
	var start,end int
	var i int
	for i  = 0; i < 100 ; i++{
		start = i*100
		end = start + 100
		subRecords := records[start:end]
		go query(cc,subRecords)
	}

	//createCC(sdk)

	time.Sleep(100000*time.Second)

	//createChannel(sdk)

	//joinChannel(sdk)
	//
	//createCC(sdk)
	//invokeChaincode(ccp)
	//
	//queryLedger(sdk)
	//downloadBlock(sdk)
	//parseBlock()
}


func readTestData() [][]string{
	//准备读取文件
	fileName := "test.csv"
	fs, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("can not open the file, err is %+v", err)
	}
	defer fs.Close()
	fs1, _ := os.Open(fileName)
	r1 := csv.NewReader(fs1)
	content, err := r1.ReadAll()
	if err != nil {
		log.Fatalf("can not readall, err is %+v", err)
	}
	return content
}

func testpath(){
	fmt.Println(runtime.GOOS)
	fmt.Println(filepath.Join("chaincode/go", "src", "github.com/testchaincode1"))
	fmt.Println(strings.Replace("opt\\go\\src\\github.com\\testchaincode1", "\\", "/", -1))
}
