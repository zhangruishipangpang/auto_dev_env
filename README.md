# 自动初始化开发环境变量工具

> 改工具旨在可以快速初始化本地的开发环境变量，消除因为下载环境变量、配置环境变量而花费的时间

**配置文件：**
```json
{
  "default_zip_dir": "默认的环境变量压缩包路径，如果configs中元素使用了use_default:true，则会使用压缩包覆盖source_path中的配置", 
  "configs": [
    {
      "env_name": "环境名称",
      "env_code": "环境的文件code，压缩包名称|资源文件路径 都需要拼接该字段以确定文件夹|文件",
      "env_source_path": "环境资源的文件夹 - 环境变量的包",
      "env_target_path": "环境资源的文件夹 - 想要最终将环境变量包放置的文件夹",
      "use_default": true, /*是否使用zip_dir*/
      "del_source": false, /*复制文件从 source-> target 后是否要删除 source 资源*/
      "env_source_check": [ /*检查文件项*/
        {
          "name" : "bin", /*名称*/
          "type" : "dir", /*类型：dir|file*/
          "path" : "$bin" /*路径。$代表是env_source_path下的文件夹*/
        }
      ],
      "env_config" : [ /*环境变量配置项*/
        {
          "key" : "TEST_JAVA_HOME", /*环境变量KEY*/
          "value" : "$", /*环境变量值。 $代表在env_target_path下*/
          "cover" : true, /*如果已经存在是否覆盖*/
          "append_path" : true, /*是否追加到path中*/
          "suffix_path" : ["bin"] /*追加到path中时，value后拼接的子路径*/
        }
      ]
    }
  ]
}
```

## Error 

### 执行 setEnv|getEv 的.ps1返回不正常码值(非0)

```shell
# 检查 PowerShell 安全策略
Get-ExecutionPolicy

#常见的执行策略包括：
#Restricted：不允许运行任何脚本。
#AllSigned：仅允许运行已由受信任的发布者签名的脚本。
#RemoteSigned：允许运行本地创建的脚本，远程下载的脚本必须由受信任的发布者签名。
#Unrestricted：允许运行所有脚本，但会提示用户确认是否运行未签名的脚本。
#Bypass：不阻止任何脚本的执行，也不会显示任何提示。
#Undefined：没有设置执行策略。

# 将策略改成 Unrestricted
Set-ExecutionPolicy Unrestricted
```