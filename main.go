package main

import (
	ctx "context"
	"encoding/csv"
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/spf13/viper"
	"log"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
	_ "time"
)

type lineObj struct {
	tps       int
	minAndSec string
}

func main() {
	//开启pprof服务
	//go func() {
	//	_ = http.ListenAndServe("localhost:6060", nil)
	//}()
	//加载性能测试配置文件
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config/")
	viper.SetConfigName("benchmark")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	//初始化sdk
	sdk := initSDK()
	//初始化context.ChannelProvider
	ccp := initCCP(sdk)
	//RunInvoke(ccp)
	//
	invokeChaincode(ccp)
	queryChaincode(ccp)

	//createCC(sdk)
	//createCCLifecycle()

}

func generateTpsChart(items []opts.LineData, name string, title string, subTitle string) {
	// create a new line instance
	line := charts.NewLine()
	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeRomantic}),
		charts.WithTitleOpts(opts.Title{
			Title:    title,
			Subtitle: subTitle,
		}))

	var strs []string
	// Put data into instance
	for _, item := range items {
		strs = append(strs, item.Name)
	}

	fmt.Println(len(strs))
	fmt.Println(len(items))
	line.SetXAxis(strs).
		AddSeries("Category A", items).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	f, _ := os.Create(name)
	_ = line.Render(f)
}

func generateTpsKlineChart(items []opts.KlineData) {
	line := charts.NewKLine()
	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Line example in Westeros theme",
			Subtitle: "Line chart rendered by the http server this time",
		}))

	var strs []string

	// Put data into instance
	for _, item := range items {
		strs = append(strs, item.Name)
	}

	fmt.Println(len(strs))
	fmt.Println(len(items))
	line.SetXAxis(strs).
		AddSeries("Category A", items).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	f, _ := os.Create("invoke.html")
	_ = line.Render(f)
}

func listenBlockEvent(ccp context.ChannelProvider, ch chan int) {
	ec, err := event.New(ccp, event.WithBlockEvents())
	if err != nil {
		_ = fmt.Errorf("init event client error %s", err)
	}

	reg, notifier, err := ec.RegisterBlockEvent()

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

func readTestData() [][]string {
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

func RunInvoke(ccp context.ChannelProvider) {
	//读取测试数据
	records := readTestData()
	var start, end int
	var i int
	//query测试开启的协程的数量
	invokeGoroutineNums := viper.GetInt("invoke.goroutine.count")
	//每个协程执行的测试用例的数量
	invokeTestCaseNums := viper.GetInt("invoke.goroutine.test_case_count")

	ctx3, cancel := ctx.WithTimeout(ctx.Background(), 20*time.Second)

	defer cancel()

	//items1 := make([]opts.LineData, 0)

	var workCounter1 int
	var count int
	var times int

	workChan1 := make(chan int, 200)

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

	for i = 0; i < invokeGoroutineNums; i++ {
		start = i * invokeTestCaseNums
		end = start + invokeTestCaseNums
		subRecords := records[start:end]
		go invoke(ccp, subRecords, ctx3, workChan1)
	}

	time.Sleep(20 * time.Second)
	t1.Stop()
	fmt.Println(count)
	fmt.Println(times)
	fmt.Println(count / times)

	//generateTpsChart(items,"query.html","Hyperledger Fabric Query Transactions Per Second","每秒查询交易数")
	//generateTpsChart(items1,"invoke.html","Hyperledger Fabric Invoke Transactions Per Five Second","每5秒交易数")

	//defer close(workChan)
	defer close(workChan1)
}

func RunQuery(ccp context.ChannelProvider) {
	//读取测试数据
	records := readTestData()
	var start, end int
	var i int
	var workCounter int
	workChan := make(chan int, 200)
	querySyncContext, _ := ctx.WithTimeout(ctx.Background(), 50*time.Second)
	//query测试
	//query测试开启的协程的数量
	queryGoroutineNums := viper.GetInt("query.goroutine.count")
	//每个协程执行的测试用力的数量
	queryTestCaseNums := viper.GetInt("query.goroutine.test_case_count")
	var l sync.Mutex
	//count := 0
	for i = 0; i < queryGoroutineNums; i++ {
		start = i * queryTestCaseNums
		end = start + queryTestCaseNums
		subRecords := records[start:end]
		go query(querySyncContext, ccp, subRecords, workChan)
	}

	//items := make([]opts.LineData, 0)
	t := time.NewTicker(time.Second * 1)

	go func() {
		for {
			select {
			case <-t.C:
				l.Lock()
				fmt.Println("query tps :", workCounter)
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
				fmt.Println("current goroutine num : ", runtime.NumGoroutine())
			}
		}
	}()

	time.Sleep(50 * time.Second)
	t.Stop()

	defer close(workChan)
}

func testpath() {
	fmt.Println(runtime.GOOS)
	fmt.Println(filepath.Join("chaincode/go", "src", "github.com/testchaincode1"))
	fmt.Println(strings.Replace("opt\\go\\src\\github.com\\testchaincode1", "\\", "/", -1))
}
