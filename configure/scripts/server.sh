#!/bin/bash

start(){
    ./start.sh gate ./config.yaml 1 release
    sleep 1
    ./start.sh db ./config.yaml 1 release
    sleep 1
    ./start.sh game ./config.yaml 1 release
}

stop(){
    ./stop.sh gate 1
    sleep 1
    ./stop.sh db 1
    sleep 1
    ./stop.sh game 1
}

case $1 in
start)
    start
    ;;
restart)
    stop
    start
    ;;
stop)
    stop
    ;;
esac
