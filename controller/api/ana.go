package api

import (
    "github.com/gin-gonic/gin"
    "bycrodata_grab/module/ana"
    "net/http"
    "github.com/d4l3k/talib"
    //"bycrodata_grab/module/util"
)

func AnaBBands(ctx *gin.Context) {
    code := ctx.Query("code")
    //collectionName = "code_" + collectionName + "_minutely"
    dateStr := ctx.Query("date");

    baseData := ana.QueryData(code, dateStr)
    //util.Logger.Info("test: %v", baseData)

    up, mid, low := talib.BBands(baseData.Close, 5, 2, 2, 0)
    head := make([]float64, 4)
    up = append(head, up[:len(up) - 4]...)
    mid = append(head, mid[:len(mid) - 4]...)
    low = append(head, low[:len(low) - 4]...)

    res := gin.H{
        "code": 0,
        "data": gin.H{
            "quote": baseData,
            "up": up,
            "mid": mid,
            "low": low,
        },
    }
    ctx.JSON(http.StatusOK, res)
}
