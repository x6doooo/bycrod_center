package yahoo

import (
    "bycrod_center/model"
    //"math"
    "bycrod_center/module/mongo"
    "time"
    //"fmt"
    "errors"
    "bycrod_center/module/util"
    "sync"
    //"sort"
    "gopkg.in/mgo.v2/bson"
)

const (
    interval_daily = "1d"
    the_range_daily = "3072d"

    interval_hourly = "1h"
    the_range_hourly = "1024d"

    interval_minutely = "1m"
    the_range_minutely = "100d"

    interval_realtime = "1m"
    the_range_realtime = "5m"
)

type ByHigh []model.QuoteItem
func (me ByHigh) Len() int {
    return len(me)
}
func (me ByHigh) Swap(i, j int) {
    me[i], me[j] = me[j], me[i]
}
func (me ByHigh) Less(i, j int) bool {
    return me[i].High < me[i].High
}

type ByLow []model.QuoteItem
func (me ByLow) Len() int {
    return len(me)
}
func (me ByLow) Swap(i, j int) {
    me[i], me[j] = me[j], me[i]
}
func (me ByLow) Less(i, j int) bool {
    return me[i].Low < me[i].Low
}

func LoadData(dataType string) {

    codes, err := GetCodes()
    if err != nil {
        panic(err)
    }

    startTime := time.Now()
    count := 0
    for {
        codes = loop(codes, startTime, dataType)
        if len(codes) == 0 {
            break
        }
        count += 1
        if count > 10 {
            break
        }
    }

    util.Logger.Info("daily format done!")
}

func loop(codes []string, startTime time.Time, dataType string) (fails []string) {

    var interval string
    var the_range string
    switch dataType {
    case "daily":
        interval = interval_daily
        the_range = the_range_daily
    case "hourly":
        interval = interval_hourly
        the_range = the_range_hourly
    case "minutely":
        interval = interval_minutely
        the_range = the_range_minutely
    case "realtime":
        interval = interval_realtime
        the_range = the_range_realtime
    default:
        interval = interval_daily
        the_range = the_range_daily
    }


    processNum := 2
    step := len(codes) / processNum

    wg := sync.WaitGroup{}
    wg.Add(processNum)

    for i := 0; i < processNum; i++ {
        var codesOfStep []string
        if i == processNum - 1 {
            codesOfStep = codes[i * step : ]
        } else {
            codesOfStep = codes[i * step : (i + 1) * step]
        }
        go func(codes []string) {
            defer wg.Done()
            size := len(codes);
            for idx, code := range codes {

                collectionName := "code_" + code + "_" + dataType

                util.Logger.Info("%s %d/%d %s", code, idx, size, time.Since(startTime).String())
                //fmt.Println(code, idx, "/", size, time.Since(startTime).String())
                //mongo.DB.C(collectionName).DropCollection()

                startTimeOfCurrentCode := time.Now()
                results, err := Get(code, interval, the_range)
                timeUsedOfFetch := time.Now()
                if err != nil {
                    util.Logger.Info("%s failed: %s", err.Error())
                    //fmt.Println(code, "failed")
                    fails = append(fails, code)
                    continue
                }
                timeUsedOfCompute := time.Now()

                var theLastTs int64 = 0

                var dataListHasBeenInserted []model.QuoteItem
                mongo.DB.C(collectionName).Find(bson.M{}).Sort("-ts").All(&dataListHasBeenInserted)
                if len(dataListHasBeenInserted) != 0 {
                    theLastTs = dataListHasBeenInserted[0].Ts
                }

                dataList, err := handle(results, theLastTs)
                if err != nil {
                    continue
                }
                if len(dataList) > 0 {
                    //if dataType == "realtime" {
                    //    mongo.DB.C(collectionName).DropCollection()
                    //}
                    mongo.DB.C(collectionName).Insert(dataList...)
                }
                timeUsedOfInsert := time.Now()

                util.Logger.Info(" - fetch: %s", timeUsedOfFetch.Sub(startTimeOfCurrentCode).String())
                util.Logger.Info(" - compute: %s", timeUsedOfCompute.Sub(timeUsedOfFetch).String())
                util.Logger.Info(" - insert: %s", timeUsedOfInsert.Sub(timeUsedOfCompute).String())
                //fmt.Println(" - fetch:", timeUsedOfFetch.Sub(startTimeOfCurrentCode).String())
                //fmt.Println(" - compute:", timeUsedOfCompute.Sub(timeUsedOfFetch).String())
                //fmt.Println(" - insert:", timeUsedOfInsert.Sub(timeUsedOfCompute).String())
            }

        }(codesOfStep)
    }

    wg.Wait()


    return
}

func InitBaseValue(idx int, ts int64, quotes model.QuoteData) (item model.QuoteItem, err error) {

    item.Ts = ts
    d := time.Unix(ts, 0)
    item.Date = d.UTC().Format("2006-01-02 15:04:05")

    theVolume := quotes.Volume[idx]
    if theVolume != nil {
        item.Volume = *theVolume
    } else {
        err = errors.New("volume error")
        return
    }
    theOpen := quotes.Open[idx]
    if theOpen != nil {
        item.Open = *theOpen
    } else {
        err = errors.New("open error")
        return
    }
    theClose := quotes.Close[idx]
    if theClose != nil {
        item.Close = *theClose
    } else {
        err = errors.New("close error")
        return
    }
    theHigh := quotes.High[idx]
    if theHigh != nil {
        item.High = *theHigh
    } else {
        err = errors.New("high error")
        return
    }
    theLow := quotes.Low[idx]
    if theLow != nil {
        item.Low = *theLow
    } else {
        err = errors.New("low error")
        return
    }
    return
}

//func InitSimpleMovingAverage(sum_set map[string]float64, item *model.QuoteItem, list []model.QuoteItem) (isValid bool) {
//
//    currentSize := len(list)
//    var theLastOne model.QuoteItem
//
//    // init the last record
//    if currentSize > 0 {
//        theLastOne = list[currentSize - 1]
//        item.Close_chg = (item.Close - theLastOne.Close) / theLastOne.Close
//    }
//
//    keys := []int{3, 5, 6, 10, 20}
//    keysStr := []string{"3", "5", "6", "10", "20"}
//    for idx, kStr := range keysStr {
//        if _, ok := sum_set[kStr]; !ok {
//            sum_set[kStr] = 0
//        }
//
//        sum_set[kStr] += item.Close
//
//        kNum := keys[idx]
//        if currentSize >= kNum {
//            maKey := "Ma_" + kStr
//            item.SetFloat64ByFieldName(maKey, sum_set[kStr] / float64(kNum))
//            sum_set[kStr] -= list[currentSize - kNum].Close
//            if currentSize > kNum {
//                v_item, _ := item.GetFloat64ByFieldName(maKey)
//                v_last, _ := theLastOne.GetFloat64ByFieldName(maKey)
//                v_pct := (v_item - v_last) / v_last
//                item.SetFloat64ByFieldName(maKey + "_chg_pct", v_pct)
//            }
//        }
//    }
//
//    if _, ok := sum_set["diff_20"]; !ok {
//        sum_set["diff_20"] = 0
//    }
//    if currentSize >= 20 {
//        item.DiffSquare = math.Pow(item.Close - item.Ma_20, 2)
//        sum_set["diff_20"] += item.DiffSquare
//    }
//
//    if currentSize >= 40 {
//        md := math.Sqrt(sum_set["diff_20"] / 20)
//        sum_set["diff_20"] -= list[currentSize - 20].DiffSquare
//        item.Boll_up_20 = item.Ma_20 + 2 * md
//        item.Boll_dn_20 = item.Ma_20 - 2 * md
//        item.Boll_pct_b = (item.Close - item.Boll_dn_20) / (item.Boll_up_20 - item.Boll_dn_20)
//        item.Boll_bandwidth = (item.Boll_up_20 - item.Boll_dn_20) / item.Ma_20
//        if currentSize > 40 {
//            // 所有指标都有的数据为可用数据
//            isValid = true
//            item.Boll_up_20_chg_pct = (item.Boll_up_20 - theLastOne.Boll_up_20) / theLastOne.Boll_up_20
//            item.Boll_dn_20_chg_pct = (item.Boll_dn_20 - theLastOne.Boll_dn_20) / theLastOne.Boll_dn_20
//            item.Boll_pct_b_chg_pct = (item.Boll_pct_b - theLastOne.Boll_pct_b) / theLastOne.Boll_pct_b
//            item.Boll_bandwidth_chg_pct = (item.Boll_bandwidth - theLastOne.Boll_bandwidth) / theLastOne.Boll_bandwidth
//        }
//    }
//
//    return
//
//}
//
//func InitExponentialMovingAverage(item *model.QuoteItem, list []model.QuoteItem) (isValid bool) {
//    if len(list) == 0 {
//        item.Ema_close_12 = item.Close
//        item.Ema_close_26 = item.Close
//        item.Macd_dem = 0
//        item.Tr = item.High - item.Low
//        return
//    }
//
//    theLastOne := list[len(list) - 1]
//
//
//    // macd
//    var k12 float64 = 2.0 / (12.0 + 1.0)
//    item.Ema_close_12 = (item.Close * k12) + theLastOne.Ema_close_12 * (1 - k12)
//
//    var k26 float64 = 2.0 / (26.0 + 1.0)
//    item.Ema_close_26 = (item.Close * k26) + theLastOne.Ema_close_26 * (1 - k26)
//
//    var k9 float64 = 2.0 / (9.0 + 1.0)
//    ema_diff := item.Ema_close_12 - item.Ema_close_26
//    item.Macd_dem = (ema_diff * k9) + theLastOne.Macd_dem * (1 - k9)
//
//    item.Macd_osc = ema_diff - item.Macd_dem
//
//    // Atr
//    var k14 float64 = 2.0 / (14.0 + 1.0)
//    item.Tr = math.Max(item.High, theLastOne.Close) - math.Min(item.Low, theLastOne.Close)
//    item.Atr_14 = (item.Tr * k14) + theLastOne.Atr_14 * (1 - k14)
//    if theLastOne.Atr_14 != 0 {
//        item.Atr_14_chg = (item.Atr_14 - theLastOne.Atr_14) / theLastOne.Atr_14
//    }
//
//
//    // 达到真正数目之前的值都是不准确的
//    if len(list) < 26 {
//        return
//    }
//
//
//    isValid = true
//    return
//}
//
//func InitRsi(item *model.QuoteItem, list []model.QuoteItem) (isValid bool) {
//    if len(list) == 0 {
//        item.Ema_D_14 = 0
//        item.Ema_U_14 = 0
//        return
//    }
//    theLastOne := list[len(list) - 1]
//
//    var U float64 = 0.0
//    var D float64 = 0.0
//    diff := item.Close - theLastOne.Close
//    if diff > 0 {
//        U = diff
//    } else {
//        D = -diff
//    }
//
//    var k float64 = 2.0 / (14.0 + 1.0)
//    item.Ema_U_14 = U * k + theLastOne.Ema_U_14 * (1 - k)
//    item.Ema_D_14 = D * k + theLastOne.Ema_D_14 * (1 - k)
//
//    rs := item.Ema_U_14 / item.Ema_D_14
//    item.Rsi = 1 - 1 / (1 + rs)
//
//    isValid = true
//    return
//}
//
//func InitBias(item *model.QuoteItem) (isValid bool) {
//    v1 := false
//    v2 := false
//    if item.Ma_3 != 0 && item.Ma_6 != 0 {
//        item.Bias_6 = (item.Close - item.Ma_3) / item.Ma_3
//        item.Bias_3_6 = (item.Ma_3 - item.Ma_6) / item.Ma_6
//        v1 = true
//    }
//    if item.Ma_10 != 0 && item.Ma_20 != 0 {
//        item.Bias_20 = (item.Close - item.Ma_20) / item.Ma_20
//        item.Bias_10_20 = (item.Ma_10 - item.Ma_20) / item.Ma_20
//        v2 = true
//    }
//    isValid = v1 && v2
//    return
//}
//
//func InitRsv(item *model.QuoteItem, list []model.QuoteItem) (isValid bool) {
//    listSize := len(list)
//    if listSize >= 9 {
//        dataSet := list[listSize - 9 : listSize]
//        sort.Sort(ByHigh(dataSet))
//        maxOfHeight := dataSet[8].High
//        sort.Sort(ByLow(dataSet))
//        minOfLow := dataSet[0].Low
//
//        if (maxOfHeight - minOfLow != 0) {
//            item.Rsv_9 = (item.Close - minOfLow) / (maxOfHeight - minOfLow)
//
//            theLastOne := list[listSize - 1]
//
//            alfa :=  1.0 / 3.0
//            item.Rsv_K_9 = alfa * item.Rsv_9 + (1 - alfa) * theLastOne.Rsv_K_9
//            item.Rsv_D_9 = alfa * item.Rsv_K_9 + (1 - alfa) * theLastOne.Rsv_D_9
//
//            isValid = true
//
//        }
//    }
//    return
//}
//

func handle(respData model.RespResult, theLastTs int64) (listInterface []interface{}, err error) {

    resultArr := respData.Chart.Result
    if len(resultArr) == 0 {
        //fmt.Println(respData)
        err = errors.New("results is empty")
        return
    }

    result := respData.Chart.Result[0]
    timestamps := result.Timestamp
    quotes := result.Indicators["quote"][0]

    // init
    //sum_set := map[string]float64{}
    list := make([]model.QuoteItem, 0, len(timestamps))
    for idx, ts := range timestamps {

        if ts <= theLastTs {
            continue
        }
        item, err := InitBaseValue(idx, ts, quotes)
        if err != nil {
            continue
        }

        //// ratio计算
        //validMap := map[string]bool{}
        //validMap["sma"] = InitSimpleMovingAverage(sum_set, &item, list)
        //validMap["ema"] = InitExponentialMovingAverage(&item, list)
        //validMap["rsi"] = InitRsi(&item, list)
        //validMap["rsv"] = InitRsv(&item, list)
        //
        //if validMap["sma"] {
        //    validMap["bias"] = InitBias(&item)
        //}


        list = append(list, item)


        //isValid := true
        //for _, v := range validMap {
        //    if !v {
        //        isValid = false
        //    }
        //}
        //
        //if !isValid {
        //    continue
        //}

        listInterface = append(listInterface, item)

    }

    return
}
