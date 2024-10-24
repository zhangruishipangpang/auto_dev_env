# 新增用户环境变量
param (
    [string]$key
)

# 检查参数是否为空
if ([string]::IsNullOrEmpty($Key)) {
    Write-Error "Both -Key parameters are required."
    exit -1
}

$value = [System.Environment]::GetEnvironmentVariable($key, [System.EnvironmentVariableTarget]::User)

Write-Output $value