package knn

//import (
//    "bycrod_center/module/mongo"
//    "gopkg.in/mgo.v2/bson"
//    "strings"
//    "bycrod_center/model"
//    "bycrod_center/module/util"
//    "sort"
//)
//
//type forecastCount struct {
//    Name string
//    Right float64
//    Total float64
//}
//
//type ByRateSort []forecastCount
//func (me ByRateSort) Len() int {
//    return len(me)
//}
//func (me ByRateSort) Swap(i, j int) {
//    me[i], me[j] = me[j], me[i]
//}
//func (me ByRateSort) Less(i, j int) bool {
//    ra := me[i].Right / me[i].Total
//    rb := me[j].Right / me[j].Total
//    return ra > rb
//}
//
//func DataFormatTest(dataType string) {
//    names, _ := mongo.DB.CollectionNames()
//    name := ""
//    countByCode := map[string]forecastCount{}
//    for _, n := range names {
//        if strings.Contains(n, "code_") && strings.Contains(n, "_" + dataType) {
//            name = n
//            countByCode[name] = forecastCount{
//                Name: name,
//            }
//            //break
//        }
//    }
//
//    dates := []map[string]string{}
//    mongo.DB.C(name).Pipe([]bson.M{
//        bson.M{
//            "$group": bson.M{
//                "_id": "$date",
//            },
//        },
//        bson.M{
//            "$sort": bson.M{
//                "_id": -1,
//            },
//        },
//        bson.M{
//            "$limit": 10,
//        },
//    }).All(&dates)
//
//    dates = dates[2:]
//
//    count := map[string]float64{
//        "total": 0,
//        "right": 0,
//    }
//
//
//    for _, d := range dates {
//        //fmt.Println(d)
//        DataAnalysis(d["_id"], 0, dataType)
//
//        res := []forecastResultItem{}
//        var cn string
//        switch dataType {
//        case "daily":
//            cn = forecast_daily_collection
//        case "hourly":
//            cn = forecast_hourly_collection
//        }
//        mongo.DB.C(cn).Find(bson.M{
//            "next_close_chg": bson.M{
//                "$gt": 0,
//            },
//            "next_next_close_chg": bson.M{
//                "$gt": 0,
//            },
//        }).Sort("-total_chg").Limit(10).All(&res)
//        //fmt.Println(res)
//        for _, r := range res {
//            shouldBe := []model.QuoteItem{}
//
//            mongo.DB.C(r.Name).Find(bson.M{
//                "date": bson.M{
//                    "$gt": d["_id"],
//                },
//            }).Limit(1).All(&shouldBe)
//
//
//            baseNum := float64(1)
//            for _, item := range shouldBe {
//                baseNum *= (1 + item.Close_chg)
//            }
//
//
//            count["total"] += 1
//            codeCountItem := countByCode[r.Name]
//            codeCountItem.Total += 1
//            if baseNum > 1 {
//                codeCountItem.Right += 1
//                count["right"] += 1
//            }
//            countByCode[r.Name] = codeCountItem
//            util.Logger.Info("...............................", codeCountItem)
//
//            //if _, ok := countByCode[r.Name]; !ok {
//            //    countByCode[r.Name] = forecastCount{
//            //        Name: r.Name,
//            //        Right: 0,
//            //        Total: 0,
//            //    }
//            //}
//
//            //op := r.Total_chg * shouldBe.Close_chg
//            ////fmt.Println(r.Name, shouldBe)
//            //count["total"] += 1
//            //codeCountItem := countByCode[r.Name]
//            //codeCountItem.Total += 1
//            //if op > 0 {
//            //    count["right"] += 1
//            //    codeCountItem.Right += 1
//            //} else if op == 0 && r.Total_chg ==0 && shouldBe.Close_chg == 0 {
//            //    count["right"] += 1
//            //    codeCountItem.Right += 1
//            //}
//            //countByCode[r.Name] = codeCountItem
//            //util.Logger.Info("...............................", codeCountItem)
//        }
//        util.Logger.Info("total: %f", count["total"])
//        util.Logger.Info("right: %f", count["right"])
//        util.Logger.Info("rate: %f%", count["right"] / count["total"] * 100)
//    }
//
//
//
//
//    util.Logger.Info("--------------------------")
//    util.Logger.Info("total: %f", count["total"])
//    util.Logger.Info("right: %f", count["right"])
//    util.Logger.Info("rate: %f%", count["right"] / count["total"] * 100)
//
//
//    arrTem := []forecastCount{}
//    for _, item := range countByCode {
//        if item.Total != 0 && item.Right != 0 {
//            arrTem = append(arrTem, item)
//        }
//    }
//
//    sort.Sort(ByRateSort(arrTem))
//    size := len(arrTem)
//    if size > 50 {
//        size = 50
//    }
//
//    util.Logger.Info("--------------------------")
//    for i := 0; i < size; i++ {
//        item := arrTem[i]
//        rate := item.Right / item.Total * 100
//        util.Logger.Info("%d: %s  %f% (%f/%f)", i, item.Name, rate, item.Right, item.Total)
//    }
//
//    //fmt.Println("total:", count["total"])
//    //fmt.Println("right:", count["right"])
//    //fmt.Println("rate:", count["right"] / count["total"] * 100, "%")
//
//}
