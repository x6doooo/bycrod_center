package xueqiu

import (
    "bycrodata_grab/conf"
    "github.com/x6doooo/xueqiu_api"
    "strings"
    "bycrodata_grab/module/mongo"
    "bycrodata_grab/module/util"
)

func InitStockList() {

    mongo.DB.C("summaries").DropCollection()

    controller := xueqiu_api.New(conf.MainConf.Xueqiu.Username, conf.MainConf.Xueqiu.Password)
    controller.Login()

    list := controller.GetCodeList()
    listSize := len(list);

    //fmt.Println(listSize)
    for i := 0; i < listSize; i += 50 {

        util.Logger.Info("init xueqiu stock list %d / %d", i, listSize)

        var codes []string
        if listSize - i < 50 {
            codes = list[i:]
        } else {
            codes = list[i:i + 50]
        }
        details := controller.GetDetail(strings.Join(codes, ","))
        mongo.DB.C("summaries").Insert(details...)
    }
}
