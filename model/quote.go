package model

import (
    "gopkg.in/mgo.v2/bson"
    "reflect"
)

type QuoteItem struct {
    Id                     bson.ObjectId `bson:"_id,omitempty"`
    Ts                     int64
    Date                   string
    Volume                 uint64
    High                   float64
    Low                    float64
    Open                   float64
    Close                  float64
    //Close_chg              float64
    //
    //                               // moving average
    //Ma_3                   float64
    //Ma_3_chg_pct           float64
    //Ma_6                   float64
    //Ma_6_chg_pct           float64
    //
    //Ma_5                   float64
    //Ma_5_chg_pct           float64
    //Ma_10                  float64
    //Ma_10_chg_pct          float64
    //Ma_20                  float64
    //Ma_20_chg_pct          float64
    //
    //                               // RSV
    //Rsv_9                  float64
    //Rsv_K_9                float64
    //Rsv_D_9                float64
    //
    //                               // 真实波动率
    //Tr                     float64
    //Atr_14                 float64
    //Atr_14_chg             float64
    //
    //                               // bias
    //Bias_3_6               float64
    //Bias_10_20             float64
    //Bias_6                 float64
    //Bias_20                float64
    //
    //                               // MACD
    //                               //Ema_diff_9             float64
    //Ema_close_12           float64
    //Ema_close_26           float64
    //Macd_dem               float64
    //Macd_osc               float64
    //
    //                               // rsi
    //Ema_U_14               float64
    //Ema_D_14               float64
    //Rsi                    float64
    //
    //                               // bollinger bands
    //DiffSquare             float64 // diff = close - ma
    //Boll_up_20             float64
    //Boll_up_20_chg_pct     float64
    //Boll_dn_20             float64
    //Boll_dn_20_chg_pct     float64
    //Boll_pct_b             float64 // %b
    //Boll_pct_b_chg_pct     float64 // %b
    //Boll_bandwidth         float64
    //Boll_bandwidth_chg_pct float64
}

func (me *QuoteItem) SetFloat64ByFieldName(field string, value float64) (ok bool) {
    r := reflect.ValueOf(me)
    f := r.Elem().FieldByName(field)
    ok = f.IsValid() && f.CanSet() && f.Kind() == reflect.Float64
    if ok {
        f.SetFloat(value)
    }
    return
}

func (me *QuoteItem) GetFloat64ByFieldName(field string) (value float64, ok bool) {
    r := reflect.ValueOf(me)
    f := r.Elem().FieldByName(field)
    ok = f.IsValid()
    if ok {
        value = f.Float()
    }
    return
}

type QuoteData struct {
    Volume []*uint64
    High   []*float64
    Low    []*float64
    Open   []*float64
    Close  []*float64
}

type ResultItem struct {
    Meta       interface{}
    Timestamp  []int64
    Indicators map[string]([]QuoteData)
}

type ChartItem struct {
    Result []ResultItem
    Error  *string
}

type RespResult struct {
    Chart ChartItem
}