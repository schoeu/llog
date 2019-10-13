#!/bin/bash
set -e

sys_bit=`getconf LONG_BIT`

if [ "$sys_bit" = "64" ]
then
  # wget 64bit version.
  wget http://qiniucdn.schoeu.com/lla_64bit
  mv lla_64bit lla
elif [ "$sys_bit" = "32" ]
then
  # wget 32bit version.
  wget http://qiniucdn.schoeu.com/lla_32bit
  mv lla_32bit lla
fi

chmod +x lla
#kill -9 `ps aux | grep lla | head -n 1 | awk '{print $2}'`
#nohup ./lla >> ./lla_nohup.log 2>&1 &
