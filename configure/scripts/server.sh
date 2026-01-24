#!/bin/bash

start(){
    ./start.sh gate ./config.yaml 1 release
    ./start.sh db ./config.yaml 1 release
    ./start.sh game ./config.yaml 1 release
}

stop(){
    ./stop.sh gate 1
    ./stop.sh db 1
    ./stop.sh game 1
}

case $1 in
start)
    start
    ;;
stop)
    stop
    ;;
esac
