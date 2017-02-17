package main

import (
    //"bycrodata_grab/module/xueqiu"
    //"fmt"
    //"os"
    //"bycrodata_grab/module/yahoo"
    //"bycrodata_grab/module/yahoo/knn"
    //"bycrodata_grab/module/realtime"
)
import "bycrodata_grab/server"

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
