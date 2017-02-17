package noname_0

//import (
//    "bycrodata_grab/module/mongo"
//    "strings"
//    "gopkg.in/mgo.v2/bson"
//    "bycrodata_grab/model"
//    "os"
//    "encoding/csv"
//    "bycrodata_grab/module/util"
//    "strconv"
//)
//
////func Do() {
////
////}
//
//func Do() {
//
//    train_file, err := os.Create("data_daily_train.csv")
//    test_file, err := os.Create("data_daily_test.csv")
//
//    if err != nil {
//        util.Logger.Info("create csv file error: %v", err)
//        return
//    }
//    defer train_file.Close()
//    defer test_file.Close()
//
//    train_writer := csv.NewWriter(train_file)
//    defer train_writer.Flush()
//    test_writer := csv.NewWriter(test_file)
//    defer test_writer.Flush()
//
//    names, _ := mongo.DB.CollectionNames()
//
//    for idx, name := range names {
//        if idx > 10 {
//            return
//        }
//        util.Logger.Info("analysis: %s", name)
//        if strings.Contains(name, "code") && strings.Contains(name, "_daily") {
//            var data []model.QuoteItem
//            mongo.DB.C(name).Find(bson.M{}).All(&data)
//            for idx, item := range data {
//                if idx == len(data) - 1{
//                    continue
//                }
//                rows := []string{}
//
//                var v string
//                v = strconv.FormatFloat(item.Close_chg, 'f', 2, 64)
//                rows = append(rows, v)
//                v = strconv.FormatFloat(item.Ma_3_chg_pct, 'f', 2, 64)
//                rows = append(rows, v)
//                v = strconv.FormatFloat(item.Ma_5_chg_pct, 'f', 2, 64)
//                rows = append(rows, v)
//                v = strconv.FormatFloat(item.Ma_10_chg_pct, 'f', 2, 64)
//                rows = append(rows, v)
//                v = strconv.FormatFloat(item.Ma_20_chg_pct, 'f', 2, 64)
//                rows = append(rows, v)
//                v = strconv.FormatFloat(item.Rsv_9, 'f', 2, 64)
//                rows = append(rows, v)
//                v = strconv.FormatFloat(item.Rsv_K_9, 'f', 2, 64)
//                rows = append(rows, v)
//                v = strconv.FormatFloat(item.Rsv_D_9, 'f', 2, 64)
//                rows = append(rows, v)
//                v = strconv.FormatFloat(item.Tr, 'f', 2, 64)
//                rows = append(rows, v)
//                v = strconv.FormatFloat(item.Atr_14, 'f', 2, 64)
//                rows = append(rows, v)
//                v = strconv.FormatFloat(item.Bias_3_6, 'f', 2, 64)
//                rows = append(rows, v)
//                v = strconv.FormatFloat(item.Bias_10_20, 'f', 2, 64)
//                rows = append(rows, v)
//                v = strconv.FormatFloat(item.Macd_dem, 'f', 2, 64)
//                rows = append(rows, v)
//                v = strconv.FormatFloat(item.Macd_osc, 'f', 2, 64)
//                rows = append(rows, v)
//                v = strconv.FormatFloat(item.Rsi, 'f', 2, 64)
//                rows = append(rows, v)
//                v = strconv.FormatFloat(item.Ema_U_14, 'f', 2, 64)
//                rows = append(rows, v)
//                v = strconv.FormatFloat(item.Ema_D_14, 'f', 2, 64)
//                rows = append(rows, v)
//                v = strconv.FormatFloat(item.Boll_up_20_chg_pct, 'f', 2, 64)
//                rows = append(rows, v)
//                v = strconv.FormatFloat(item.Boll_dn_20_chg_pct, 'f', 2, 64)
//                rows = append(rows, v)
//                v = strconv.FormatFloat(item.Boll_pct_b, 'f', 2, 64)
//                rows = append(rows, v)
//                v = strconv.FormatFloat(item.Boll_bandwidth, 'f', 2, 64)
//                rows = append(rows, v)
//
//                //if item.Macd_dem < 0 {
//                //    rows = append(rows, "1")
//                //} else if item.Macd_dem > 0 {
//                //    rows = append(rows, "-1")
//                //} else {
//                //    rows = append(rows, "0")
//                //}
//                //
//                ////if item.Bias_6 < -0.03 {
//                ////    rows = append(rows, "1")
//                ////} else if item.Bias_6 > 0.035{
//                ////    rows = append(rows, "-1")
//                ////} else {
//                ////    rows = append(rows, "0")
//                ////}
//                ////
//                ////if item.Bias_20 < -0.07 {
//                ////    rows = append(rows, "1")
//                ////} else if item.Bias_20 > 0.08 {
//                ////    rows = append(rows, "-1")
//                ////} else {
//                ////    rows = append(rows, "0")
//                ////}
//                //
//                //if item.Rsi < 0.3 {
//                //    rows = append(rows, "1")
//                //} else if item.Rsi > 0.7 {
//                //    rows = append(rows, "-1")
//                //} else {
//                //    rows = append(rows, "0")
//                //}
//                //
//                //if item.Rsv_K_9 > 0.8 && item.Rsv_D_9 > 0.8 {
//                //    rows = append(rows, "-1")
//                //} else if item.Rsv_K_9 < 0.2 && item.Rsv_D_9 < 0.2 {
//                //    rows = append(rows, "1")
//                //} else if item.Rsv_K_9 > item.Rsv_D_9 {
//                //    rows = append(rows, "1")
//                //} else if item.Rsv_K_9 < item.Rsv_D_9 {
//                //    rows = append(rows, "-1")
//                //} else {
//                //    rows = append(rows, "0")
//                //}
//                //
//                //if item.Atr_14_chg > 0 {
//                //    rows = append(rows, "1")
//                //} else if item.Atr_14_chg < 0 {
//                //    rows = append(rows, "-1")
//                //} else {
//                //    rows = append(rows, "0")
//                //}
//
//                nextItem := data[idx + 1]
//                if nextItem.Close_chg > 0 {
//                    rows = append(rows, "UP")
//                } else if nextItem.Close_chg < 0 {
//                    rows = append(rows, "DN")
//                } else {
//                    rows = append(rows, "NA")
//                }
//
//                if  float64(idx) / float64(len(data)) < 0.8 {
//                    train_writer.Write(rows)
//                } else {
//                    test_writer.Write(rows)
//                }
//            }
//        }
//    }
//
//    /**
//        MACD
//            macd_dem > 0  买入
//            macd_dem < 0 卖出
//
//        Boll
//            连续突破up，卖出
//            连续突破down，买入
//
//        Bias
//            Bias_6 < -0.03  买入
//            Bias_6 > 0.035 卖出
//            Bias_20 < -0.07 买入
//            Bias_20 > 0.08 卖出
//
//        Rsi
//            > 0.7 卖出
//            < 0.3 买入
//
//        RSV
//            if: Rsv_K_9 > 0.8 && Rsv_D_9 > 0.8 卖出
//            else:  Rsv_K_9 < 0.2 && Rsv_D_9 < 0.2 买入
//            else: Rsv_K_9 > Rsv_D_9 买入
//            else: Rsv_K_9 < Rsv_D_9 卖出
//
//        真实波动率
//            Atr_14_chg > 0 买入
//            Atr_14_chg < 0 卖出
//
//
//     */
//
//}
