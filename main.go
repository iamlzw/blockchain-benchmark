package main

import (
	"github.com/lifegoeson/blockchain-benchmark/tools"
	"github.com/spf13/viper"
	"log"
	_ "net/http/pprof"
	_ "time"
)

type lineObj struct {
	tps       int
	minAndSec string
}

func main(){
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
	sdk := tools.InitSDK()
	//初始化context.ChannelProvider
	ccp := tools.InitCCP(sdk)
	//
	//tools.CreateCCLifecycle(sdk)
	//
	//tools.RunInvoke(ccp)
	//
	tools.InvokeChaincode(ccp)
	tools.QueryChaincode(ccp)
	//createCC(sdk)

}

