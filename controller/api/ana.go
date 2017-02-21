package api

import (
    "bycrod_center/module/ana"
    "fmt"
    //talib "github.com/markcheno/go-talib"
    "github.com/x6doooo/talib"
)

func init() {
    data := ana.QueryDailyData("code_AMD_daily")

    // sma
    //sma, _, _ := talib.Sma(data.Close, 20)
    ////sma = talib.FormatFloat64(sma, begIdx, nbElement)
    //fmt.Println(sma)
    //res, begIdx, nbElement := talib.Cdl2Crows(data.Open, data.High, data.Low, data.Close)
    //res = talib.FormatInt(res, begIdx, nbElement)

    start := 500
    end := 1000
    open1 := data.Open[start : end]
    close1 := data.Close[start : end]
    high1 := data.High[start : end]
    low1 := data.Low[start : end]
    fmt.Println(data.Open)
    res1, begIdx1, nbElement1 := talib.Cdl2Crows(open1, high1, low1, close1)
    res2, begIdx2, nbElement2 := talib.Cdl2Crows_bbb(open1, high1, low1, close1)


    res1_idx := []int{}
    for idx, v := range res1 {
        if v != 0 {
            res1_idx = append(res1_idx, idx);
        }
    }
    res2_idx := []int{}
    for idx, v := range res2 {
        if int(v) != 0 {
            res2_idx = append(res2_idx, idx);
        }
    }

    fmt.Println(res1_idx, len(res1), begIdx1, nbElement1)
    fmt.Println(res2_idx, len(res2), begIdx2, nbElement2)
    fmt.Println(res1)
    fmt.Println(res2)

    //res, begIdx, nbElement = talib.Cdl3BlackCrows(data.Open, data.High, data.Low, data.Close)
    //res = talib.FormatInt(res, begIdx, nbElement)
    //res, _, _ = talib.Cdl3BlackCrows(data.Open, data.High, data.Low, data.Close)
    //fmt.Println(res)
    //talib.Cr
}
