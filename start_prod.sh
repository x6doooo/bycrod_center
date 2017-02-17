#!/usr/bin/env bash
#ps aux|grep -v grep|grep bin\/bycrodata_grab_prod|awk '{print $2}'|xargs kill -9
#./bin/bycrodata_grab_prod --conf=./conf/conf.prod.toml
kill -9 `pgrep -f bin/bycrodata`
nohup ./bin/bycrod_center_prod --conf=./conf/conf.prod.toml $1 >/dev/null 2>&1 &

