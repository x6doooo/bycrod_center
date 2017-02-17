package util

import (
    "bycrod_center/conf"
    "github.com/x6doooo/smog"
)

var (
    Logger smog.LoggerInterface
)

func init() {
    logConf := conf.MainConf.Log;
    Logger = smog.NewLogger(logConf.File, logConf.Max_line, logConf.Backups, conf.IsDevMode)
}