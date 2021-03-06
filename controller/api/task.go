package api

import (
    "github.com/gin-gonic/gin"
    "bycrod_center/module/xueqiu"
    "bycrod_center/module/yahoo"
    "net/http"
    //"bycrod_center/module/yahoo/knn"
    //"bycrod_center/module/yahoo/noname_0"
    "bycrod_center/module/mongo"
    "strings"
    "time"
    "bycrod_center/module/util"
)

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

func TaskClean(ctx *gin.Context) {
    names, _ := mongo.DB.CollectionNames()

    validNames := []string{}
    for _, name := range names {
        //util.Logger.Info("%s %v %v %v", name, strings.Contains(name, "code_"), strings.Contains(name, "_daily"), strings.Contains(name, "_hourly"))
        if strings.Contains(name, "code_") && (strings.Contains(name, "_daily") || strings.Contains(name, "_hourly")||
            strings.Contains(name, "_minutely") || strings.Contains(name, "_realtime")) {
            validNames = append(validNames, name)
            mongo.DB.C(name).DropCollection()
        }
    }
    ctx.JSON(http.StatusOK, gin.H{
        "code": 0,
    })
}

func TaskXueqiuInit(ctx *gin.Context) {
    xueqiu.InitStockList()
    ctx.JSON(http.StatusOK, gin.H{
        "code": 0,
    })
}

func TaskFormat(ctx *gin.Context) {
    t := ctx.Query("type")
    yahoo.LoadData(t)
    ctx.JSON(http.StatusOK, gin.H{
        "code": 0,
    })
}

func TaskRealtime(ctx *gin.Context) {

    go func() {
        for {
            util.Logger.Info("realtime loading...")
            yahoo.LoadData("realtime")
            time.Sleep(10 * time.Second)
        }
    }()

    ctx.JSON(http.StatusOK, gin.H{
        "code": 0,
    })
}

func TaskGetXqEvents(ctx *gin.Context) {
    err := xueqiu.GetEvents();
    if err != nil {
        ctx.JSON(http.StatusOK, gin.H{
            "code": 1,
            "data": err.Error(),
        })
    } else {
        ctx.JSON(http.StatusOK, gin.H{
            "code": 0,
        })
    }
}

//func TaskDataKnnAnalysis(ctx *gin.Context) {
//    t := ctx.Query("type")
//    knn.DataAnalysis("", 0, t)
//    ctx.JSON(http.StatusOK, gin.H{
//        "code": 0,
//    })
//}
//
//func TaskDataFormatTest(ctx *gin.Context) {
//    t := ctx.Query("type")
//    knn.DataFormatTest(t)
//    ctx.JSON(http.StatusOK, gin.H{
//        "code": 0,
//    })
//}
//
//func TaskDataAnalysis_Noname_0(ctx *gin.Context) {
//    noname_0.Do()
//    ctx.JSON(http.StatusOK, gin.H{
//        "code": 0,
//    })
//}