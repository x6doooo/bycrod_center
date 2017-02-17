package knn

//import (
//    "bycrodata_grab/module/mongo"
//    "strings"
//    "runtime"
//    "sync"
//    "gopkg.in/mgo.v2/bson"
//    "bycrodata_grab/model"
//    "math"
//    "sort"
//    "time"
//    //"fmt"
//    "bycrodata_grab/module/util"
//)
//
//const (
//    /*
//        1: 0.566
//        3: 0.566
//        7: 0.566
//        9: 0.5
//     */
//    K = 2
//    forecast_daily_collection = "forecast_daily"
//    forecast_hourly_collection = "forecast_hourly"
//)
//
//type forecastResultItem struct {
//    Id        bson.ObjectId `bson:"_id,omitempty"`
//    Name      string `bson:"name"`
//    //Close_chg float64 `bson:"close_chg"`
//    Total_chg           float64
//    Next_close_chg      float64
//    Next_next_close_chg float64
//}
//
//type diffItem struct {
//    Diff                float64
//    Total_chg           float64
//    Next_close_chg      float64
//    Next_next_close_chg float64
//}
//
//type ByDiff []diffItem
//
//func (me ByDiff) Len() int {
//    return len(me)
//}
//func (me ByDiff) Swap(i, j int) {
//    me[i], me[j] = me[j], me[i]
//}
//func (me ByDiff) Less(i, j int) bool {
//    return me[i].Diff < me[j].Diff
//}
//
//func DataAnalysis(date string, ts int64, dataType string) {
//    var cn string
//    switch dataType {
//    case "daily":
//        cn = forecast_daily_collection
//    case "hourly":
//        cn = forecast_hourly_collection
//    }
//
//    names, _ := mongo.DB.CollectionNames()
//
//    mongo.DB.C(cn).DropCollection()
//
//    validNames := []string{}
//    for _, name := range names {
//        if strings.Contains(name, "code_") && strings.Contains(name, "_" + dataType) {
//            validNames = append(validNames, name)
//        }
//    }
//
//    cpuNum := runtime.NumCPU()
//    totalSize := len(validNames)
//
//    step := totalSize / cpuNum
//
//    wg := sync.WaitGroup{}
//    wg.Add(cpuNum)
//
//    startTime := time.Now()
//
//    for i := 0; i < cpuNum; i++ {
//        if i == cpuNum - 1 {
//            names = validNames[i * step :]
//        } else {
//            names = validNames[i * step : (i + 1) * step]
//        }
//
//        go func(names []string) {
//            defer wg.Done()
//            results := []interface{}{}
//            for nameIdx, name := range names {
//
//                util.Logger.Info("%s %d/%d %s", name, nameIdx, len(names), time.Since(startTime).String())
//                //fmt.Println(name, nameIdx, "/", len(names), time.Since(startTime).String())
//                condition := bson.M{}
//                if date != "" {
//                    condition["date"] = bson.M{
//                        "$lte": date,
//                    }
//                }
//                if ts != 0 {
//                    condition["ts"] = bson.M{
//                        "$lte": ts,
//                    }
//                }
//
//                var dataSet []model.QuoteItem
//                mongo.DB.C(name).Find(condition).Sort("ts").All(&dataSet)
//                dataSetSize := len(dataSet)
//
//                if dataSetSize == 0 {
//                    continue
//                }
//
//                theLabelItem := dataSet[dataSetSize - 1]
//                dataSet = dataSet[0:dataSetSize]
//
//                diffSet := []diffItem{}
//                for idx, item := range dataSet {
//                    if idx >= len(dataSet) - 2 {
//                        continue
//                    }
//                    diff := float64(0)
//                    //diff += math.Pow(theLabelItem.Ma_5_chg_pct - item.Ma_5_chg_pct, 2)
//                    //diff += math.Pow(theLabelItem.Ma_10_chg_pct - item.Ma_10_chg_pct, 2)
//                    //diff += math.Pow(theLabelItem.Ma_20_chg_pct - item.Ma_20_chg_pct, 2)
//                    diff += math.Pow(theLabelItem.Boll_up_20_chg_pct - item.Boll_up_20_chg_pct, 2)
//                    diff += math.Pow(theLabelItem.Boll_dn_20_chg_pct - item.Boll_dn_20_chg_pct, 2)
//                    diff += math.Pow(theLabelItem.Boll_pct_b_chg_pct - item.Boll_pct_b_chg_pct, 2)
//                    diff += math.Pow(theLabelItem.Macd_dem - item.Macd_dem, 2)
//                    diff += math.Pow(theLabelItem.Macd_osc - item.Macd_osc, 2)
//                    diff += math.Pow(theLabelItem.Rsi - item.Rsi, 2)
//                    diff += math.Pow(theLabelItem.Rsv_K_9 - item.Rsv_K_9, 2)
//                    diff += math.Pow(theLabelItem.Rsv_D_9 - item.Rsv_D_9, 2)
//                    diff += math.Pow(theLabelItem.Atr_14_chg - item.Atr_14_chg, 2)
//                    diff += math.Pow(theLabelItem.Atr_14 - item.Atr_14, 2)
//                    diff += math.Pow(theLabelItem.Tr - item.Tr, 2)
//                    diff += math.Pow(theLabelItem.Bias_6 - item.Bias_6, 2)
//                    diff += math.Pow(theLabelItem.Bias_20 - item.Bias_20, 2)
//                    diff += math.Pow(theLabelItem.Close_chg - item.Close_chg, 2)
//
//                    // ma5 over ma10
//                    a := (theLabelItem.Ma_5 - theLabelItem.Ma_10) / theLabelItem.Ma_10
//                    b := (item.Ma_5 - item.Ma_10) / item.Ma_10
//                    diff += math.Pow(a - b, 2)
//                    // ma10 over ma20
//                    a = (theLabelItem.Ma_10 - theLabelItem.Ma_20) / theLabelItem.Ma_20
//                    b = (item.Ma_10 - item.Ma_20) / item.Ma_20
//                    diff += math.Pow(a - b, 2)
//
//                    //
//                    chg1 := dataSet[idx + 1].Close_chg
//                    chg2 := dataSet[idx + 2].Close_chg
//                    diffSet = append(diffSet, diffItem{
//                        Diff: diff,
//                        Total_chg: (1 + chg1) * (1 + chg2) - 1,
//                        Next_close_chg: chg1,
//                        Next_next_close_chg: chg2,
//                        //Next_close_chg: dataSet[idx + 1].Close_chg,
//                        //Next_next_close_chg: dataSet[idx + 2].Close_chg,
//                    })
//                }
//
//                // sort
//                sort.Sort(ByDiff(diffSet))
//                if (len(diffSet) > K + 1) {
//                    diffSet = diffSet[0:K]
//                }
//
//                next_chg := float64(0)
//                next_next_chg := float64(0)
//                for _, item := range diffSet {
//                    next_chg += item.Next_close_chg / K
//                    next_next_chg += item.Next_next_close_chg / K
//                }
//                results = append(results, forecastResultItem{
//                    Name: name,
//                    Next_close_chg: next_chg,
//                    Next_next_close_chg: next_next_chg,
//                    Total_chg: (1 + next_chg) * (1 + next_next_chg) - 1,
//                })
//            }
//
//            //cn := forecast_daily_collection
//
//
//            err := mongo.DB.C(cn).Insert(results...)
//            if err != nil {
//                util.Logger.Info("%s", err.Error())
//                //fmt.Println(err)
//            }
//        }(names)
//    }
//
//    wg.Wait()
//
//    util.Logger.Info("%s kNN analysis done!", dataType)
//    //fmt.Println("done")
//
//}
