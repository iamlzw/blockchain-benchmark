package tools

import (
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"os"
)

func generateTpsChart(items []opts.LineData,name string,title string,subTitle string){
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
	for _,item := range items {
		strs = append(strs,item.Name)
	}

	fmt.Println(len(strs))
	fmt.Println(len(items))
	line.SetXAxis(strs).
		AddSeries("Category A", items).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	f, _ := os.Create(name)
	_ = line.Render(f)
}

// 生成tps k线图
func generateTpsKlineChart(items []opts.KlineData){
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
	for _,item := range items {
		strs = append(strs,item.Name)
	}

	fmt.Println(len(strs))
	fmt.Println(len(items))
	line.SetXAxis(strs).
		AddSeries("Category A", items).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	f, _ := os.Create("invoke.html")
	_ = line.Render(f)
}
