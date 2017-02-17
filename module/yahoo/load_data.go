package yahoo

import (
    "time"
    "bycrodata_grab/module/util"
    "bycrodata_grab/conf"
    "bycrodata_grab/model"
    "net/http"
    "io"
    "compress/gzip"
    "io/ioutil"
    "encoding/json"
    "errors"
    "bycrodata_grab/module/mongo"
    "gopkg.in/mgo.v2/bson"
)

type msgData struct {
    Code string
    Err  error
    Res  []model.QuoteItem
}
type apiData struct {
    Code int `json:"code"`
    Data []model.QuoteItem `json:"data"`
    Err  string `json:"err"`
}

var client *http.Client

func init() {
    client = &http.Client{
        Timeout: time.Second * 5,
    }
}

const (
    interval_daily = "1d"
    the_range_daily = "3072d"
    //the_range_daily = "2d"

    interval_hourly = "1h"
    the_range_hourly = "1024d"

    interval_minutely = "1m"
    the_range_minutely = "100d"

    interval_realtime = "1m"
    the_range_realtime = "5m"
)

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

    grabs := conf.MainConf.Backend.Grabs

    bufSize := len(codes)
    jobs := make(chan string, bufSize)
    results := make(chan msgData, bufSize)

    for idx, grabAddr := range grabs {
        go graber(idx, jobs, results, grabAddr, interval, the_range)
    }

    for _, code := range codes {
        jobs <- code
    }
    close(jobs)

    for i := 0; i < bufSize; i++ {
        data := <-results
        if data.Err != nil {
            fails = append(fails, data.Code)
        } else {

            // 处理返回数据
            save(data.Res, data.Code, dataType)

        }
        util.Logger.Info("%s %d/%d %s", data.Code, i, bufSize, time.Since(startTime).String())
    }
    close(results)

    return
}

func save(data []model.QuoteItem, code, dataType string) {
    collectionName := "code_" + code + "_" + dataType
    var theLastTs uint64 = 0
    var dataListHasBeenInserted []model.QuoteItem
    mongo.DB.C(collectionName).Find(bson.M{}).Sort("-ts").All(&dataListHasBeenInserted)
    if len(dataListHasBeenInserted) != 0 {
        theLastTs = dataListHasBeenInserted[0].Ts
    }
    dataList := []interface{}{}//model.QuoteItem{}
    for _, d := range data {
        if d.Ts > theLastTs {
            dataList = append(dataList, d)
        }
    }
    if len(dataList) > 0 {
        mongo.DB.C(collectionName).Insert(dataList...)
    }
}

func graber(id int, jobs chan string, results chan msgData, grabAddr, interval, the_range string) {
    theUrl := "http://" + grabAddr +
        "/api/yahoo/fetch?interval=" + interval + "&range=" + the_range + "&code="
    for code := range jobs {
        util.Logger.Info("grab %d: start fetch %s, start...", id, code)
        res, err := get(theUrl + code)
        if err != nil {
            util.Logger.Info("grab %d: start fetch %s, error! %s", id, code, err.Error())
        } else {
            util.Logger.Info("grab %d: start fetch %s, end!", id, code)
        }
        msg := msgData{
            Code: code,
            Res: res,
            Err: err,
        }
        results <- msg
    }
}

func get(theUrl string) (respData []model.QuoteItem, err error) {
    req, _ := http.NewRequest("GET", theUrl, nil)
    var resp *http.Response
    resp, err = client.Do(req)
    if err != nil {
        return
    }

    var reader io.ReadCloser
    switch resp.Header.Get("Content-Encoding") {
    case "gzip":
        reader, err = gzip.NewReader(resp.Body)
        if err != nil {
            return
        }
    default:
        reader = resp.Body
    }
    var body []byte
    body, err = ioutil.ReadAll(reader)
    if err != nil {
        return;
    }
    reader.Close()
    resp.Body.Close()

    var theApiData apiData
    err = json.Unmarshal(body, &theApiData)
    if err != nil {
        return
    }
    //fmt.Println(theApiData, 222)
    if theApiData.Code != 0 {
        err = errors.New(theApiData.Err)
        return
    }
    respData = theApiData.Data
    return;
}


