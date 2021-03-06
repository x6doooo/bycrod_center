package yahoo

import (
    "bycrod_center/module/mongo"
    "gopkg.in/mgo.v2/bson"
    //"strings"
    //"bycrod_center/module/util"
)

func GetCodes() (codes []string, err error) {
    var stocks []bson.M
    //return
    err = mongo.DB.C("summaries").Find(bson.M{
        "volume": bson.M{
            "$gt": 100 * 10000,
        },
        "current": bson.M{
            "$lte": 50,
            "$gte": 3,
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

