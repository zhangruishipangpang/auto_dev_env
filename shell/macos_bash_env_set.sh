#! /bin/bash

##### 添加环境变量程序 #####
if [ ! "$1" ] && [ ! "$2" ]; then
    echo "key or value is null"
    return 1
fi

export "$1"="$2"

# shellcheck disable=SC1090
source ~/.bash_profile
