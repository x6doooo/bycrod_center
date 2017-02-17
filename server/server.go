package server

import (
    "time"
    "github.com/gin-gonic/gin"
    "github.com/x6doooo/err_handler"
    "net/http"
    "bycrod_center/module/util"
    "bycrod_center/conf"
    "bycrod_center/controller/api"
)


func RequestLog() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()

        ip := c.ClientIP()

        c.Next()

        end := time.Now()
        latency := end.Sub(start)
        util.Logger.Info("[%d] %s %s %s %s",
            c.Writer.Status(), ip,  c.Request.Method, c.Request.RequestURI, latency.String())
    }
}

func ErrHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        var err error
        defer err_handler.Recover(&err, func() {
            if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{
                    "data": err.Error(),
                })
                c.Abort()
            }
        })
        c.Next()
    }
}


func Start() {
    engine := gin.New()
    engine.Use(ErrHandler())
    engine.Use(gin.Recovery())

    // request log
    engine.Use(RequestLog())

    // http://123.56.128.80:50090/api/task/daily-format
    apiRouter := engine.Group("/api/task")
    {
        apiRouter.GET("/xq-init", api.TaskXueqiuInit)
        apiRouter.GET("/xq-events", api.TaskGetXqEvents)
        apiRouter.GET("/format", api.TaskFormat)
        apiRouter.GET("/clean", api.TaskClean)
        apiRouter.GET("/realtime", api.TaskRealtime)
        //apiRouter.GET("/knn", api.TaskDataKnnAnalysis)
        //apiRouter.GET("/noname", api.TaskDataAnalysis_Noname_0)
        //apiRouter.GET("/knn-test", api.TaskDataFormatTest)
    }
    //apiAnaRouter := engine.Group("/api/ana")
    //{
    //    apiAnaRouter.GET("/bbands", api.AnaBBands)
    //}

    util.Logger.Info("server start! %s", conf.MainConf.Server.Addr)
    engine.Run(conf.MainConf.Server.Addr)
}

