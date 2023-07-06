#!/bin/sh
sudo service mysql start
cd ~/redis-6.0.8
sudo ./src/redis-server redis.conf