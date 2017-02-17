package ana

import (
    "gopkg.in/mgo.v2/bson"
    "bycrodata_grab/module/mongo"
    //"github.com/d4l3k/talib"
)

type QueryDataResult struct {
    //Id     interface{} `bson:"_id"`
    Open   []float64
    Close  []float64
    High   []float64
    Low    []float64
    Volume []float64
    Ts     []int64
    Date   []string
}

var (
    baseConditionSort = bson.M{
        "$sort": bson.M{
            "ts": 1,
        },
    }
    baseConditionFields = bson.M{
        "$group": bson.M{
            "_id": nil,
            "open": bson.M{
                "$push": "$open",
            },
            "close": bson.M{
                "$push": "$close",
            },
            "high": bson.M{
                "$push": "$high",
            },
            "low": bson.M{
                "$push": "$low",
            },
            "volume": bson.M{
                "$push": "$volume",
            },
            "ts": bson.M{
                "$push": "$ts",
            },
            "date": bson.M{
                "$push": "$date",
            },
        },
    }
)

func QueryData(code, dateStr string) QueryDataResult {
    //dailyCollectionName := "code_" + code + "_daily"
    minutelyCollectionName := "code_" + code + "_minutely"
    //resDaily := QueryDailyData(dailyCollectionName, dateStr)
    resMinutely := QueryMinutelyData(minutelyCollectionName, dateStr)

    return resMinutely
    //res := QueryDataResult{}
    //res.Open = append(resDaily.Open, resMinutely.Open...)
    //res.Close = append(resDaily.Close, resMinutely.Close...)
    //res.High = append(resDaily.High, resMinutely.High...)
    //res.Low = append(resDaily.Low, resMinutely.Low...)
    //res.Volume = append(resDaily.Volume, resMinutely.Volume...)
    //res.Date = append(resDaily.Date, resMinutely.Date...)
    //res.Ts = append(resDaily.Ts, resMinutely.Ts...)
    //
    //return res
}

func QueryDailyData(collectionName, dateStr string) QueryDataResult {
    condition := []bson.M{
        bson.M{
            "$match": bson.M{
                "date": bson.M{
                    "$lt": dateStr,
                },
            },
        },
        baseConditionSort,
        baseConditionFields,
    }
    res := QueryDataResult{}
    mongo.DB.C(collectionName).Pipe(condition).One(&res)
    return res
}

func QueryMinutelyData(collectionName, dateStr string) QueryDataResult {
    condition := []bson.M{
        bson.M{
            "$match": bson.M{
                "date": bson.M{
                    "$regex": dateStr,
                },
            },
        },
        baseConditionSort,
        baseConditionFields,
    }
    res := QueryDataResult{}
    mongo.DB.C(collectionName).Pipe(condition).One(&res)
    return res
}
