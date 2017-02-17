package main

import (
    //"bycrod_center/module/xueqiu"
    //"fmt"
    //"os"
    //"bycrod_center/module/yahoo"
    //"bycrod_center/module/yahoo/knn"
    //"bycrod_center/module/realtime"
)
import "bycrod_center/server"

func main() {
    //method := os.Args[2]
    //switch method {
    //case "init":
    //    xueqiu.InitStockList()
    //case "daily-format":
    //    yahoo.DailyDataFormat()
    //case "knn-analysis":
    //    knn.DailyDataAnalysis("")
    //case "knn-test":
    //    knn.DailyDataFormatTest()
    //case "realtime":
    //    realtime.Do()
    //default:
    //    fmt.Println("method not found...")
    //}
    server.Start()
}
