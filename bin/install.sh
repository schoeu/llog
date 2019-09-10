#!/bin/bash
set -e

sys_bit=`getconf LONG_BIT`

if [ "$sys_bit" = "64" ]
then
  # wget 64bit version.
  wget https://github.com/schoeu/nma/blob/master/bin/nma_64bit?raw=true
  mv nma_64bit nma
elif [ "$sys_bit" = "32" ]
then
  # wget 32bit version.
  wget https://github.com/schoeu/nma/blob/master/bin/nma_32bit?raw=true
  mv nma_32bit nma
fi

chmod +x nma
kill -9 `ps aux | grep nma | head -n 1 | awk '{print $2}'`
nohup ./nma >> ./nma_nohup.log 2>&1 &
