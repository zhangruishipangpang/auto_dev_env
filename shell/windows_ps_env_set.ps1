# 新增用户环境变量
param (
    [string]$key,
    [string]$value
)

# 检查参数是否为空
if ([string]::IsNullOrEmpty($Key) -or [string]::IsNullOrEmpty($Value)) {
    Write-Error "Both -Key and -Value parameters are required."
    exit -1
}

Write-Output "env ->  key:$key   value:$value"

# 添加用户环境变量
[System.Environment]::SetEnvironmentVariable($key, $Value, [System.EnvironmentVariableTarget]::User)

$newValue = [System.Environment]::GetEnvironmentVariable($key, [System.EnvironmentVariableTarget]::User)

Write-Output "success value --> $newValue"