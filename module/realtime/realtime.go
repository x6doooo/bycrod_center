package realtime

import (
    "time"
    "fmt"
    "bycrod_center/module/mongo"
    "bycrod_center/module/yahoo"
    "errors"
    "gopkg.in/mgo.v2/bson"
    "bycrod_center/model"
)

func Do() {
    go processing()
    select{}
}

func processing() {
    // pre handle
    fmt.Println("get codes...")
    codes, err := yahoo.GetCodes()
    if err != nil {
        panic(err)
    }
    if len(codes) == 0 {
        panic("no codes")
    }

    fmt.Println("get history data...")
    sumSets := map[string]map[string]float64{}
    historyDataSet := map[string]([]model.QuoteItem){}
    for _, c := range codes {
        res := []model.QuoteItem{}
        mongo.DB.C("code_" + c + "_daily").Find(bson.M{}).Sort("ts").All(&res)
        if len(res) == 0 {
            continue
        }
        // init
        sum_set := map[string]float64{}
        list := []model.QuoteItem{}
        for _, item := range res {
            yahoo.InitSimpleMovingAverage(sum_set, &item, list)
            list = append(list, item)
        }
        historyDataSet[c] = res[len(res) - 90:]
        sumSets[c] = sum_set
    }

    // go
    fmt.Println("ready... GO!!!")
    for {
        t := time.Now()
        line := t.Format("2006-01-02 15:04:05")
        line = "------- " + line + " -------"
        fmt.Println(line)
        go action(codes, historyDataSet, sumSets)
        time.Sleep(10 * time.Second)
    }
}

func action(codes []string, historyDataSet map[string]([]model.QuoteItem), sumSets map[string]map[string]float64) {
    for _, code := range codes {
        respData, err := yahoo.Get(code, "1m", "5m")
        if err != nil {
            continue
        }
        resultArr := respData.Chart.Result
        if len(resultArr) == 0 {
            fmt.Println(respData)
            err = errors.New("results is empty")
            return
        }

        result := respData.Chart.Result[0]
        timestamps := result.Timestamp
        quotes := result.Indicators["quote"]

        if len(quotes) == 0 {
            continue
        }

        if historyData, ok := historyDataSet[code]; ok {
            if sum_set, ok := sumSets[code]; ok {
                if len(historyData) > 90 {
                    historyDataSet[code] = handle(timestamps, quotes[0], historyData, sum_set)
                }
            }
        }
    }
}

func handle(timestamps []int64, quoteSet model.QuoteData,
historyData []model.QuoteItem, sum_set map[string]float64) ([]model.QuoteItem) {

    for idx, ts := range timestamps {
        theLastOne := historyData[len(historyData) - 1]
        if ts > theLastOne.Ts && quoteSet.Close[idx] != nil {
            item, err := yahoo.InitBaseValue(idx, ts, quoteSet)
            if err == nil {
                continue
            }
            yahoo.InitSimpleMovingAverage(sum_set, &item, historyData)
            historyData = append(historyData, item)
        }
    }

    return historyData

}

