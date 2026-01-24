#!/bin/bash

work_path=$(cd `dirname $0`; pwd)
bin_file=${work_path}/${1}
monitor_file=${work_path}/log/${1}${3}_monitor.log

if [ $# -lt 4 ]; then
    echo "eg: start.sh gate ./config.yaml 1 debug"
    exit 1
fi

if [ ! -x ${bin_file} ]; then
    echo "${bin_file}不是可执行文件"
    exit 1
fi

mkdir -p ${work_path}/log
nohup ${bin_file} -mode=${4} -config=${2} -id=${3} 2>${monitor_file} 1>/dev/null &
echo "启动服务成功：${bin_file} -mode=${4} -config=${2} -id=${3}"
exit 0