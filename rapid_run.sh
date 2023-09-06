#!/bin/sh

PID_FILE="/tmp/parkinglot.pid"

# Pid file is needless in container enviroment
rm -f $PID_FILE

./parkinglot -pid-file=$PID_FILE &
