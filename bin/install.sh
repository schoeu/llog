#!/bin/bash
set -e

# TODO download & start server.
sys_bit=`getconf LONG_BIT`

if [ "$sys_bit" = "64" ]
then
  # wget 64bit version.
elif [ "$sys_bit" = "32" ]
then
  # wget 32bit version.
fi
