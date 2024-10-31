#!/bin/bash

# 检查是否提供了参数
if [ -z "$1" ]; then
  echo "Usage: $0 <environment_variable_key>"
  exit 1
fi

# 获取环境变量的键
KEY=$1

# 获取并输出环境变量的值
VALUE=$(printenv "$KEY")

if [ -z "$VALUE" ]; then
  echo ""
else
  echo "$VALUE"
fi