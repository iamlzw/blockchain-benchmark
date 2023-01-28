package tools

import (
	ctx "context"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/spf13/viper"
	"runtime"
	"sync"
	"time"
)

func RunInvoke(ccp context.ChannelProvider){
	//读取测试数据
	records := ReadTestData()
	var start,end int
	var i int
	//query测试开启的协程的数量
	invokeGoroutineNums := viper.GetInt("invoke.goroutine.count")
	//每个协程执行的测试用例的数量
	invokeTestCaseNums := viper.GetInt("invoke.goroutine.test_case_count")

	ctx3, cancel := ctx.WithTimeout(ctx.Background(), 10*time.Second)

	defer cancel()

	//items1 := make([]opts.LineData, 0)

	var workCounter1 int

	var count int

	var times int

	workChan1 := make(chan int,200)

	//go listenBlockEvent(ccp,workChan1)

	var l sync.Mutex

	t1 := time.NewTicker(time.Second * 1)

	//count := 0
	go func() {
		for {
			select {
			case <-t1.C:
				l.Lock()
				//fmt.Println(workCounter1)
				count += workCounter1
				times++
				workCounter1 = 0
				//items1 = append(items1, opts.LineData{Value: workCounter1,Name: strconv.Itoa(time.Now().Minute())+":"+strconv.Itoa(time.Now().Second())})
				l.Unlock()
			case <-workChan1:
				l.Lock()
				workCounter1 += 1
				l.Unlock()
			}
		}
	}()

	for i  = 0; i < invokeGoroutineNums ; i++{
		start = i * invokeTestCaseNums
		end = start + invokeTestCaseNums
		subRecords := records[start:end]
		go invoke(ccp,subRecords,ctx3,workChan1)
	}

	time.Sleep(10* time.Second)
	t1.Stop()
	fmt.Println(count)
	fmt.Println(times)
	fmt.Println(count/times)

	//generateTpsChart(items,"query.html","Hyperledger Fabric Query Transactions Per Second","每秒查询交易数")
	//generateTpsChart(items1,"invoke.html","Hyperledger Fabric Invoke Transactions Per Five Second","每5秒交易数")

	//defer close(workChan)
	defer close(workChan1)
}

func RunQuery(ccp context.ChannelProvider){
	//读取测试数据
	records := ReadTestData()
	var start,end int
	var i int
	var workCounter int
	workChan := make(chan int, 200)
	querySyncContext, _ := ctx.WithTimeout(ctx.Background(),50*time.Second)
	//query测试
	//query测试开启的协程的数量
	queryGoroutineNums := viper.GetInt("query.goroutine.count")
	//每个协程执行的测试用力的数量
	queryTestCaseNums := viper.GetInt("query.goroutine.test_case_count")
	var l sync.Mutex
	//count := 0
	for i  = 0; i < queryGoroutineNums; i++{
		start = i * queryTestCaseNums
		end = start + queryTestCaseNums
		subRecords := records[start:end]
		go query(querySyncContext,ccp,subRecords, workChan)
	}


	//items := make([]opts.LineData, 0)
	t := time.NewTicker(time.Second * 1)

	go func() {
		for {
			select {
			case <-t.C:
				l.Lock()
				fmt.Println("query tps :",workCounter)
				workCounter = 0
				//items = append(items, opts.LineData{Value: workCounter,Name: strconv.Itoa(time.Now().Minute())+":"+strconv.Itoa(time.Now().Second())})
				l.Unlock()
				//l.Lock()
				//
				//l.Unlock()
			case <-workChan:
				l.Lock()
				workCounter += 1
				l.Unlock()
			default:
				break
			}
		}
	}()

	//打印当前的协程数量,
	timeTickerChan2 := time.Tick(time.Second * 1)
	go func() {
		for {
			select {
			case <-timeTickerChan2:
				fmt.Println("current goroutine num : ",runtime.NumGoroutine())
			}
		}
	}()

	time.Sleep(50*time.Second)
	t.Stop()

	defer close(workChan)
}
