package yahoo

import (
    "bycrodata_grab/module/mongo"
    "gopkg.in/mgo.v2/bson"
    //"strings"
    //"bycrodata_grab/module/util"
)

func GetCodes() (codes []string, err error) {
    var stocks []bson.M
    //return
    err = mongo.DB.C("summaries").Find(bson.M{
        "volume": bson.M{
            "$gt": 50 * 10000,
        },
        "current": bson.M{
            "$lte": 8, //50
            "$gte": 5, //3
        },
        //"instOwn": bson.M{
        //    "$lt": 50,
        //},
    }).All(&stocks)
    if err != nil {
        return
    }
    for _, item := range stocks {
        code := item["code"]
        switch code.(type) {
        case string:
            c := code.(string)
            codes = append(codes, c)
        }
    }
    return
}

