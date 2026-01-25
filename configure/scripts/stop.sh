#!/bin/bash

work_path=$(cd `dirname $0`; pwd)
bin_file=${pwd_path}/${1}

if [ $# -lt 2 ]; then
    echo "stop.sh [服务] [节点id]"
    exit 1
fi

PIDS=$(ps -ef | grep ${bin_file} | grep "id=${3}" | grep -v grep | awk '{print $2}')
if [ -z "${PIDS}" ]; then
    echo "服务已经关闭 ${1}"
    exit 1
fi

kill -15 ${PIDS} || true
echo "关闭服务：${1}${2}"
exit 0