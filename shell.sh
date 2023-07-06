#!/bin/bash

ulimit -c unlimited
PORT=7200
PM="slotjelly"
DIR="SlotJellyServer"
LOGDIR="/LogDir/${DIR}"
LOGFILE1="${LOGDIR}/pm2.${PORT}.log"
LOGFILE2="${LOGDIR}/pm2.1${PORT}.log"
CONFIG="/data/gameserver/${DIR}/config.yaml"

# group A
echo "sudo pm2 start --name ${PM}-a -e $LOGFILE1 -o $LOGFILE1 ./slots -- -game 1 -level 306 -port $PORT -f $CONFIG -log_dir $LOGFILE1  -alsologtostderr"

# group B
PORT=17200
echo "sudo pm2 start --name ${PM}-a -e $LOGFILE2 -o $LOGFILE2 ./slots -- -game 1 -level 306 -port $PORT -f $CONFIG -log_dir $LOGFILE2  -alsologtostderr"