#!/bin/sh

top -c | grep -w 'PID' -A 10 > watch.txt