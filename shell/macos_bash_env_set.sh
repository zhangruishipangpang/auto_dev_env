#!/bin/bash

# 检查是否提供了两个参数
if [ -z "$1" ] || [ -z "$2" ]; then
  echo "Usage: $0 <environment_variable_key> <environment_variable_value>"
  exit 1
fi

# 获取环境变量的键和值
KEY=$1
VALUE=$2

# 检查当前使用的 Shell
if [ -f ~/.bash_profile ]; then
  CONFIG_FILE=~/.bash_profile
elif [ -f ~/.zshrc ]; then
  CONFIG_FILE=~/.zshrc
else
  echo "No supported shell configuration file found. Please create one manually."
  exit 1
fi

# 添加环境变量到配置文件
echo "export $KEY=\"$VALUE\"" >> "$CONFIG_FILE"

# 提示用户重新加载配置文件
echo "Please reload your shell configuration by running:"
echo "  source $CONFIG_FILE"